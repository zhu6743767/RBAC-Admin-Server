package user_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"
	"rbac_admin_server/utils"
	"rbac_admin_server/utils/captcha"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录接口
// @Description 用户登录系统获取访问令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body struct{Username string, Password string, Captcha string, CaptchaId string} true "登录信息"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"token":string, "refresh_token":string, "user":models.User, "is_admin":bool}}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 401 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /public/login [post]
func (u *UserApi) Login(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		CaptchaID   string `json:"captchaID"`
		CaptchaCode string `json:"captchaCode"`
	}

	// 如果启用了验证码，需要验证
	if global.Config.Captcha.Enable {
		// 调整绑定规则，增加验证码字段的必填要求
		type LoginWithCaptchaReq struct {
			Username    string `json:"username" binding:"required"`
			Password    string `json:"password" binding:"required"`
			CaptchaID   string `json:"captchaID" binding:"required"`
			CaptchaCode string `json:"captchaCode" binding:"required"`
		}
		var captchaReq LoginWithCaptchaReq
		if err := c.ShouldBindJSON(&captchaReq); err != nil {
			global.Logger.Error("登录参数错误: " + err.Error())
			c.JSON(400, gin.H{"code": utils.ERROR_INVALID_PARAM, "msg": utils.GetErrMsg(utils.ERROR_INVALID_PARAM)})
			return
		}
		// 验证验证码
		if !captcha.CaptchaStore.Verify(captchaReq.CaptchaID, captchaReq.CaptchaCode, true) {
			global.Logger.Error("验证码错误: " + captchaReq.Username)
			c.JSON(400, gin.H{"code": utils.ERROR_CAPTCHA_WRONG, "msg": utils.GetErrMsg(utils.ERROR_CAPTCHA_WRONG)})
			return
		}
		// 复制数据到原始请求结构
		req.Username = captchaReq.Username
		req.Password = captchaReq.Password
	} else {
		// 不启用验证码时，只绑定基础字段
		if err := c.ShouldBindJSON(&req); err != nil {
			global.Logger.Error("登录参数错误: " + err.Error())
			c.JSON(400, gin.H{"code": utils.ERROR_INVALID_PARAM, "msg": utils.GetErrMsg(utils.ERROR_INVALID_PARAM)})
			return
		}
	}

	// 查询用户
	var user models.User
	err := global.DB.Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		global.Logger.Error("用户不存在: " + req.Username)
		c.JSON(401, gin.H{"code": utils.ERROR_USER_NOT_EXIST, "msg": utils.GetErrMsg(utils.ERROR_USER_NOT_EXIST)})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		global.Logger.Error("用户已被禁用: " + req.Username)
		c.JSON(401, gin.H{"code": utils.ERROR, "msg": "用户已被禁用"})
		return
	}

	// 验证密码
	if !utils.ComparePassword(user.Password, req.Password) {
		global.Logger.Error("密码验证失败: " + req.Username)
		c.JSON(401, gin.H{"code": utils.ERROR_PASSWORD_WRONG, "msg": utils.GetErrMsg(utils.ERROR_PASSWORD_WRONG)})
		return
	}

	// 获取用户角色列表
	roleList, err := global.GetUserRoles(user.ID)
	if err != nil {
		global.Logger.Error("获取用户角色失败: " + err.Error())
		c.JSON(500, gin.H{"code": utils.ERROR, "msg": utils.GetErrMsg(utils.ERROR)})
		return
	}

	// 生成访问令牌
	accessToken, err := global.GenerateToken(global.ClaimsUserInfo{
		UserID:   user.ID,
		Username: user.Username,
		RoleList: roleList,
	})
	if err != nil {
		global.Logger.Error("生成访问令牌失败: ", err)
		c.JSON(500, gin.H{"code": utils.ERROR, "msg": utils.GetErrMsg(utils.ERROR)})
		return
	}

	// 生成刷新令牌
	refreshToken, err := global.GenerateRefreshToken(user.ID)
	if err != nil {
		global.Logger.Error("生成刷新令牌失败: ", err)
		c.JSON(500, gin.H{"code": utils.ERROR, "msg": utils.GetErrMsg(utils.ERROR)})
		return
	}

	global.Logger.Infof("用户登录成功: %s", user.Username)
	c.JSON(200, gin.H{
		"code": utils.SUCCESS,
		"msg":  utils.GetErrMsg(utils.SUCCESS),
		"data": gin.H{
			"token":         accessToken,
			"refresh_token": refreshToken,
			"user":          user,
			"is_admin":      user.IsAdmin,
		},
	})
}

// RefreshToken 刷新JWT令牌
// @Summary 刷新JWT令牌接口
// @Description 基于有效刷新令牌生成新的访问令牌和刷新令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param refresh body struct{RefreshToken string} true "刷新令牌请求"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"token":string, "refresh_token":string}}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 401 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /public/refresh-token [post]
func (u *UserApi) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("刷新令牌参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": utils.ERROR, "msg": "参数错误"})
		return
	}

	// 刷新令牌
	newAccessToken, newRefreshToken, err := global.RefreshToken(req.RefreshToken)
	if err != nil {
		global.Logger.Error("刷新令牌失败: " + err.Error())
		c.JSON(401, gin.H{"code": utils.ERROR_TOKEN_INVALID, "msg": err.Error()})
		return
	}

	global.Logger.Info("令牌刷新成功")
	c.JSON(200, gin.H{
		"code": utils.SUCCESS,
		"msg":  "令牌刷新成功",
		"data": gin.H{
			"token":         newAccessToken,
			"refresh_token": newRefreshToken,
		},
	})
}