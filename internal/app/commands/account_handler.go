package commands

import (
	"context"
	"example/internal/app/handlers/assembler"
	"example/internal/app/handlers/dto"
	"example/internal/domain"
	"example/internal/domain/types"
	"example/internal/pkg/logger"
	"example/internal/pkg/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/modood/table"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// AccountHandler 账户控制器
type AccountHandler struct {
	srv domain.AccountService
	lg  *logger.Logger
}

// NewAccountHandler 实例化
func NewAccountHandler(srv domain.AccountService, lg *logger.Logger) *AccountHandler {
	return &AccountHandler{srv: srv, lg: lg}
}

// Create 创建账户
func (c *AccountHandler) Create(email, pass, role string) error {
	data := &dto.AccountCreateDTO{
		Email:    email,
		Password: pass,
		Role:     role,
	}
	if err := utils.Validator(data); err != nil {
		return err
	}
	account, err := new(assembler.Account).ToEntityFromCreateDTO(data)
	if err != nil {
		return err
	}
	if err = c.srv.Create(context.Background(), account); err != nil {
		return err
	}
	log.Infof("账号创建成功，角色为：%s", account.Roles)
	return nil
}

// ShowAccounts 显示账户列表
func (c *AccountHandler) ShowAccounts() error {
	list, err := c.srv.GetAccounts(context.Background(), nil)
	if err != nil {
		return err
	}
	var result []dto.ResAccountInfo
	for _, account := range list {
		result = append(result, new(assembler.Account).ToAccountInfo(account))
	}
	fmt.Println(table.Table(result))
	return nil
}

// UpdateRole 修改管理员角色
func (c *AccountHandler) UpdateRole(email, role string) error {
	if !govalidator.IsEmail(email) {
		return errors.New("请输入正确的邮箱账号")
	}
	roles, err := types.ParseRoles(role, false)
	if err != nil {
		return err
	}
	if err = c.srv.UpdateRoleByEmail(context.Background(), email, roles); err != nil {
		return err
	}
	log.Infof("角色已修改为：%s", roles)
	return nil
}

// UpdatePass 修改管理员密码
func (c *AccountHandler) UpdatePass(email, pass string) error {
	var data = &struct {
		Email string `valid:"valid_email~请输入正确的邮箱账号"`
		Pass  string `valid:"valid_pass~密码不符合要求"`
	}{
		Email: email,
		Pass:  pass,
	}
	if err := utils.Validator(data); err != nil {
		return err
	}
	if err := c.srv.UpdatePassByEmail(context.Background(), email, pass); err != nil {
		return err
	}
	log.Info("密码修改成功")
	return nil
}

// ShowAccountRole 显示账户角色
func (c *AccountHandler) ShowAccountRole(email string) error {
	if !govalidator.IsEmail(email) {
		return errors.New("请输入正确的邮箱账号")
	}

	account, err := c.srv.GetAccountByEmail(context.Background(), email)
	if err != nil {
		return err
	}
	c.lg.Infof("管理员：%s, 角色：%s", email, account.Roles.String())
	return nil
}

// ShowAccountPerms 显示账户拥有路由权限
func (c *AccountHandler) ShowAccountPerms(email string) error {
	if !govalidator.IsEmail(email) {
		return errors.New("请输入正确的邮箱账号")
	}

	account, err := c.srv.GetAccountByEmail(context.Background(), email)
	if err != nil {
		return err
	}
	c.lg.Infof("管理员：%s, 角色：%s, 拥有路由权限：", email, account.Roles.String())
	account.Roles.ShowPerms()
	return nil
}
