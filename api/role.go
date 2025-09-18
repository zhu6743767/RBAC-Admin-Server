package api

import (
	"rbac.admin/global"
	"rbac.admin/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRoleList 获取角色列表
func GetRoleList(c *gin.Context) {
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

	var roles []models.Role
	query := global.DB.Model(&models.Role{})

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&roles).Error; err != nil {
		logrus.Error("获取角色列表失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":  roles,
			"total": total,
			"page":  req.Page,
			"page_size": req.PageSize,
		},
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status,default=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查角色名是否已存在
	var count int64
	global.DB.Model(&models.Role{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{"code": 400, "msg": "角色名已存在"})
		return
	}

	// 创建角色
	role := models.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := global.DB.Create(&role).Error; err != nil {
		logrus.Error("创建角色失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "创建成功"})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	var req struct {
		ID          uint   `json:"id" binding:"required"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 更新角色
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"status":      req.Status,
	}

	if err := global.DB.Model(&models.Role{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		logrus.Error("更新角色失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "更新成功"})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.Role{}, req.ID).Error; err != nil {
		logrus.Error("删除角色失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}