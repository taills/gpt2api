package user

import (
"context"
"errors"
"strconv"
"strings"

"github.com/gin-gonic/gin"

"github.com/432539/gpt2api/internal/audit"
"github.com/432539/gpt2api/internal/middleware"
"github.com/432539/gpt2api/pkg/resp"
)

// PasswordService 由 auth 包实现。通过接口解耦,避免 user<->auth 循环依赖。
type PasswordService interface {
HashPassword(plain string) (string, error)
VerifyPassword(ctx context.Context, userID uint64, password string) error
}

// AdminHandler 管理员视角下的用户管理接口。
type AdminHandler struct {
dao      *DAO
auth     PasswordService
auditDAO *audit.DAO
}

// NewAdminHandler 构造。
func NewAdminHandler(dao *DAO, authSvc PasswordService, auditDAO *audit.DAO) *AdminHandler {
return &AdminHandler{dao: dao, auth: authSvc, auditDAO: auditDAO}
}

// List GET /api/admin/users
func (h *AdminHandler) List(c *gin.Context) {
limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

items, total, err := h.dao.ListPage(c.Request.Context(), ListFilter{
Keyword: c.Query("q"),
Role:    c.Query("role"),
Status:  c.Query("status"),
}, limit, offset)
if err != nil {
resp.Internal(c, err.Error())
return
}
resp.OK(c, gin.H{"items": items, "total": total, "limit": limit, "offset": offset})
}

// createReq POST /api/admin/users 的 body。
type createReq struct {
Email    string `json:"email" binding:"required,email"`
Password string `json:"password" binding:"required,min=6"`
Nickname string `json:"nickname"`
Role     string `json:"role"`   // user | admin,默认 user
Status   string `json:"status"` // active | banned,默认 active
}

// Create POST /api/admin/users
func (h *AdminHandler) Create(c *gin.Context) {
var req createReq
if err := c.ShouldBindJSON(&req); err != nil {
resp.BadRequest(c, err.Error())
return
}

email := strings.ToLower(strings.TrimSpace(req.Email))
if email == "" {
resp.BadRequest(c, "email required")
return
}
role := strings.ToLower(strings.TrimSpace(req.Role))
if role == "" {
role = "user"
}
if role != "user" && role != "admin" {
resp.BadRequest(c, "role must be user or admin")
return
}
status := strings.ToLower(strings.TrimSpace(req.Status))
if status == "" {
status = "active"
}
if status != "active" && status != "banned" {
resp.BadRequest(c, "status must be active or banned")
return
}

ctx := c.Request.Context()
if n, err := h.dao.CountByEmail(ctx, email); err != nil {
resp.Internal(c, err.Error())
return
} else if n > 0 {
resp.BadRequest(c, "邮箱已存在")
return
}

hash, err := h.auth.HashPassword(req.Password)
if err != nil {
resp.BadRequest(c, err.Error())
return
}

u := &User{
Email:        email,
PasswordHash: hash,
Nickname:     strings.TrimSpace(req.Nickname),
GroupID:      1,
Role:         role,
Status:       status,
}
id, err := h.dao.Create(ctx, u)
if err != nil {
resp.Internal(c, err.Error())
return
}
u.ID = id

audit.Record(c, h.auditDAO, "users.create", strconv.FormatUint(id, 10),
gin.H{"email": email, "role": role, "status": status})
resp.OK(c, gin.H{
"id":     id,
"email":  u.Email,
"role":   u.Role,
"status": u.Status,
})
}

// Get GET /api/admin/users/:id
func (h *AdminHandler) Get(c *gin.Context) {
id, ok := parseID(c)
if !ok {
return
}
u, err := h.dao.GetByID(c.Request.Context(), id)
if err != nil {
if errors.Is(err, ErrNotFound) {
resp.NotFound(c, "user not found")
return
}
resp.Internal(c, err.Error())
return
}
resp.OK(c, u)
}

// updateReq PATCH /api/admin/users/:id 的 body。全部字段可选。
type updateReq struct {
Nickname *string `json:"nickname,omitempty"`
Role     *string `json:"role,omitempty"`
Status   *string `json:"status,omitempty"`
}

