package email_api

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"rbac_admin_server/global"
	"rbac_admin_server/utils"
	"rbac_admin_server/utils/captcha"
	"rbac_admin_server/utils/email"
)

// EmailApi 邮件API控制器
// 提供邮件相关的API接口
type EmailApi struct {}

// SendEmailRequest 发送邮件验证码请求结构
// 定义客户端发送邮件验证码时需要提供的参数
type SendEmailRequest struct {
	Email       string `json:"email" binding:"required,email"` // 收件人邮箱（必填，需验证邮箱格式）
	CaptchaID   string `json:"captchaID"`                      // 图片验证码ID
	CaptchaCode string `json:"captchaCode"`                    // 图片验证码内容
}

// SendEmailResponse 发送邮件验证码响应结构
// 定义服务器返回的邮件发送结果
type SendEmailResponse struct {
	EmailID string `json:"emailID"` // 邮件唯一标识，用于后续验证
}

// SendEmailView 发送邮件验证码接口
// 处理客户端发送邮件验证码的请求
func (e *EmailApi) SendEmailView(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": utils.ERROR_INVALID_PARAM,
			"msg":  utils.GetErrMsg(utils.ERROR_INVALID_PARAM) + ": " + err.Error(),
		})
		return
	}

	// 检查邮箱配置是否完整
	if !global.Config.Email.Verify() {
		c.JSON(400, gin.H{
			"code": utils.ERROR_EMAIL_CONFIG,
			"msg":  utils.GetErrMsg(utils.ERROR_EMAIL_CONFIG),
		})
		return
	}

	// 如果启用了验证码，验证图片验证码
	if global.Config.Captcha.Enable {
		if req.CaptchaID == "" || req.CaptchaCode == "" {
			c.JSON(400, gin.H{
				"code": utils.ERROR_INVALID_PARAM,
				"msg":  "请输入验证码",
			})
			return
		}

		// 验证图片验证码
		if !captcha.CaptchaStore.Verify(req.CaptchaID, req.CaptchaCode, true) {
			c.JSON(400, gin.H{
				"code": utils.ERROR_CAPTCHA_WRONG,
				"msg":  utils.GetErrMsg(utils.ERROR_CAPTCHA_WRONG),
			})
			return
		}
	}

	// 生成唯一的邮箱验证ID
	emailID := uuid.New().String()

	// 生成6位数字验证码
	code := captcha.EmailStore.GenerateAndStoreEmailCode(req.Email, 5*time.Minute)

	// 存储邮箱验证码信息
	email.Set(emailID, req.Email, code)

	// 构建邮件内容
	content := fmt.Sprintf("您正在完成用户注册，您的验证码为 %s ，请在5分钟内使用，过时无效！", code)

	// 发送邮件
	err := email.SendEmail(req.Email, "用户注册", content)
	if err != nil {
		global.Logger.Error("发送验证码邮件失败: " + err.Error())
		c.JSON(500, gin.H{
			"code": utils.ERROR_EMAIL_SEND,
			"msg":  utils.GetErrMsg(utils.ERROR_EMAIL_SEND),
		})
		return
	}

	// 记录日志
	global.Logger.Info("验证码邮件发送成功: " + req.Email)

	// 返回成功响应
	c.JSON(200, gin.H{
		"code": utils.SUCCESS,
		"msg":  "验证码邮件发送成功",
		"data": SendEmailResponse{
			EmailID: emailID,
		},
	})
}

// RegisterRoutes 注册邮件相关路由
// 将邮件API接口注册到路由系统中
func (e *EmailApi) RegisterRoutes(r *gin.RouterGroup) {
	emailGroup := r.Group("email")
	{
		emailGroup.POST("/send", e.SendEmailView) // 发送邮件验证码接口
	}
}