// local_cache.go — 图片本地磁盘缓存。
//
// 写入路径:  <baseDir>/YYYYMMDD/<md5_of_url>
// Content-Type:  <baseDir>/YYYYMMDD/<md5_of_url>.ct  (sidecar)
//
// 读取策略:  用 glob 扫描 <baseDir>/*/<md5_of_url>,
//
//	无需提前知道写入日期,天然幂等。
package image

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// fetchImageHTTP 用标准 http.Client 下载公开签名 URL(files.oaiusercontent.com)。
// 仅用于 runner 内的预热缓存;超时 60s;最大 32MB。
func fetchImageHTTP(url string) (data []byte, contentType string, err error) {
	const maxBytes = 32 * 1024 * 1024
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url) //nolint:gosec // URL comes from chatgpt.com signed URL
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("unexpected status: " + resp.Status)
	}
	data, err = io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	if err != nil {
		return nil, "", err
	}
	ct := resp.Header.Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(data)
	}
	return data, ct, nil
}

// imageCacheFilename 根据原始 URL 计算缓存文件名(MD5 hex)。
func imageCacheFilename(rawURL string) string {
	sum := md5.Sum([]byte(rawURL))
	return hex.EncodeToString(sum[:])
}

// ImageCachePath 返回应写入的缓存路径(basDir/YYYYMMDD/<md5>)。
// 写入时按"今天"决定目录。
func ImageCachePath(baseDir, rawURL string) string {
	date := time.Now().Format("20060102")
	name := imageCacheFilename(rawURL)
	return filepath.Join(baseDir, date, name)
}

// ReadImageCache 从磁盘缓存读取图片。
// 返回 (data, contentType, true) 命中;(nil, "", false) 未命中。
func ReadImageCache(baseDir, rawURL string) ([]byte, string, bool) {
	name := imageCacheFilename(rawURL)
	// glob 扫描所有日期目录
	pattern := filepath.Join(baseDir, "*", name)
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return nil, "", false
	}
	// 使用第一个命中
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, "", false
	}
	// 读取 sidecar content-type
	ct := "image/png"
	ctPath := matches[0] + ".ct"
	if ctBytes, err := os.ReadFile(ctPath); err == nil && len(ctBytes) > 0 {
		ct = string(ctBytes)
	} else {
		// sidecar 不存在时嗅探
		sniffed := http.DetectContentType(data)
		if sniffed != "" {
			ct = sniffed
		}
	}
	return data, ct, true
}

// WriteImageCache 将图片写入本地缓存。
// 目录按需创建;写入失败不影响主流程,调用方可忽略错误。
func WriteImageCache(baseDir, rawURL string, data []byte, contentType string) error {
	if len(data) == 0 {
		return errors.New("empty data")
	}
	date := time.Now().Format("20060102")
	dir := filepath.Join(baseDir, date)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	name := imageCacheFilename(rawURL)
	imgPath := filepath.Join(dir, name)
	// 先写临时文件再原子 rename,防止写一半被读到
	tmp := imgPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, imgPath); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	// 写 content-type sidecar
	ctPath := imgPath + ".ct"
	_ = os.WriteFile(ctPath, []byte(contentType), 0o644)
	return nil
}

// ImageCacheExists 快速检查缓存是否存在(不读取内容)。
func ImageCacheExists(baseDir, rawURL string) bool {
	name := imageCacheFilename(rawURL)
	pattern := filepath.Join(baseDir, "*", name)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false
	}
	for _, m := range matches {
		if _, err := os.Stat(m); !errors.Is(err, fs.ErrNotExist) {
			return true
		}
	}
	return false
}
