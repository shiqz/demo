// Package mysqlreposimpl 仓库MySQL实现
package mysqlreposimpl

import (
	"github.com/doug-martin/goqu/v9"
	// MySQL驱动
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var dialect = goqu.Dialect("mysql")
