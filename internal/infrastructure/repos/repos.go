// Package repos 仓库实现
package repos

import (
	log "github.com/sirupsen/logrus"
)

// format 转义SQL
func format(sql string) string {
	//sql = strings.ReplaceAll(sql, "\"", "`")
	log.Tracef("[DB]%s", sql)
	return sql
}
