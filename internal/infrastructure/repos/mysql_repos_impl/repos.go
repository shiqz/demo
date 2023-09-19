// Package mysql_repos_impl 仓库MySQL实现
package mysql_repos_impl

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var dialect = goqu.Dialect("mysql")
