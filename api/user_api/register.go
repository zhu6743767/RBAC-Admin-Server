package user_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"
	"rbac_admin_server/utils"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// @Summary 用户注册接口
// @Description 创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param register body struct{Username string, Password string, Nickname string, Email string, Phone string} true "注册信息"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":models.User}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/public/register [post]
func (u *UserApi) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("注册参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		global.Logger.Error("用户名已存在: " + req.Username)
		c.JSON(400, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	global.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		global.Logger.Error("邮箱已存在: " + req.Email)
		c.JSON(400, gin.H{"code": 400, "msg": "邮箱已存在"})
		return
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		global.DB.Model(&models.User{}).Where("phone = ?", req.Phone).Count(&count)
		if count > 0 {
			global.Logger.Error("手机号已存在: " + req.Phone)
			c.JSON(400, gin.H{"code": 400, "msg": "手机号已存在"})
			return
		}
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: utils.MakePassword(req.Password),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1, // 默认为启用状态
	}

	// 保存用户到数据库
	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("创建用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "注册失败"})
		return
	}

	global.Logger.Infof("用户注册成功: %s", req.Username)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": user,
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表接口
// @Description 查询系统中的用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.User}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/user/list [get]
func (u *UserApi) GetUserList(c *gin.Context) {
	var users []models.User
	if err := global.DB.Find(&users).Error; err != nil {
		global.Logger.Error("获取用户列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取用户列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": users,
	})
}

// CreateUser 创建用户
// @Summary 创建用户接口
// @Description 管理员创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/user/create [post]
func (u *UserApi) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		global.Logger.Error("创建用户参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 密码加密
	user.Password = utils.MakePassword(user.Password)

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("创建用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	global.Logger.Infof("管理员创建用户成功: %s", user.Username)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateUser 更新用户
// @Summary 更新用户接口
// @Description 管理员更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.User true "更新的用户信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/user/update [put]
func (u *UserApi) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		global.Logger.Error("更新用户参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 不允许直接更新密码
	if user.Password != "" {
		user.Password = utils.MakePassword(user.Password)
	}

	if err := global.DB.Save(&user).Error; err != nil {
		global.Logger.Error("更新用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	global.Logger.Infof("管理员更新用户成功: %s", user.Username)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteUser 删除用户
// @Summary 删除用户接口
// @Description 管理员删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/user/delete [delete]
func (u *UserApi) DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除用户参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.User{}, id).Error; err != nil {
		global.Logger.Error("删除用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("管理员删除用户成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}