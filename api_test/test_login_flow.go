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

func main() {
	// 1. 获取验证码
	captchaResp, err := getCaptcha()
	if err != nil {
		fmt.Printf("获取验证码失败: %v\n", err)
		return
	}
	fmt.Printf("获取验证码成功: CaptchaID=%s, Answer=%s\n", captchaResp.Data.CaptchaID, captchaResp.Data.Answer)

	// 2. 使用验证码进行登录
	loginResp, err := login(LoginRequest{
		Username:    "admin",
		Password:    "admin123",
		CaptchaID:   captchaResp.Data.CaptchaID,
		CaptchaCode: captchaResp.Data.Answer, // 使用返回的answer作为验证码
	})
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}

	if loginResp.Code == 200 {
		fmt.Printf("登录成功! Token=%s\n", loginResp.Data.Token)
	} else {
		fmt.Printf("登录失败: %s\n", loginResp.Msg)
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