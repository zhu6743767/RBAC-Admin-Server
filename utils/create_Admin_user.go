package main

import (
	"fmt"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/pwd"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// 初始化系统
	config := core.InitSystem()
	if config == nil {
		fmt.Println("系统初始化失败")
		return
	}
	defer core.CleanupSystem()
	
	// 创建管理员用户
	CreateAdminUser()
}

func CreateAdminUser() {
	// 创建用户
	fmt.Println("请输入用户名")
	var username string
	fmt.Scanln(&username)
	
	// 检查用户是否已存在
	var user models.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err == nil {
		fmt.Println("用户已存在")
		logrus.Error("用户已存在")
		return
	}
	
	fmt.Println("请输入密码")
	password, err := terminal.ReadPassword(0) // 从标准输入读取密码
	if err != nil {
		fmt.Println("输入密码时出错")
		logrus.Errorf("输入密码时出错: %v", err)
		return
	}
	
	fmt.Println("请再次输入密码")
	rePassword, err := terminal.ReadPassword(0) // 从标准输入再次读取密码
	if err != nil {
		fmt.Println("输入密码时出错")
		logrus.Errorf("输入密码时出错: %v", err)
		return
	}
	
	if string(password) != string(rePassword) {
		fmt.Println("两次输入密码不一致")
		logrus.Error("两次输入密码不一致")
		return
	}

	// 密码加密
	hashPwd := pwd.HashedPassword(string(password))
	
	// 创建管理员用户
	newUser := models.User{
		Username: username,
		Password: hashPwd,
		Status:   1, // 正常状态
	}
	
	err = global.DB.Create(&newUser).Error
	if err != nil {
		fmt.Println("创建用户时出错")
		logrus.Errorf("创建用户时出错: %v", err)
		return
	}
	
	fmt.Println("创建用户成功")
	logrus.Info("创建用户成功")
}
