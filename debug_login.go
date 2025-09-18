package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("🔍 登录调试测试")
	fmt.Println("================")

	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("❌ JSON编码错误: %v\n", err)
		return
	}

	fmt.Printf("📤 请求数据: %s\n", string(jsonData))

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/public/login", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ 创建请求错误: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求发送错误: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应错误: %v\n", err)
		return
	}

	fmt.Printf("📊 状态码: %d\n", resp.StatusCode)
	fmt.Printf("📋 响应头: %v\n", resp.Header)

	// 格式化JSON响应
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
		fmt.Printf("📄 响应体:\n%s\n", prettyJSON.String())
	} else {
		fmt.Printf("📄 响应体: %s\n", string(body))
	}

	// 解析关键信息
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if code, ok := result["code"].(float64); ok {
			fmt.Printf("🔑 响应代码: %.0f\n", code)
		}
		if msg, ok := result["msg"].(string); ok {
			fmt.Printf("💬 响应消息: %s\n", msg)
		}
		if data, ok := result["data"].(map[string]interface{}); ok {
			if token, ok := data["token"].(string); ok {
				fmt.Printf("🔐 Token: %s...\n", token[:20])
			}
			if isAdmin, ok := data["is_admin"].(bool); ok {
				fmt.Printf("👑 管理员: %t\n", isAdmin)
			}
		}
	}
}
