package main

import (
	"fmt"
	"os"
	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/pwd"

	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	cfg, err := config.Load("../settings_dev.yaml")
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		os.Exit(1)
	}
	
	// 初始化系统
	err = core.InitSystem(cfg)
	if err != nil {
		fmt.Printf("系统初始化失败: %v\n", err)
		os.Exit(1)
	}
	
	// 创建测试管理员用户
	createTestAdmin()
}

func createTestAdmin() {
	username := "admin"
	password := "admin123"
	
	// 检查用户是否已存在
	var user models.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err == nil {
		fmt.Printf("用户 %s 已存在\n", username)
		logrus.Info("用户已存在")
		return
	}
	
	// 密码加密
	hashPwd := pwd.HashedPassword(password)
	
	// 创建管理员用户
	newUser := models.User{
		Username: username,
		Password: hashPwd,
		Nickname: "管理员",
		Status:   1, // 正常状态
		IsAdmin:  true,
	}
	
	err = global.DB.Create(&newUser).Error
	if err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		logrus.Errorf("创建用户失败: %v", err)
		return
	}
	
	fmt.Printf("管理员用户 %s 创建成功！\n", username)
	fmt.Printf("用户名: %s\n", username)
	fmt.Printf("密码: %s\n", password)
	logrus.Info("创建用户成功")
}