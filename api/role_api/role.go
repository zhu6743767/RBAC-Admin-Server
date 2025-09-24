package role_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/gin-gonic/gin"
)

// GetRoleList 获取角色列表
// @Summary 获取角色列表接口
// @Description 查询系统中的角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.Role}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/list [get]
func (r *RoleApi) GetRoleList(c *gin.Context) {
	var roles []models.Role
	if err := global.DB.Find(&roles).Error; err != nil {
		global.Logger.Error("获取角色列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取角色列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": roles,
	})
}

// CreateRole 创建角色
// @Summary 创建角色接口
// @Description 管理员创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body models.Role true "角色信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/create [post]
func (r *RoleApi) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		global.Logger.Error("创建角色参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查角色名是否已存在
	var count int64
	global.DB.Model(&models.Role{}).Where("name = ?", role.Name).Count(&count)
	if count > 0 {
		global.Logger.Error("角色名已存在: " + role.Name)
		c.JSON(400, gin.H{"code": 400, "msg": "角色名已存在"})
		return
	}

	if err := global.DB.Create(&role).Error; err != nil {
		global.Logger.Error("创建角色失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	global.Logger.Infof("管理员创建角色成功: %s", role.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateRole 更新角色
// @Summary 更新角色接口
// @Description 管理员更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body models.Role true "更新的角色信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/update [put]
func (r *RoleApi) UpdateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		global.Logger.Error("更新角色参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if role.ID == 0 {
		global.Logger.Error("更新角色参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Save(&role).Error; err != nil {
		global.Logger.Error("更新角色失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	global.Logger.Infof("管理员更新角色成功: %s", role.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteRole 删除角色
// @Summary 删除角色接口
// @Description 管理员删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id query int true "角色ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/delete [delete]
func (r *RoleApi) DeleteRole(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除角色参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查角色是否有用户关联
	var userRoles []models.UserRole
	global.DB.Where("role_id = ?", id).Find(&userRoles)
	if len(userRoles) > 0 {
		global.Logger.Error("删除角色失败: 角色有用户关联")
		c.JSON(400, gin.H{"code": 400, "msg": "角色有用户关联，无法删除"})
		return
	}

	// 检查角色是否有权限关联
	var rolePermissions []models.RolePermission
	global.DB.Where("role_id = ?", id).Find(&rolePermissions)
	if len(rolePermissions) > 0 {
		// 先删除权限关联
		global.DB.Delete(&models.RolePermission{}, "role_id = ?", id)
	}

	if err := global.DB.Delete(&models.Role{}, id).Error; err != nil {
		global.Logger.Error("删除角色失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("管理员删除角色成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色权限接口
// @Description 查询指定角色的权限列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role_id query int true "角色ID"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.Permission}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/permissions [get]
func (r *RoleApi) GetRolePermissions(c *gin.Context) {
	roleID := c.Query("role_id")
	if roleID == "" {
		global.Logger.Error("获取角色权限参数错误: 角色ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var permissions []models.Permission
	if err := global.DB.Table("permissions").
		Select("permissions.*").
		Joins("join role_permissions on permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		global.Logger.Error("获取角色权限失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// SetRolePermissions 设置角色权限
// @Summary 设置角色权限接口
// @Description 为角色分配权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body struct{RoleID int, PermissionIDs []int} true "角色和权限信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/set-permissions [post]
func (r *RoleApi) SetRolePermissions(c *gin.Context) {
	var req struct {
		RoleID        int   `json:"role_id" binding:"required"`
		PermissionIDs []int `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("设置角色权限参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 事务处理
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除原有权限关联
	if err := tx.Delete(&models.RolePermission{}, "role_id = ?", req.RoleID).Error; err != nil {
		global.Logger.Error("删除原有角色权限关联失败: " + err.Error())
		tx.Rollback()
		c.JSON(500, gin.H{"code": 500, "msg": "设置失败"})
		return
	}

	// 添加新的权限关联
	for _, permissionID := range req.PermissionIDs {
		rolePermission := models.RolePermission{
			RoleID:       req.RoleID,
			PermissionID: permissionID,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			global.Logger.Error("添加角色权限关联失败: " + err.Error())
			tx.Rollback()
			c.JSON(500, gin.H{"code": 500, "msg": "设置失败"})
			return
		}
	}

	// 提交事务
	tx.Commit()

	global.Logger.Infof("管理员设置角色权限成功: 角色ID=%d", req.RoleID)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "设置成功",
	})
}

// GetRoleUsers 获取角色用户
// @Summary 获取角色用户接口
// @Description 查询指定角色下的用户列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role_id query int true "角色ID"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.User}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/role/users [get]
func (r *RoleApi) GetRoleUsers(c *gin.Context) {
	roleID := c.Query("role_id")
	if roleID == "" {
		global.Logger.Error("获取角色用户参数错误: 角色ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var users []models.User
	if err := global.DB.Table("users").
		Select("users.*").
		Joins("join user_roles on users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error; err != nil {
		global.Logger.Error("获取角色用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": users,
	})
}