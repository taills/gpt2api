-- SQLite schema — combined from all migrations, idempotent (IF NOT EXISTS everywhere).
-- Generated for SQLite. No MySQL-specific syntax.

-- ============================================================
-- user_groups
-- ============================================================
CREATE TABLE IF NOT EXISTS user_groups (
    id                  INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    name                TEXT     NOT NULL,
    ratio               REAL     NOT NULL DEFAULT 1.0,
    daily_limit_credits INTEGER  NOT NULL DEFAULT 0,
    rpm_limit           INTEGER  NOT NULL DEFAULT 60,
    tpm_limit           INTEGER  NOT NULL DEFAULT 60000,
    remark              TEXT     NOT NULL DEFAULT '',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_user_groups_name ON user_groups(name);

INSERT OR IGNORE INTO user_groups (name, ratio, rpm_limit, tpm_limit, remark) VALUES
    ('default', 1.0,  60,  60000, '默认分组'),
    ('vip',     0.8, 300, 300000, 'VIP 分组,倍率 0.8'),
    ('svip',    0.6, 600, 600000, 'SVIP 分组,倍率 0.6');

-- ============================================================
-- users
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    email           TEXT     NOT NULL,
    password_hash   TEXT     NOT NULL,
    nickname        TEXT     NOT NULL DEFAULT '',
    group_id        INTEGER  NOT NULL DEFAULT 1,
    role            TEXT     NOT NULL DEFAULT 'user',
    status          TEXT     NOT NULL DEFAULT 'active',
    credit_balance  INTEGER  NOT NULL DEFAULT 0,
    credit_frozen   INTEGER  NOT NULL DEFAULT 0,
    version         INTEGER  NOT NULL DEFAULT 0,
    last_login_at   DATETIME,
    last_login_ip   TEXT     NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_users_email  ON users(email);
CREATE        INDEX IF NOT EXISTS idx_users_group  ON users(group_id);
CREATE        INDEX IF NOT EXISTS idx_users_status ON users(status);

-- ============================================================
-- api_keys
-- ============================================================
CREATE TABLE IF NOT EXISTS api_keys (
    id              INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER  NOT NULL,
    name            TEXT     NOT NULL DEFAULT '',
    key_prefix      TEXT     NOT NULL,
    key_hash        TEXT     NOT NULL,
    quota_limit     INTEGER  NOT NULL DEFAULT 0,
    quota_used      INTEGER  NOT NULL DEFAULT 0,
    allowed_models  TEXT,
    allowed_ips     TEXT,
    rpm             INTEGER  NOT NULL DEFAULT 0,
    tpm             INTEGER  NOT NULL DEFAULT 0,
    expires_at      DATETIME,
    enabled         INTEGER  NOT NULL DEFAULT 1,
    last_used_at    DATETIME,
    last_used_ip    TEXT     NOT NULL DEFAULT '',
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_api_keys_key_hash ON api_keys(key_hash);
CREATE        INDEX IF NOT EXISTS idx_api_keys_user     ON api_keys(user_id);
CREATE        INDEX IF NOT EXISTS idx_api_keys_enabled  ON api_keys(enabled);

-- ============================================================
-- proxies
-- ============================================================
CREATE TABLE IF NOT EXISTS proxies (
    id            INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    scheme        TEXT     NOT NULL DEFAULT 'http',
    host          TEXT     NOT NULL,
    port          INTEGER  NOT NULL,
    username      TEXT     NOT NULL DEFAULT '',
    password_enc  TEXT     NOT NULL DEFAULT '',
    country       TEXT     NOT NULL DEFAULT '',
    isp           TEXT     NOT NULL DEFAULT '',
    health_score  INTEGER  NOT NULL DEFAULT 100,
    last_probe_at DATETIME,
    last_error    TEXT     NOT NULL DEFAULT '',
    enabled       INTEGER  NOT NULL DEFAULT 1,
    remark        TEXT     NOT NULL DEFAULT '',
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    DATETIME
);
CREATE INDEX IF NOT EXISTS idx_proxies_enabled_health ON proxies(enabled, health_score);

-- ============================================================
-- oai_accounts  (final shape: init + migration 20260418000005 + 20260423000004)
-- Note: active_email generated column NOT included — replaced by partial unique index below.
-- ============================================================
CREATE TABLE IF NOT EXISTS oai_accounts (
    id                      INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    email                   TEXT     NOT NULL,
    auth_token_enc          TEXT     NOT NULL DEFAULT '',
    refresh_token_enc       TEXT,
    session_token_enc       TEXT,
    token_expires_at        DATETIME,
    oai_session_id          TEXT     NOT NULL DEFAULT '',
    oai_device_id           TEXT     NOT NULL DEFAULT '',
    client_id               TEXT     NOT NULL DEFAULT 'app_EMoamEEZ73f0CkXaXp7hrann',
    chatgpt_account_id      TEXT     NOT NULL DEFAULT '',
    account_type            TEXT     NOT NULL DEFAULT 'codex',
    plan_type               TEXT     NOT NULL DEFAULT 'plus',
    daily_image_quota       INTEGER  NOT NULL DEFAULT 100,
    status                  TEXT     NOT NULL DEFAULT 'healthy',
    warned_at               DATETIME,
    cooldown_until          DATETIME,
    last_used_at            DATETIME,
    today_used_count        INTEGER  NOT NULL DEFAULT 0,
    today_used_date         DATE,
    last_refresh_at         DATETIME,
    last_refresh_source     TEXT     NOT NULL DEFAULT '',
    refresh_error           TEXT     NOT NULL DEFAULT '',
    image_quota_remaining   INTEGER  NOT NULL DEFAULT -1,
    image_quota_total       INTEGER  NOT NULL DEFAULT -1,
    image_quota_reset_at    DATETIME,
    image_quota_updated_at  DATETIME,
    notes                   TEXT     NOT NULL DEFAULT '',
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              DATETIME
);
-- Partial unique index: only active (non-deleted) rows must have unique email
CREATE UNIQUE INDEX IF NOT EXISTS uk_active_email        ON oai_accounts(email) WHERE deleted_at IS NULL;
CREATE        INDEX IF NOT EXISTS idx_oai_status_used    ON oai_accounts(status, last_used_at);
CREATE        INDEX IF NOT EXISTS idx_oai_token_expires  ON oai_accounts(token_expires_at);
CREATE        INDEX IF NOT EXISTS idx_oai_quota_updated  ON oai_accounts(image_quota_updated_at);

-- ============================================================
-- oai_account_cookies
-- ============================================================
CREATE TABLE IF NOT EXISTS oai_account_cookies (
    account_id       INTEGER  NOT NULL PRIMARY KEY,
    cookie_json_enc  TEXT     NOT NULL,
    updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- account_proxy_bindings
-- ============================================================
CREATE TABLE IF NOT EXISTS account_proxy_bindings (
    account_id INTEGER  NOT NULL PRIMARY KEY,
    proxy_id   INTEGER  NOT NULL,
    bound_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_apb_proxy ON account_proxy_bindings(proxy_id);

-- ============================================================
-- models
-- ============================================================
CREATE TABLE IF NOT EXISTS models (
    id                      INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    slug                    TEXT     NOT NULL,
    type                    TEXT     NOT NULL,
    upstream_model_slug     TEXT     NOT NULL,
    input_price_per_1m      INTEGER  NOT NULL DEFAULT 0,
    output_price_per_1m     INTEGER  NOT NULL DEFAULT 0,
    cache_read_price_per_1m INTEGER  NOT NULL DEFAULT 0,
    image_price_per_call    INTEGER  NOT NULL DEFAULT 0,
    description             TEXT     NOT NULL DEFAULT '',
    enabled                 INTEGER  NOT NULL DEFAULT 1,
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_models_slug ON models(slug) WHERE deleted_at IS NULL;
CREATE        INDEX IF NOT EXISTS idx_models_type ON models(type);

INSERT OR IGNORE INTO models (slug, type, upstream_model_slug, input_price_per_1m, output_price_per_1m, image_price_per_call, description) VALUES
    ('gpt-5',            'chat',  'gpt-5',            25000,  75000,       0, 'GPT-5 主力模型'),
    ('gpt-5-mini',       'chat',  'gpt-5-mini',        5000,  15000,       0, 'GPT-5 轻量'),
    ('gpt-5-codex-max',  'chat',  'gpt-5-codex-max',  50000, 150000,       0, '代码专用'),
    ('gpt-image-1',      'image', 'auto',                  0,      0, 500000, '[deprecated] 旧占位'),
    ('gpt-image-2',      'image', 'gpt-5-3',               0,      0, 500000, 'GPT-Image-2 高清生图(picture_v2)'),
    ('gpt-5-1',          'chat',  'gpt-5-1',          25000,  75000,       0, 'GPT-5.1'),
    ('gpt-5-2',          'chat',  'gpt-5-2',          25000,  75000,       0, 'GPT-5.2'),
    ('gpt-5-2-instant',  'chat',  'gpt-5-2-instant',  25000,  75000,       0, 'GPT-5.2 Instant'),
    ('gpt-5-3',          'chat',  'gpt-5-3',          25000,  75000,       0, 'GPT-5.3'),
    ('gpt-5-3-instant',  'chat',  'gpt-5-3-instant',  25000,  75000,       0, 'GPT-5.3 Instant'),
    ('gpt-5-2-thinking', 'chat',  'gpt-5-2-thinking', 40000, 120000,       0, 'GPT-5.2 Thinking'),
    ('gpt-5-4-thinking', 'chat',  'gpt-5-4-thinking', 40000, 120000,       0, 'GPT-5.4 Thinking'),
    ('gpt-5-3-mini',     'chat',  'gpt-5-3-mini',      5000,  15000,       0, 'GPT-5.3 Mini'),
    ('gpt-5-t-mini',     'chat',  'gpt-5-t-mini',      5000,  15000,       0, 'GPT-5 Thinking Mini'),
    ('gpt-5-4-t-mini',   'chat',  'gpt-5-4-t-mini',    5000,  15000,       0, 'GPT-5.4 Thinking Mini'),
    ('o3',               'chat',  'o3',               15000,  60000,       0, 'o3'),
    ('research',         'chat',  'research',         30000,  90000,       0, 'Deep Research'),
    ('agent-mode',       'chat',  'agent-mode',       30000,  90000,       0, 'Agent');

-- ============================================================
-- usage_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS usage_logs (
    id                  INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id             INTEGER  NOT NULL,
    key_id              INTEGER  NOT NULL DEFAULT 0,
    model_id            INTEGER  NOT NULL,
    account_id          INTEGER  NOT NULL DEFAULT 0,
    request_id          TEXT     NOT NULL DEFAULT '',
    type                TEXT     NOT NULL DEFAULT 'chat',
    input_tokens        INTEGER  NOT NULL DEFAULT 0,
    output_tokens       INTEGER  NOT NULL DEFAULT 0,
    cache_read_tokens   INTEGER  NOT NULL DEFAULT 0,
    cache_write_tokens  INTEGER  NOT NULL DEFAULT 0,
    image_count         INTEGER  NOT NULL DEFAULT 0,
    credit_cost         INTEGER  NOT NULL DEFAULT 0,
    duration_ms         INTEGER  NOT NULL DEFAULT 0,
    status              TEXT     NOT NULL DEFAULT 'success',
    error_code          TEXT     NOT NULL DEFAULT '',
    ip                  TEXT     NOT NULL DEFAULT '',
    ua                  TEXT     NOT NULL DEFAULT '',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_usage_logs_user_time    ON usage_logs(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_usage_logs_key_time     ON usage_logs(key_id, created_at);
CREATE INDEX IF NOT EXISTS idx_usage_logs_model_time   ON usage_logs(model_id, created_at);
CREATE INDEX IF NOT EXISTS idx_usage_logs_account_time ON usage_logs(account_id, created_at);
CREATE INDEX IF NOT EXISTS idx_usage_logs_request_id   ON usage_logs(request_id);

-- ============================================================
-- image_tasks  (+ upscale column from migration 20260423000003)
-- ============================================================
CREATE TABLE IF NOT EXISTS image_tasks (
    id               INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    task_id          TEXT     NOT NULL,
    user_id          INTEGER  NOT NULL,
    key_id           INTEGER  NOT NULL DEFAULT 0,
    model_id         INTEGER  NOT NULL,
    account_id       INTEGER  NOT NULL DEFAULT 0,
    prompt           TEXT     NOT NULL,
    n                INTEGER  NOT NULL DEFAULT 1,
    size             TEXT     NOT NULL DEFAULT '1024x1024',
    upscale          TEXT     NOT NULL DEFAULT '',
    status           TEXT     NOT NULL DEFAULT 'queued',
    conversation_id  TEXT     NOT NULL DEFAULT '',
    file_ids         TEXT,
    result_urls      TEXT,
    error            TEXT     NOT NULL DEFAULT '',
    estimated_credit INTEGER  NOT NULL DEFAULT 0,
    credit_cost      INTEGER  NOT NULL DEFAULT 0,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at       DATETIME,
    finished_at      DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_image_tasks_task_id    ON image_tasks(task_id);
CREATE        INDEX IF NOT EXISTS idx_image_tasks_user_time ON image_tasks(user_id, created_at);
CREATE        INDEX IF NOT EXISTS idx_image_tasks_status    ON image_tasks(status);

-- ============================================================
-- system_settings
-- ============================================================
CREATE TABLE IF NOT EXISTS system_settings (
    k           TEXT     NOT NULL PRIMARY KEY,
    v           TEXT,
    description TEXT     NOT NULL DEFAULT '',
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO system_settings (k, v, description) VALUES
    ('site.name',           'GPT2API',               '站点名称,展示在登录页/顶栏'),
    ('site.description',    '企业级 OpenAI 兼容网关', '登录页副标题'),
    ('site.logo_url',       '',                       '站点 Logo URL,空则显示默认图标'),
    ('site.footer',         '',                       '页脚版权/备案文本(支持纯文本)'),
    ('site.contact_email',  '',                       '站点对外联系邮箱'),
    ('auth.allow_register',        'true',  '是否开放用户自助注册'),
    ('auth.default_group_id',      '1',     '新用户默认分组 ID'),
    ('auth.signup_bonus_credits',  '0',     '新用户注册赠送积分(厘)'),
    ('limit.default_rpm',          '60',    '用户默认 RPM'),
    ('limit.default_tpm',          '60000', '用户默认 TPM'),
    ('mail.enabled_display',       'auto',  '邮件开关展示(auto/true/false)'),
    ('proxy.probe_enabled',      'true',  '代理池健康探测总开关'),
    ('proxy.probe_interval_sec', '300',   '两轮定时探测间隔(秒)'),
    ('proxy.probe_timeout_sec',  '10',    '单条代理探测超时(秒)'),
    ('proxy.probe_target_url',   '',      '探测目标 URL;留空走内置候选链'),
    ('proxy.probe_concurrency',  '8',     '并发数(1~64)'),
    ('account.refresh_enabled',          'true',                          '账号 AT 自动刷新总开关'),
    ('account.refresh_interval_sec',     '120',                           '后台扫描间隔(秒)'),
    ('account.refresh_ahead_sec',        '900',                           '距过期多少秒内触发预刷新'),
    ('account.refresh_concurrency',      '4',                             '同时刷新账号数(1~32)'),
    ('account.quota_probe_enabled',      'true',                          '账号图片额度自动探测总开关'),
    ('account.quota_probe_interval_sec', '18000',                         '额度探测最小间隔(秒);默认 5h'),
    ('account.default_client_id',        'app_EMoamEEZ73f0CkXaXp7hrann', '导入账号时默认 client_id');

-- ============================================================
-- admin_audit_logs
-- ============================================================
CREATE TABLE IF NOT EXISTS admin_audit_logs (
    id           INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    actor_id     INTEGER  NOT NULL DEFAULT 0,
    actor_email  TEXT     NOT NULL DEFAULT '',
    action       TEXT     NOT NULL,
    method       TEXT     NOT NULL,
    path         TEXT     NOT NULL,
    status_code  INTEGER  NOT NULL DEFAULT 0,
    ip           TEXT     NOT NULL DEFAULT '',
    ua           TEXT     NOT NULL DEFAULT '',
    target       TEXT     NOT NULL DEFAULT '',
    meta         TEXT,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_audit_actor_time  ON admin_audit_logs(actor_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_action_time ON admin_audit_logs(action, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_created_at  ON admin_audit_logs(created_at);


