package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/432539/gpt2api/internal/account"
	"github.com/432539/gpt2api/internal/apikey"
	"github.com/432539/gpt2api/internal/audit"
	"github.com/432539/gpt2api/internal/auth"
	"github.com/432539/gpt2api/internal/config"
	"github.com/432539/gpt2api/internal/db"
	"github.com/432539/gpt2api/internal/gateway"
	"github.com/432539/gpt2api/internal/image"
	modelpkg "github.com/432539/gpt2api/internal/model"
	"github.com/432539/gpt2api/internal/proxy"
	gwratelimit "github.com/432539/gpt2api/internal/ratelimit"
	"github.com/432539/gpt2api/internal/scheduler"
	"github.com/432539/gpt2api/internal/server"
	"github.com/432539/gpt2api/internal/settings"
	"github.com/432539/gpt2api/internal/usage"
	"github.com/432539/gpt2api/internal/user"
	"github.com/432539/gpt2api/pkg/crypto"
	pkgjwt "github.com/432539/gpt2api/pkg/jwt"
	"github.com/432539/gpt2api/pkg/lock"
	"github.com/432539/gpt2api/pkg/logger"
	"github.com/432539/gpt2api/pkg/mailer"
	pkgratelimit "github.com/432539/gpt2api/pkg/ratelimit"
)

var (
	configPath = flag.String("c", "configs/config.yaml", "config file path")
	showVer    = flag.Bool("v", false, "show version and exit")
)

var (
	version   = "0.2.0-dev"
	buildTime = "unknown"
)

