package api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/gin-gonic/gin"
)

// GetPermissionList 获取权限列表
func GetPermissionList(c *gin.Context) {
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

	var permissions []models.Permission
	query := global.DB.Model(&models.Permission{})

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&permissions).Error; err != nil {
		global.Logger.Error("获取权限列表失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":  permissions,
			"total": total,
			"page":  req.Page,
			"page_size": req.PageSize,
		},
	})
}

// CreatePermission 创建权限
func CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status,default=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查权限名是否已存在
	var count int64
	global.DB.Model(&models.Permission{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(400, gin.H{"code": 400, "msg": "权限名已存在"})
		return
	}

	// 创建权限
	permission := models.Permission{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := global.DB.Create(&permission).Error; err != nil {
		global.Logger.Error("创建权限失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "创建成功"})
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
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

	// 更新权限
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"status":      req.Status,
	}

	if err := global.DB.Model(&models.Permission{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		global.Logger.Error("更新权限失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "更新成功"})
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.Permission{}, req.ID).Error; err != nil {
		global.Logger.Error("删除权限失败: ", err)
		c.JSON(500, gin.H{"code": 500, "msg": "系统错误"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}