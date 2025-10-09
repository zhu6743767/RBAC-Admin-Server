package captcha_api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"rbac_admin_server/global"
	"rbac_admin_server/utils/captcha"
)

// CaptchaApi 验证码API控制器
// 提供验证码图片生成和验证功能
type CaptchaApi struct {}

// CaptchaResponse 验证码响应结构
// 包含验证码ID、图片Base64编码和答案
type CaptchaResponse struct {
	CaptchaID string `json:"captchaId"` // 验证码唯一标识
	Image     string `json:"image"`     // 验证码图片Base64编码
	Answer    string `json:"answer"`    // 验证码答案（仅用于调试）
}

// GetCaptcha 生成验证码图片
// 处理客户端获取验证码图片的请求
func (c *CaptchaApi) GetCaptcha(ctx *gin.Context) {
	// 配置验证码驱动参数
	driver := base64Captcha.DriverString{
		Height:          80,          // 验证码图片高度
		Width:           240,         // 验证码图片宽度
		NoiseCount:      2,           // 干扰线数量
		ShowLineOptions: 4,           // 显示线条选项
		Length:          4,           // 验证码长度
		Source:          "1234567890", // 验证码字符集
	}

	// 创建验证码实例
	cp := base64Captcha.NewCaptcha(&driver, captcha.CaptchaStore)

	// 生成验证码
	id, b64s, answer, err := cp.Generate()
	if err != nil {
		global.Logger.Error("生成验证码失败: " + err.Error())
		ctx.JSON(500, gin.H{
			"code": 500,
			"msg":  "生成验证码失败",
		})
		return
	}

	// 返回验证码信息
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "生成验证码成功",
		"data": CaptchaResponse{
			CaptchaID: id,
			Image:     b64s,
			Answer:    answer, // 实际生产环境中应移除此字段
		},
	})
}

// RegisterRoutes 注册验证码相关路由
// 将验证码API接口注册到路由系统中
func (c *CaptchaApi) RegisterRoutes(r *gin.RouterGroup) {
	captchaGroup := r.Group("captcha")
	{
		captchaGroup.GET("/get", c.GetCaptcha) // 获取验证码图片接口
	}
}