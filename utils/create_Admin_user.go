package main

import (
	"fmt"
	"os"
	"rbac.admin/config"
	"rbac.admin/core"
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/pwd"
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
	
	// 创建管理员用户
	CreateAdminUser()
}

func CreateAdminUser() {
	// 创建用户
	fmt.Println("请输入用户名:")
	var username string
	fmt.Scanln(&username)
	
	if username == "" {
		fmt.Println("用户名不能为空")
		return
	}
	
	// 检查用户是否已存在
	var user models.User
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err == nil {
		fmt.Println("用户已存在")
		global.Logger.Error("用户已存在")
		return
	}
	
	fmt.Println("请输入密码:")
	var password string
	fmt.Scanln(&password)
	
	fmt.Println("请再次输入密码:")
	var rePassword string
	fmt.Scanln(&rePassword)
	
	if password != rePassword {
		fmt.Println("两次输入密码不一致")
		global.Logger.Error("两次输入密码不一致")
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
		global.Logger.Errorf("创建用户时出错: %v", err)
		return
	}
	
	fmt.Println("创建用户成功")
	global.Logger.Info("创建用户成功")
}
