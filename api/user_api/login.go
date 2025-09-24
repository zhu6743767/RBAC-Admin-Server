package user_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"
	"rbac_admin_server/utils"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
// @Summary 用户登录接口
// @Description 验证用户身份并返回JWT令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body struct{Username string, Password string} true "登录信息"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"token":string, "user":models.User, "is_admin":bool}}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 401 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/public/login [post]
func (u *UserApi) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("登录参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": utils.ERROR, "msg": "参数错误"})
		return
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
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	// 生成token
	token, err := global.GenerateToken(global.ClaimsUserInfo{
		UserID:   user.ID,
		Username: user.Username,
		RoleList: roleList,
	})
	if err != nil {
		global.Logger.Error("生成token失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	global.Logger.Infof("用户登录成功: %s", user.Username)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"user":     user,
			"is_admin": user.IsAdmin,
		},
	})
}

// RefreshToken 刷新JWT令牌
// @Summary 刷新JWT令牌接口
// @Description 基于有效令牌生成新令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param refresh body struct{Token string} true "刷新令牌请求"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"token":string}}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 401 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/public/refresh-token [post]
func (u *UserApi) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("刷新令牌参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": utils.ERROR, "msg": "参数错误"})
		return
	}

	// 解析现有令牌
	claims, err := global.ParseToken(req.Token)
	if err != nil {
		global.Logger.Error("令牌解析失败: " + err.Error())
		c.JSON(401, gin.H{"code": utils.ERROR_TOKEN_INVALID, "msg": utils.GetErrMsg(utils.ERROR_TOKEN_INVALID)})
		return
	}

	// 生成新令牌
	newToken, err := global.GenerateToken(claims)
	if err != nil {
		global.Logger.Error("生成新令牌失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	global.Logger.Infof("用户 %s 令牌刷新成功", claims.Username)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "令牌刷新成功",
		"data": gin.H{"token": newToken},
	})
}