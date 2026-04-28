package server

import (
"io/fs"
"net/http"
"os"
"path/filepath"
"strings"

"github.com/gin-gonic/gin"
)

// mountSPA 把前端 Vite 产物挂到 `/` 上,并实现 SPA 回退(deep link 刷新)。
//
// 优先级:
//  1. 环境变量 GPT2API_WEB_DIR 指定的磁盘目录(开发 / 调试覆盖)
//  2. 磁盘候选路径:./web/dist
//  3. 编译期嵌入的 embed.FS(internal/server/web/)
//  4. 三者均不可用 → 退化为纯 API 服务
//
// 注意:
//   - 只有 GET/HEAD 的 NoRoute 请求才会被 fallback 到 index.html。其它方法保持 404。
//   - 明确排除 /api/、/v1/、/healthz、/readyz 等 API 前缀,避免 404 被 index.html 掩盖。
func mountSPA(r *gin.Engine) bool {
// 优先尝试磁盘路径(便于开发期热更前端)
if dir := resolveWebDir(); dir != "" {
return mountSPAFromDisk(r, dir)
}
// 回退到 embed.FS
return mountSPAFromEmbed(r)
}

// mountSPAFromDisk 从磁盘目录挂载 SPA。
func mountSPAFromDisk(r *gin.Engine, dir string) bool {
indexPath := filepath.Join(dir, "index.html")
if _, err := os.Stat(indexPath); err != nil {
return false
}

entries, _ := os.ReadDir(dir)
for _, ent := range entries {
name := ent.Name()
if name == "" || name == "." || name == ".." {
continue
}
full := filepath.Join(dir, name)
if ent.IsDir() {
r.Static("/"+name, full)
} else if name != "index.html" {
r.StaticFile("/"+name, full)
}
}

r.GET("/", func(c *gin.Context) { c.File(indexPath) })
r.NoRoute(spaFallbackDisk(indexPath))
return true
}

// mountSPAFromEmbed 从编译期嵌入的 FS 挂载 SPA。
// embed.FS 路径形如 "web/index.html"(包含 "web/" 前缀),
// 所以用 fs.Sub 剥掉前缀后再使用。
func mountSPAFromEmbed(r *gin.Engine) bool {
sub, err := fs.Sub(webFS, "web")
if err != nil {
return false
}
// 探测 index.html 是否真的存在(空目录时 embed 仍然有效但无文件)
if _, err := sub.Open("index.html"); err != nil {
return false
}

httpFS := http.FS(sub)
fileServer := http.FileServer(httpFS)

// 注册根路径
r.GET("/", func(c *gin.Context) {
c.FileFromFS("index.html", httpFS)
})

// 注册已知的静态目录/文件前缀,让 gin 走 FileServer
entries, _ := fs.ReadDir(sub, ".")
for _, ent := range entries {
name := ent.Name()
if name == "" || name == "index.html" {
continue
}
captured := name
if ent.IsDir() {
r.GET("/"+captured+"/*filepath", func(c *gin.Context) {
c.Request.URL.Path = "/" + captured + c.Param("filepath")
fileServer.ServeHTTP(c.Writer, c.Request)
})
} else {
r.GET("/"+captured, func(c *gin.Context) {
c.Request.URL.Path = "/" + captured
fileServer.ServeHTTP(c.Writer, c.Request)
})
}
}

r.NoRoute(spaFallbackEmbed(httpFS))
return true
}

// spaFallbackDisk 返回从磁盘提供 index.html 的 NoRoute handler。
func spaFallbackDisk(indexPath string) gin.HandlerFunc {
return func(c *gin.Context) {
if !isSPARequest(c) {
c.Status(http.StatusNotFound)
return
}
c.File(indexPath)
}
}

// spaFallbackEmbed 返回从 embed.FS 提供 index.html 的 NoRoute handler。
func spaFallbackEmbed(httpFS http.FileSystem) gin.HandlerFunc {
return func(c *gin.Context) {
if !isSPARequest(c) {
c.Status(http.StatusNotFound)
return
}
c.FileFromFS("index.html", httpFS)
}
}

// isSPARequest 判断是否应该返回 index.html(GET/HEAD 且不命中 API 前缀)。
func isSPARequest(c *gin.Context) bool {
if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
return false
}
p := c.Request.URL.Path
for _, prefix := range apiPrefixes {
if strings.HasPrefix(p, prefix) {
return false
}
}
return true
}

// API 前缀白名单:命中这里的请求不做 SPA fallback。
var apiPrefixes = []string{
"/api/",
"/v1/",
"/p/",
"/healthz",
"/readyz",
}

// resolveWebDir 查找磁盘上的 web/dist 目录。
// 返回空字符串表示不使用磁盘路径(回退到 embed)。
func resolveWebDir() string {
if d := os.Getenv("GPT2API_WEB_DIR"); d != "" {
if isDir(d) {
return d
}
}
for _, d := range []string{"./web/dist"} {
if isDir(d) {
abs, _ := filepath.Abs(d)
return abs
}
}
return ""
}

func isDir(p string) bool {
st, err := os.Stat(p)
if err != nil {
return false
}
return st.IsDir()
}
