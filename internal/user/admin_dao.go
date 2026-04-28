package user

import (
	"context"
	"fmt"
	"strings"
)

// ListFilter admin 查询用户列表的过滤条件。
type ListFilter struct {
	Keyword string // 模糊匹配 email/nickname
	Role    string // "" / "user" / "admin"
	Status  string // "" / "active" / "banned"
}

// ListPage admin 分页列出用户。返回结果带最大值兜底,limit 超过 500 强制 500。
func (d *DAO) ListPage(ctx context.Context, f ListFilter, limit, offset int) ([]User, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}

	if f.Keyword != "" {
		where = append(where, "(email LIKE ? OR nickname LIKE ?)")
		kw := "%" + f.Keyword + "%"
		args = append(args, kw, kw)
	}
	if f.Role != "" {
		where = append(where, "role = ?")
		args = append(args, f.Role)
	}
	if f.Status != "" {
		where = append(where, "status = ?")
		args = append(args, f.Status)
	}

	whereSQL := "WHERE " + strings.Join(where, " AND ")

	var total int64
	if err := d.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM users `+whereSQL, args...); err != nil {
		return nil, 0, err
	}

	var items []User
	argsWithLimit := append(args, limit, offset)
	if err := d.db.SelectContext(ctx, &items,
		`SELECT id, email, password_hash, nickname, group_id, role, status,
                credit_balance, credit_frozen, version, last_login_at, last_login_ip,
                created_at, updated_at, deleted_at
           FROM users `+whereSQL+`
          ORDER BY id DESC
          LIMIT ? OFFSET ?`, argsWithLimit...); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdatePatch 用于 admin 更新用户基础字段(允许仅传部分)。
type UpdatePatch struct {
	Nickname *string
	Role     *string
	Status   *string
}

// Update 执行 UpdatePatch;每个非 nil 字段才会被更新。返回受影响行数。
func (d *DAO) Update(ctx context.Context, id uint64, p UpdatePatch) (int64, error) {
	sets := []string{}
	args := []interface{}{}

	if p.Nickname != nil {
		sets = append(sets, "nickname = ?")
		args = append(args, *p.Nickname)
	}
	if p.Role != nil {
		if *p.Role != "user" && *p.Role != "admin" {
			return 0, fmt.Errorf("invalid role: %s", *p.Role)
		}
		sets = append(sets, "role = ?")
		args = append(args, *p.Role)
	}
	if p.Status != nil {
		if *p.Status != "active" && *p.Status != "banned" {
			return 0, fmt.Errorf("invalid status: %s", *p.Status)
		}
		sets = append(sets, "status = ?")
		args = append(args, *p.Status)
	}
	if len(sets) == 0 {
		return 0, nil
	}
	sets = append(sets, "version = version + 1")
	q := "UPDATE users SET " + strings.Join(sets, ", ") + " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, id)

	res, err := d.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

// ResetPassword 覆盖 password_hash。hash 由上层用 bcrypt 生成。
func (d *DAO) ResetPassword(ctx context.Context, id uint64, hash string) error {
	res, err := d.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ?, version = version + 1
          WHERE id = ? AND deleted_at IS NULL`, hash, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDelete 将用户标记为删除(并不物理删除,也不回收其 api keys/usage 等)。
// 被删除的用户无法登录,api key 按策略可单独吊销。
func (d *DAO) SoftDelete(ctx context.Context, id uint64) error {
	res, err := d.db.ExecContext(ctx,
		`UPDATE users SET deleted_at = CURRENT_TIMESTAMP, status = 'banned', version = version + 1
          WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
