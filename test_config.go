package main

import (
	"fmt"
	"log"
	"rbac.admin/config"
)

func main() {
	// 测试开发配置
	fmt.Println("🧪 测试配置加载...")
	
	// 测试开发配置
	fmt.Println("\n📋 测试开发配置:")
	cfgDev, err := config.Load("settings_dev.yaml")
	if err != nil {
		log.Fatalf("开发配置加载失败: %v", err)
	}
	fmt.Printf("✅ 开发配置加载成功: %+v\n", cfgDev.App)
	
	// 测试测试配置
	fmt.Println("\n📋 测试测试配置:")
	cfgTest, err := config.Load("settings_test.yaml")
	if err != nil {
		log.Fatalf("测试配置加载失败: %v", err)
	}
	fmt.Printf("✅ 测试配置加载成功: %+v\n", cfgTest.App)
	
	// 测试生产配置
	fmt.Println("\n📋 测试生产配置:")
	cfgProd, err := config.Load("settings_prod.yaml")
	if err != nil {
		log.Fatalf("生产配置加载失败: %v", err)
	}
	fmt.Printf("✅ 生产配置加载成功: %+v\n", cfgProd.App)
	
	fmt.Println("\n🎉 所有配置测试通过!")
}