// Update PATCH /api/admin/users/:id
func (h *AdminHandler) Update(c *gin.Context) {
id, ok := parseID(c)
if !ok {
return
}
var req updateReq
if err := c.ShouldBindJSON(&req); err != nil {
resp.BadRequest(c, err.Error())
return
}
// 防自我锁死:不允许把自己从 admin 改成 user(避免最后一个 admin 流失)
actor := middleware.UserID(c)
if req.Role != nil && *req.Role != "admin" && id == actor {
resp.BadRequest(c, "cannot downgrade your own admin role")
return
}

n, err := h.dao.Update(c.Request.Context(), id, UpdatePatch{
Nickname: req.Nickname,
Role:     req.Role,
Status:   req.Status,
})
if err != nil {
resp.BadRequest(c, err.Error())
return
}
if n == 0 {
resp.NotFound(c, "user not found")
return
}
audit.Record(c, h.auditDAO, "users.update", strconv.FormatUint(id, 10), req)
resp.OK(c, gin.H{"updated": n})
}

// resetPwdReq 重置密码的 body。
type resetPwdReq struct {
NewPassword   string `json:"new_password" binding:"required,min=6"`
AdminPassword string `json:"admin_password" binding:"required"`
}

// ResetPassword POST /api/admin/users/:id/reset-password
func (h *AdminHandler) ResetPassword(c *gin.Context) {
id, ok := parseID(c)
if !ok {
return
}
var req resetPwdReq
if err := c.ShouldBindJSON(&req); err != nil {
resp.BadRequest(c, err.Error())
return
}
actor := middleware.UserID(c)
if err := h.auth.VerifyPassword(c.Request.Context(), actor, req.AdminPassword); err != nil {
resp.Forbidden(c, "admin password mismatch")
return
}
hash, err := h.auth.HashPassword(req.NewPassword)
if err != nil {
resp.BadRequest(c, err.Error())
return
}
if err := h.dao.ResetPassword(c.Request.Context(), id, hash); err != nil {
if errors.Is(err, ErrNotFound) {
resp.NotFound(c, "user not found")
return
}
resp.Internal(c, err.Error())
return
}
audit.Record(c, h.auditDAO, "users.reset_password", strconv.FormatUint(id, 10), nil)
resp.OK(c, gin.H{"ok": true})
}

// Delete DELETE /api/admin/users/:id (软删除)
func (h *AdminHandler) Delete(c *gin.Context) {
id, ok := parseID(c)
if !ok {
return
}
actor := middleware.UserID(c)
if id == actor {
resp.BadRequest(c, "cannot delete yourself")
return
}
if err := confirmAdmin(c, h.auth); err != nil {
resp.Forbidden(c, err.Error())
return
}
if err := h.dao.SoftDelete(c.Request.Context(), id); err != nil {
if errors.Is(err, ErrNotFound) {
resp.NotFound(c, "user not found")
return
}
resp.Internal(c, err.Error())
return
}
audit.Record(c, h.auditDAO, "users.delete", strconv.FormatUint(id, 10), nil)
resp.OK(c, gin.H{"deleted": id})
}

// ---- helpers ----

func parseID(c *gin.Context) (uint64, bool) {
id, err := strconv.ParseUint(c.Param("id"), 10, 64)
if err != nil || id == 0 {
resp.BadRequest(c, "invalid id")
return 0, false
}
return id, true
}

// confirmAdmin 从 header X-Admin-Confirm 里拿密码做二次校验。
func confirmAdmin(c *gin.Context, authSvc PasswordService) error {
pwd := c.GetHeader("X-Admin-Confirm")
if pwd == "" {
pwd = c.PostForm("admin_password")
}
if pwd == "" {
return errors.New("X-Admin-Confirm header required for this destructive operation")
}
actor := middleware.UserID(c)
if actor == 0 {
return errors.New("not authenticated")
}
if err := authSvc.VerifyPassword(c.Request.Context(), actor, pwd); err != nil {
return errors.New("admin password mismatch")
}
return nil
}
