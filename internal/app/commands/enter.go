package commands

import (
	"example/internal/domain/types"
	"fmt"
	"github.com/modood/table"
	"github.com/spf13/cobra"
)

// Version 打印版本号
func Version(_ *cobra.Command, _ []string) {
	commander.version()
}

// RunAPIServer 启动API服务
func RunAPIServer(_ *cobra.Command, _ []string) {
	commander.runAPIServer()
}

// AdminCreate 创建管理员
func AdminCreate(cmd *cobra.Command, _ []string) error {
	email := cmd.Flag("email").Value.String()
	pass := cmd.Flag("pass").Value.String()
	role := cmd.Flag("role").Value.String()
	return commander.account().Create(email, pass, role)
}

// ShowAdmins 创建管理员
func ShowAdmins(_ *cobra.Command, _ []string) error {
	return commander.account().ShowAccounts()
}

// UpdateAdminRole 修改管理员角色
func UpdateAdminRole(cmd *cobra.Command, _ []string) error {
	return commander.account().UpdateRole(
		cmd.Flag("email").Value.String(),
		cmd.Flag("role").Value.String(),
	)
}

// UpdateAdminPass 修改管理员密码
func UpdateAdminPass(cmd *cobra.Command, _ []string) error {
	return commander.account().UpdatePass(
		cmd.Flag("email").Value.String(),
		cmd.Flag("pass").Value.String(),
	)
}

// ShowAllRoles 查看系统所有角色
func ShowAllRoles(_ *cobra.Command, _ []string) error {
	roles := types.AllRoles
	type item struct {
		Role string
		Name string
	}
	var list []item
	for userRole, s := range roles {
		list = append(list, item{
			Role: userRole.String(),
			Name: s,
		})
	}
	fmt.Println(table.Table(list))
	return nil
}

// ShowAdminRoles 查看系统账号拥有角色
func ShowAdminRoles(cmd *cobra.Command, _ []string) error {
	return commander.account().ShowAccountRole(cmd.Flag("email").Value.String())
}

// ShowAllPerms 查看系统所有路由权限
func ShowAllPerms(_ *cobra.Command, _ []string) error {
	data := types.AllPerms
	type item struct {
		Method string
		Route  string
		Name   string
	}
	var list []item
	for _, info := range data {
		list = append(list, item{
			Method: info.Method,
			Route:  info.Path,
			Name:   info.Name,
		})
	}
	fmt.Println(table.Table(list))
	return nil
}

// ShowRolePerms 显示角色拥有路由权限
func ShowRolePerms(cmd *cobra.Command, _ []string) error {
	role := cmd.Flag("role").Value.String()
	roles, err := types.ParseRoles(role, false)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return fmt.Errorf("请输入要查看的角色名称")
	}
	roles.ShowPerms()
	return nil
}

// ShowAdminPerms 显示账户拥有路由权限
func ShowAdminPerms(cmd *cobra.Command, _ []string) error {
	return commander.account().ShowAccountPerms(cmd.Flag("email").Value.String())
}

// UpdateUserSession 更新会话过期时间
func UpdateUserSession(cmd *cobra.Command, _ []string) error {
	return commander.user().UpdateSession(
		cmd.Flag("id").Value.String(),
		cmd.Flag("time").Value.String(),
	)
}

// RemoveUserSession 删除会话
func RemoveUserSession(cmd *cobra.Command, _ []string) error {
	return commander.user().RemoveSession(
		cmd.Flag("id").Value.String(),
	)
}
