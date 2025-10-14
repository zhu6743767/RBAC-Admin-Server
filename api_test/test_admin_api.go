package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://127.0.0.1:8080"

// CaptchaResponse 验证码响应结构
type CaptchaResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CaptchaID string `json:"captchaId"`
		Image     string `json:"image"`
		Answer    string `json:"answer"`
	} `json:"data"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaID   string `json:"captchaID"`
	CaptchaCode string `json:"captchaCode"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		User         struct {
			ID uint `json:"id"`
		} `json:"user"`
		IsAdmin bool `json:"is_admin"`
	} `json:"data"`
}

// UserListResponse 用户列表响应结构
type UserListResponse struct {
	Code int                       `json:"code"`
	Msg  string                    `json:"msg"`
	Data []map[string]interface{}  `json:"data"`
}

// RoleListResponse 角色列表响应结构
type RoleListResponse struct {
	Code int                       `json:"code"`
	Msg  string                    `json:"msg"`
	Data []map[string]interface{}  `json:"data"`
}

func main() {
	// 1. 获取验证码
	captchaResp, err := getCaptcha()
	if err != nil {
		fmt.Printf("获取验证码失败: %v\n", err)
		return
	}
	fmt.Printf("获取验证码成功: CaptchaID=%s\n", captchaResp.Data.CaptchaID)

	// 2. 使用验证码进行登录
	loginResp, err := login(LoginRequest{
		Username:    "admin",
		Password:    "admin123",
		CaptchaID:   captchaResp.Data.CaptchaID,
		CaptchaCode: captchaResp.Data.Answer, // 使用返回的answer作为验证码
	})
	if err != nil || loginResp.Code != 200 {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	fmt.Printf("登录成功!\n")
	fmt.Printf("Token: %s\n", loginResp.Data.Token)
	fmt.Printf("是否为管理员: %v\n", loginResp.Data.IsAdmin)

	// 3. 测试管理员接口
	token := loginResp.Data.Token

	// 3.1 测试获取用户列表
	fmt.Println("\n--- 测试获取用户列表 ---")
	userList, err := getUserList(token)
	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
	} else {
		fmt.Printf("用户总数: %d\n", len(userList.Data))
		fmt.Printf("用户列表预览: 用户1: %v\n", userList.Data[0]["username"])
		if len(userList.Data) > 1 {
			fmt.Printf("              用户2: %v\n", userList.Data[1]["username"])
		}
	}

	// 3.2 测试获取角色列表
	fmt.Println("\n--- 测试获取角色列表 ---")
	roleList, err := getRoleList(token)
	if err != nil {
		fmt.Printf("获取角色列表失败: %v\n", err)
	} else {
		fmt.Printf("角色总数: %d\n", len(roleList.Data))
		fmt.Printf("角色列表预览: ")
		for i, role := range roleList.Data {
			fmt.Printf("角色%d: %v ", i+1, role["name"])
		}
		fmt.Println()
	}
}

// getCaptcha 调用获取验证码接口
func getCaptcha() (*CaptchaResponse, error) {
	resp, err := http.Get(baseURL + "/public/captcha/get")
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var captchaResp CaptchaResponse
	if err := json.Unmarshal(body, &captchaResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w, 原始响应: %s", err, string(body))
	}

	if captchaResp.Code != 200 {
		return nil, fmt.Errorf("获取验证码失败: %s", captchaResp.Msg)
	}

	return &captchaResp, nil
}

// login 调用登录接口
func login(req LoginRequest) (*LoginResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/public/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w, 原始响应: %s", err, string(body))
	}

	return &loginResp, nil
}

// getUserList 获取用户列表
func getUserList(token string) (*UserListResponse, error) {
	httpReq, err := http.NewRequest("GET", baseURL+"/admin/user/list", nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var userListResp UserListResponse
	if err := json.Unmarshal(body, &userListResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w, 原始响应: %s", err, string(body))
	}

	if userListResp.Code != 200 {
		return nil, fmt.Errorf("获取用户列表失败: %s", userListResp.Msg)
	}

	return &userListResp, nil
}

// getRoleList 获取角色列表
func getRoleList(token string) (*RoleListResponse, error) {
	httpReq, err := http.NewRequest("GET", baseURL+"/admin/role/list", nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	var roleListResp RoleListResponse
	if err := json.Unmarshal(body, &roleListResp); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %w, 原始响应: %s", err, string(body))
	}

	if roleListResp.Code != 200 {
		return nil, fmt.Errorf("获取角色列表失败: %s", roleListResp.Msg)
	}

	return &roleListResp, nil
}