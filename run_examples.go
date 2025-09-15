package main

import (
	"fmt"
	"rbac.admin/examples"
)

func main() {
	fmt.Println("=== RBAC Admin Server 示例程序 ===")
	fmt.Println("\n可用的示例:")
	fmt.Println("1. 数据库使用示例")
	fmt.Println("2. 日志工厂使用示例")
	fmt.Println("3. 完整示例")
	fmt.Println("4. 退出")

	var choice int
	fmt.Print("\n请选择要运行的示例 (1-4): ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("\n--- 运行数据库使用示例 ---")
		examples.DatabaseUsageExample()
	case 2:
		fmt.Println("\n--- 运行日志工厂使用示例 ---")
		examples.FactoryUsageExample()
	case 3:
		fmt.Println("\n--- 运行完整示例 ---")
		examples.MainExample()
	case 4:
		fmt.Println("退出程序")
		return
	default:
		fmt.Println("无效的选择")
	}

	fmt.Println("\n示例运行完成！")
}
