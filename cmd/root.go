package cmd

import (
	hc "example/internal/app/commands"
	_ "example/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// command 指令类型
type command string

// 命令集合（子命令）
type commands map[command]*cobra.Command

// 管理系统账号命令
const (
	adminCreate     command = "create"           // 创建系统账号
	adminList       command = "list"             // 显示系统已创建系统账号
	adminUpdateRole command = "update-role"      // 更新账号角色
	adminUpdatePass command = "update-pass"      // 更新账号密码
	showAllRoles    command = "show-roles"       // 查看系统所有角色
	showAdminRoles  command = "show-admin-role"  // 查看系统账号拥有角色
	showAllPerms    command = "show-perms"       // 查看系统所有路由权限
	showRolePerms   command = "show-role-perms"  // 显示角色拥有路由权限
	showAdminPerms  command = "show-admin-perms" // 显示管理员拥有路由权限
)

// 管理系统用户命令
const (
	setUserSession command = "session-update" // 更新用户会话
	remUserSession command = "session-remove" // 删除用户会话
)

var (
	// 主命令
	root = &cobra.Command{
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableDescriptions: true,
			DisableNoDescFlag:   true,
			HiddenDefaultCmd:    true,
		},
		SilenceUsage: true,
	}
	version      = &cobra.Command{Use: "version", Short: "显示当前脚本版本号", Run: hc.Version}
	runAPIServer = &cobra.Command{Use: "run", Short: "运行API服务", Run: hc.RunAPIServer}
	// 管理系统账号
	admin            = &cobra.Command{Use: "admin", Short: "管理系统账号"}
	adminSubCommands = commands{
		adminCreate:     {Use: "create", Short: "创建系统账号", RunE: hc.AdminCreate},
		adminList:       {Use: "list", Short: "显示系统已创建系统账号", RunE: hc.ShowAdmins},
		adminUpdateRole: {Use: "update-role", Short: "更新账号角色", RunE: hc.UpdateAdminRole},
		adminUpdatePass: {Use: "update-pass", Short: "更新账号密码", RunE: hc.UpdateAdminPass},
		showAllRoles:    {Use: "show-roles", Short: "查看系统所有角色", RunE: hc.ShowAllRoles},
		showAdminRoles:  {Use: "show-admin-role", Short: "查看系统账号拥有角色", RunE: hc.ShowAdminRoles},
		showAllPerms:    {Use: "show-perms", Short: "查看系统所有路由权限", RunE: hc.ShowAllPerms},
		showRolePerms:   {Use: "show-role-perms", Short: "显示角色拥有路由权限", RunE: hc.ShowRolePerms},
		showAdminPerms:  {Use: "show-admin-perms", Short: "显示管理员拥有路由权限", RunE: hc.ShowAdminPerms},
	}
	// 管理用户
	user            = &cobra.Command{Use: "user", Short: "管理系统用户"}
	userSubCommands = commands{
		setUserSession: {Use: "session-update", Short: "更新会话过期时间", RunE: hc.UpdateUserSession},
		remUserSession: {Use: "session-remove", Short: "删除用户会话", RunE: hc.RemoveUserSession},
	}
)

var (
	cfgFile string
)

// 获取命令集合
func (cmds commands) all() []*cobra.Command {
	var list []*cobra.Command
	for _, cmd := range cmds {
		list = append(list, cmd)
	}
	return list
}

func init() {
	// 初始化
	cobra.OnInitialize(func() {
		hc.Init(cfgFile)
	})
	// 基本参数
	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "运行配置文件")
	// 管理系统账号相关参数
	admin.PersistentFlags().StringP("email", "e", "", "邮箱账号")
	// 创建账号参数
	adminSubCommands[adminCreate].Flags().StringP("pass", "p", "", "登录密码")
	adminSubCommands[adminCreate].Flags().StringP("role", "r", "", "设置角色")
	// 修改账号角色参数
	adminSubCommands[adminUpdateRole].Flags().StringP("role", "r", "", "设置角色")
	// 修改密码参数
	adminSubCommands[adminUpdatePass].Flags().StringP("pass", "p", "", "登录密码")
	// 角色路由权限参数
	adminSubCommands[showRolePerms].Flags().StringP("role", "r", "", "角色名称")
	admin.AddCommand(adminSubCommands.all()...)
	// 管理用户相关参数
	user.PersistentFlags().String("id", "", "用户ID")
	userSubCommands[setUserSession].Flags().IntP("time", "t", 0, "会话过期时间（时间戳）")
	user.AddCommand(userSubCommands.all()...)
	// 关联到主命令
	root.AddCommand(version, runAPIServer, admin, user)
}

// Exec 运行command程序
func Exec() {
	if err := root.Execute(); err != nil {
		log.Errorln(err)
	}
}
