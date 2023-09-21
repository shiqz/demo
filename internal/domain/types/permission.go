// Package types 应用路由权限相关
package types

import (
	"fmt"
	"github.com/modood/table"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
)

// UserRole 角色类型
type UserRole string

// Roles 角色类型集合
type Roles []UserRole

// Permission 路由权限类型
type Permission int

// Permissions 路由权限集合类型
type Permissions []Permission

// Route 路由信息
type Route struct {
	Path   string
	Name   string
	Method string
}

// 路由表
const (
	Users          Permission = iota + 1 // 用户列表接口
	UsersResetPass                       // 重置用户密码
	UsersSetStatus                       // 更改用户状态
)

// 角色表
const (
	UserRoleAll     UserRole = "*"
	UserRoleDefault UserRole = "default"
	UserRoleManager UserRole = "manager"
	UserRoleTest    UserRole = "test"
)

var (
	// AllRoles 系统所有角色信息
	AllRoles = map[UserRole]string{
		UserRoleAll:     "所有权限",
		UserRoleManager: "超级管理员",
		UserRoleTest:    "测试角色",
		UserRoleDefault: "默认角色（无任何权限）",
	}
	// AllPerms 所有路由权限信息
	AllPerms = map[Permission]Route{
		Users:          {Path: "/api/admin/users", Method: http.MethodGet, Name: "获取用户列表接口"},
		UsersResetPass: {Path: "/api/admin/users/passwd", Method: http.MethodPatch, Name: "重置用户密码"},
		UsersSetStatus: {Path: "/api/admin/users/status", Method: http.MethodPatch, Name: "更改用户状态"},
	}
	// RolePermissionsMap 角色路由权限集合
	RolePermissionsMap = map[UserRole]Permissions{
		UserRoleManager: {
			Users, UsersResetPass, UsersSetStatus,
		},
		UserRoleTest: {
			Users,
		},
	}
)

func init() {
	// 初始化超集角色
	for perm := range AllPerms {
		RolePermissionsMap[UserRoleAll] = append(RolePermissionsMap[UserRoleAll], perm)
	}
}

// Valid 校验角色有效性
func (r UserRole) Valid() bool {
	_, ok := AllRoles[r]
	return ok
}

// GetPermissions 获取角色路由权限
func (r UserRole) GetPermissions() Permissions {
	return RolePermissionsMap[r]
}

// GetPermissionsTable 获取角色路由权限表
func (r UserRole) GetPermissionsTable() string {
	type field struct {
		ID     int
		Name   string
		Method string
		Route  string
	}
	var keys Permissions
	for id := range AllPerms {
		keys = append(keys, id)
	}
	sort.Slice(keys, func(i, j int) bool {
		return int(keys[i]) < int(keys[j])
	})
	var data []field
	for _, key := range keys {
		data = append(data, field{
			ID:     int(key),
			Route:  AllPerms[key].Path,
			Name:   AllPerms[key].Name,
			Method: AllPerms[key].Method,
		})
	}
	return table.Table(data)
}

// Name 获取角色描述名称
func (r UserRole) Name() string {
	return AllRoles[r]
}

// Name 获取角色描述名称
func (r UserRole) String() string {
	return string(r)
}

// String 角色集合字符串序列化
func (rs Roles) String() string {
	var str []string
	for _, r := range rs {
		str = append(str, string(r))
	}
	sort.Strings(str)
	if len(str) == 0 {
		str = append(str, UserRoleDefault.String())
	}
	return strings.Join(str, ",")
}

// Eq 判断是否不相等
func (rs Roles) Eq(cr Roles) bool {
	return rs.String() == cr.String()
}

// HasPermission 验证路由权限
func (rs Roles) HasPermission(route Route) bool {
	for _, role := range rs {
		for _, perm := range role.GetPermissions() {
			owner := perm.GetRouteInfo()
			if owner.Method == route.Method && owner.Path == route.Path {
				return true
			}
		}
	}
	return false
}

// ShowPerms 显示路由权限
func (rs Roles) ShowPerms() {
	allPerms := make(map[Permission]Route)
	for _, userRole := range rs {
		for _, permission := range userRole.GetPermissions() {
			if _, ok := allPerms[permission]; !ok {
				allPerms[permission] = permission.GetRouteInfo()
			}
		}
	}
	if len(allPerms) == 0 {
		log.Warnf("角色：%s 无任何路由权限", rs)
		return
	}
	type item struct {
		Method string
		Route  string
		Name   string
	}
	var list []item
	for _, info := range allPerms {
		list = append(list, item{
			Method: info.Method,
			Route:  info.Path,
			Name:   info.Name,
		})
	}
	fmt.Println(table.Table(list))
}

// GetRouteInfo 获取路由权限信息
func (p Permission) GetRouteInfo() Route {
	return AllPerms[p]
}

// ParseRoles 将字符串解析为角色集合
func ParseRoles(str string, skipInvalid bool) (Roles, error) {
	var roles Roles
	if str == "" {
		return roles, nil
	}
	for _, roleStr := range strings.Split(str, ",") {
		role := UserRole(strings.TrimSpace(roleStr))
		if role == "" {
			continue
		}
		if !role.Valid() {
			if skipInvalid {
				continue
			}
			return nil, fmt.Errorf("invalid role: %s", role)
		}
		roles = append(roles, role)
	}
	return roles, nil
}
