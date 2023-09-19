package commands

import "github.com/spf13/cobra"

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
	return commander.depends().adminCreate(email, pass, role)
}
