package examples

import "fmt"

// RunAllExamples 运行所有示例
func RunAllExamples() {
	fmt.Println("=== 运行所有示例 ===")
	
	// 运行数据库示例
	fmt.Println("\n1. 运行数据库示例...")
	DatabaseUsageExample()
	
	// 运行工厂模式示例
	fmt.Println("\n2. 运行工厂模式示例...")
	FactoryUsageExample()
	
	// 运行主示例
	fmt.Println("\n3. 运行主示例...")
	MainExample()
	
	fmt.Println("\n=== 所有示例运行完成 ===")
}

// 独立运行
// go run examples/run_examples.go
// func main() {
// 	RunAllExamples()
// }