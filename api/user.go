package api

import (
	"rbac.admin/global"
	"rbac.admin/models"
	"rbac.admin/pwd"
	"rbac.admin/utils"

	"github.com/gin-gonic/gin"
)

// Login 用户登录
func Login(c *gin.Context) {
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
	if !pwd.ComparePassword(user.Password, req.Password) {
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

// Register 用户注册
func Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	// 密码加密
	hashPwd := pwd.HashedPassword(req.Password)

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashPwd,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
		IsAdmin:  false,
	}

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("创建用户失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "注册成功"})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	var req struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=10"`
		Keyword  string `form:"keyword"`
		Status   int    `form:"status"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var users []models.User
	query := global.DB.Model(&models.User{})

	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", 
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		global.Logger.Error("获取用户列表失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":  users,
			"total": total,
			"page":  req.Page,
			"page_size": req.PageSize,
		},
	})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone"`
		Status   int    `json:"status,default=1"`
		IsAdmin  bool   `json:"is_admin,default=false"`
		DepartmentID uint   `json:"department_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	// 密码加密
	hashPwd := pwd.HashedPassword(req.Password)

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashPwd,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
		IsAdmin:  req.IsAdmin,
		DepartmentID: req.DepartmentID,
	}

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Error("创建用户失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "创建成功"})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	var req struct {
		ID       uint   `json:"id" binding:"required"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Status   int    `json:"status"`
		IsAdmin  bool   `json:"is_admin"`
		DeptID   uint   `json:"dept_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 更新用户
	updates := map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
		"is_admin": req.IsAdmin,
		"dept_id":  req.DeptID,
	}

	if err := global.DB.Model(&models.User{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		global.Logger.Error("更新用户失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "更新成功"})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.User{}, req.ID).Error; err != nil {
		global.Logger.Error("删除用户失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}

// RefreshToken 刷新令牌
func RefreshToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{
			"code": 401,
			"msg": "请提供token",
		})
		return
	}

	// 移除可能的Bearer前缀
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	newToken, err := global.RefreshToken(token)
	if err != nil {
		global.Logger.Warnf("刷新token失败: %s", err.Error())
		c.JSON(401, gin.H{
			"code": 401,
			"msg": "刷新失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg": "刷新成功",
		"data": gin.H{
			"token": newToken,
		},
	})