func main() {
	flag.Parse()
	if *showVer {
		fmt.Printf("gpt2api %s (build %s)\n", version, buildTime)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Format, cfg.Log.Output); err != nil {
		fmt.Fprintf(os.Stderr, "init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	log := logger.L()
	log.Info("boot gpt2api",
		zap.String("version", version),
		zap.String("env", cfg.App.Env),
		zap.String("listen", cfg.App.Listen),
	)

	sqldb, err := db.NewSQLite(cfg.SQLite)
	if err != nil {
		log.Fatal("sqlite init", zap.Error(err))
	}
	defer sqldb.Close()

	if err := db.Migrate(sqldb); err != nil {
		log.Fatal("db migrate", zap.Error(err))
	}

	cipher, err := crypto.NewAESGCM(cfg.Crypto.AESKey)
	if err != nil {
		log.Fatal("crypto init", zap.Error(err))
	}

	// ---- DAO / Service ----
	userDAO := user.NewDAO(sqldb)
	proxyDAO := proxy.NewDAO(sqldb)
	proxySvc := proxy.NewService(proxyDAO, cipher)

	accDAO := account.NewDAO(sqldb)
	accSvc := account.NewService(accDAO, cipher)

	keyDAO := apikey.NewDAO(sqldb)
	keySvc := apikey.NewService(keyDAO)

	modelDAO := modelpkg.NewDAO(sqldb)
	modelReg := modelpkg.NewRegistry(modelDAO)
	if err := modelReg.Preload(context.Background()); err != nil {
		log.Warn("model preload failed", zap.Error(err))
	}

	rl := lock.NewMemoryLock()
	sched := scheduler.New(accSvc, proxySvc, rl, cfg.Scheduler)

	tb := pkgratelimit.NewMemoryBucket()
	limiter := gwratelimit.New(tb)

	usageLogger := usage.New(sqldb, usage.Options{})
	defer usageLogger.Close()

	jm := pkgjwt.NewManager(pkgjwt.Config{
		Secret:        cfg.JWT.Secret,
		Issuer:        cfg.JWT.Issuer,
		AccessTTLSec:  cfg.JWT.AccessTTLSec,
		RefreshTTLSec: cfg.JWT.RefreshTTLSec,
	})
	authSvc := auth.NewService(userDAO, jm, cfg.Security.BcryptCost)

	gwH := &gateway.Handler{
		Models:    modelReg,
		Keys:      keySvc,
		Scheduler: sched,
		Limiter:   limiter,
		Usage:     usageLogger,
		AccSvc:    accSvc,
	}

	imageDAO := image.NewDAO(sqldb)
	imageRunner := image.NewRunner(sched, imageDAO)
	imageRunner.SetQuotaDecrementor(accDAO)
	imageRunner.SetCacheDir("data")
	imagesH := &gateway.ImagesHandler{
		Handler:  gwH,
		Runner:   imageRunner,
		DAO:      imageDAO,
		CacheDir: "data",
	}
	gwH.Images = imagesH

	image.SetProxyURLBuilder(func(taskID string, idx int) string {
		return gateway.BuildImageProxyURL(taskID, idx, gateway.ImageProxyTTL)
	})

	auditDAO := audit.NewDAO(sqldb)
	auditH := audit.NewHandler(auditDAO)

	adminUserH := user.NewAdminHandler(userDAO, authSvc, auditDAO)

	adminModelH := modelpkg.NewAdminHandler(modelDAO, modelReg, auditDAO)
	adminKeyH := apikey.NewAdminHandler(keySvc, keyDAO, sqldb)
	usageQDAO := usage.NewQueryDAO(sqldb)
	adminUsageH := usage.NewAdminHandler(usageQDAO)
	meUsageH := usage.NewMeHandler(usageQDAO)
	meImageH := image.NewMeHandler(imageDAO)
	adminImageH := image.NewAdminHandler(imageDAO)

	mailSvc := mailer.New(mailer.Config{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
		FromName: cfg.SMTP.FromName,
		UseTLS:   cfg.SMTP.UseTLS,
	}, log)
	defer mailSvc.Close()
	if mailSvc.Disabled() {
		log.Info("mail channel disabled (smtp.host empty)")
	} else {
		log.Info("mail channel ready", zap.String("host", cfg.SMTP.Host))
	}
	authSvc.SetMailer(mailSvc, cfg.App.BaseURL)

	settingsDAO := settings.NewDAO(sqldb)
	settingsSvc := settings.NewService(settingsDAO)
	if err := settingsSvc.Reload(context.Background()); err != nil {
		log.Warn("settings reload failed, using defaults", zap.Error(err))
	}
	settingsH := settings.NewHandler(settingsSvc, mailSvc, auditDAO)
	authSvc.SetSettings(settingsSvc)

	keySvc.SetSettings(settingsSvc)
	gwH.Settings = settingsSvc
	sched.SetRuntime(scheduler.RuntimeParams{
		DailyUsageRatio: settingsSvc.DailyUsageRatio,
		Cooldown429Sec:  settingsSvc.Cooldown429Sec,
		WarnedPauseHrs:  settingsSvc.WarnedPauseHours,
		QueueWaitSec:    settingsSvc.DispatchQueueWaitSec,
	})
	jm.SetTTLProvider(func() (int, int) {
		return settingsSvc.JWTAccessTTLSec(), settingsSvc.JWTRefreshTTLSec()
	})

	proxyH := proxy.NewHandler(proxySvc)
	prober := proxy.NewProber(proxySvc, settingsSvc, log.Named("proxy-prober"))
	proxyH.SetProber(prober)

	proberCtx, cancelProber := context.WithCancel(context.Background())
	defer cancelProber()
	go prober.Run(proberCtx)

	accountH := account.NewHandler(accSvc)
	accRefresher := account.NewRefresher(accSvc, settingsSvc, log.Named("account-refresh"))
	accQuota := account.NewQuotaProber(accSvc, settingsSvc, log.Named("account-quota"))

	acctProxyResolver := &accountProxyResolver{accSvc: accSvc, proxySvc: proxySvc}
	accRefresher.SetProxyResolver(acctProxyResolver)
	accQuota.SetProxyResolver(acctProxyResolver)

	accountH.SetRefresher(accRefresher)
	accountH.SetProber(accQuota)
	accountH.SetSettings(settingsSvc)
	accountH.SetProxyResolver(acctProxyResolver)

	imagesH.ImageAccResolver = acctProxyResolver

	accBgCtx, cancelAccBg := context.WithCancel(context.Background())
	defer cancelAccBg()
	go accRefresher.Run(accBgCtx)
	go accQuota.Run(accBgCtx)

	deps := &server.Deps{
		Config: cfg,
		JWT:    jm,

		AuthH: auth.NewHandler(authSvc),
		UserH: user.NewHandler(userDAO),

		KeySvc:   keySvc,
		KeyH:     apikey.NewHandler(keySvc),
		ProxyH:   proxyH,
		AccountH: accountH,

		GatewayH: gwH,
		ImagesH:  imagesH,

		AuditH:     auditH,
		AuditDAO:   auditDAO,
		AdminUserH: adminUserH,

		AdminModelH: adminModelH,
		AdminKeyH:   adminKeyH,
		AdminUsageH: adminUsageH,

		MeUsageH:    meUsageH,
		MeImageH:    meImageH,
		AdminImageH: adminImageH,

		SettingsH: settingsH,
	}

	r := server.New(deps)
	srv := &http.Server{
		Addr:              cfg.App.Listen,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("http server started", zap.String("addr", cfg.App.Listen))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http listen", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown", zap.Error(err))
	}
	log.Info("bye")
}

// accountProxyResolver wraps account + proxy services for resolving proxy URLs.
type accountProxyResolver struct {
	accSvc   *account.Service
	proxySvc *proxy.Service
}

func (r *accountProxyResolver) ProxyURLForAccount(ctx context.Context, accountID uint64) string {
	if r == nil || r.accSvc == nil || r.proxySvc == nil {
		return ""
	}
	b, err := r.accSvc.GetBinding(ctx, accountID)
	if err != nil || b == nil {
		return ""
	}
	p, err := r.proxySvc.Get(ctx, b.ProxyID)
	if err != nil || p == nil || !p.Enabled {
		return ""
	}
	u, err := r.proxySvc.BuildURL(p)
	if err != nil {
		return ""
	}
	return u
}

func (r *accountProxyResolver) ProxyURLByID(ctx context.Context, proxyID uint64) string {
	if r == nil || r.proxySvc == nil || proxyID == 0 {
		return ""
	}
	p, err := r.proxySvc.Get(ctx, proxyID)
	if err != nil || p == nil || !p.Enabled {
		return ""
	}
	u, err := r.proxySvc.BuildURL(p)
	if err != nil {
		return ""
	}
	return u
}

func (r *accountProxyResolver) AuthToken(ctx context.Context, accountID uint64) (string, string, string, error) {
	if r == nil || r.accSvc == nil {
		return "", "", "", fmt.Errorf("account service not ready")
	}
	a, err := r.accSvc.Get(ctx, accountID)
	if err != nil {
		return "", "", "", err
	}
	at, err := r.accSvc.DecryptAuthToken(a)
	if err != nil {
		return "", "", "", err
	}
	cookies, _ := r.accSvc.DecryptCookies(ctx, accountID)
	return at, a.OAIDeviceID, cookies, nil
}

func (r *accountProxyResolver) ProxyURL(ctx context.Context, accountID uint64) string {
	return r.ProxyURLForAccount(ctx, accountID)
}
