package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var schemas = []string{`
CREATE TABLE IF NOT EXISTS users (
    user_id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username varchar(32) NOT NULL DEFAULT '' COMMENT '登录用户名',
    passwd varchar(32) NOT NULL DEFAULT '' COMMENT '登录密码',
    salt char(8) NOT NULL DEFAULT '' COMMENT '密码加密盐',
    gender tinyint(1) NOT NULL DEFAULT 0 COMMENT '性别；0: 未知; 1: 男; 2: 女',
    create_time bigint NOT NULL DEFAULT 0 COMMENT '记录创建时间, 使用unixtimestamp',
    nickname varchar(64) NOT NULL DEFAULT '' COMMENT '显示昵称',
    status tinyint(1) NOT NULL DEFAULT 1 COMMENT '用户状态；1： 启用；2：禁用',
    CONSTRAINT ix_users_username UNIQUE (username)
) Engine=Innodb DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT = 1000000 COMMENT '用户表'
`,
	`CREATE TABLE IF NOT EXISTS admins (
    admin_id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email varchar(64) NOT NULL DEFAULT '' COMMENT '登录邮箱',
    passwd varchar(64) NOT NULL DEFAULT '' COMMENT '登录密码',
    roles varchar(127) NOT NULL DEFAULT '' COMMENT '用户角色',
    create_time bigint(11) NOT NULL,
    CONSTRAINT uk_admins_emails UNIQUE (email)
) Engine=Innodb DEFAULT CHARSET=utf8mb4 COMMENT '系统管理员';`,
}

// 迁移
func migrate(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, schema := range schemas {
		if _, err = tx.Exec(schema); err != nil {
			return errors.WithStack(err)
		}
	}
	return tx.Commit()
}
