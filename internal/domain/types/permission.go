// Package types 应用路由权限相关
package types

import (
	"fmt"
	"github.com/modood/table"
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

// Info 路由信息
type Info struct {
	Route  string
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
	AllPerms = map[Permission]Info{
		Users:          {Route: "/api/admin/users", Method: http.MethodGet, Name: "获取用户列表接口"},
		UsersResetPass: {Route: "/api/admin/users/passwd", Method: http.MethodPatch, Name: "重置用户密码"},
		UsersSetStatus: {Route: "/api/admin/users/status", Method: http.MethodPatch, Name: "更改用户状态"},
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
			Route:  AllPerms[key].Route,
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
	if len(str) == 0 {
		str = append(str, UserRoleDefault.String())
	}
	return strings.Join(str, ",")
}

// ParseRoles 将字符串解析为角色集合
func ParseRoles(str string, skipInvalid bool) (Roles, error) {
	var roles Roles
	if str == "" {
		return roles, nil
	}
	for _, roleStr := range strings.Split(str, ",") {
		role := UserRole(roleStr)
		if !role.Valid() {
			if skipInvalid {
				continue
			}
			return nil, fmt.Errorf("invalid role: %s", roleStr)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// HasPermission 路由权限校验
func HasPermission(raw, uri, method string) bool {
	// 判断路由是否需要校验
	var per Permission
	for id, info := range AllPerms {
		if info.Route == uri && info.Method == method {
			per = id
			break
		}
	}
	// 路由不需要做权限校验
	if per == 0 {
		return true
	}
	roles, _ := ParseRoles(raw, true)
	for _, role := range roles {
		for _, permission := range role.GetPermissions() {
			if permission == per {
				return true
			}
		}
	}
	return false
}
