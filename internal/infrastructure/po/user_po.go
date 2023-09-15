// Package po 定义相关数据库结构映射
package po

// User 用户表映射
type User struct {
	UserID     uint   `db:"user_id"`
	Username   string `db:"username"`
	Passwd     string `db:"passwd"`
	Salt       string `db:"salt"`
	Nickname   string `db:"nickname"`
	Gender     uint   `db:"gender"`
	Status     uint   `db:"status"`
	CreateTime int64  `db:"create_time"`
}
