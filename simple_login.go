package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("🧪 登录测试开始...")

	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	jsonBytes, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("❌ JSON编码失败: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/api/public/login", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Printf("❌ 创建请求失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ 读取响应失败: %v\n", err)
		return
	}

	if resp.StatusCode == 200 && strings.Contains(string(body), "token") {
		fmt.Println("✅ 登录成功!")
		fmt.Printf("   状态码: %d\n", resp.StatusCode)

		// 尝试解析token（简化显示）
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			if data, ok := result["data"].(map[string]interface{}); ok {
				if token, ok := data["token"].(string); ok {
					fmt.Printf("   Token: %s...\n", token[:20])
				}
			}
		}
	} else {
		fmt.Printf("❌ 登录失败! 状态码: %d\n", resp.StatusCode)
		fmt.Printf("   响应: %s\n", string(body))
	}
}
