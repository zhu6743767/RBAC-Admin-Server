package permission_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/gin-gonic/gin"
)

// GetPermissionList 获取权限列表
// @Summary 获取权限列表接口
// @Description 查询系统中的权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.Permission}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/list [get]
func (p *PermissionApi) GetPermissionList(c *gin.Context) {
	var permissions []models.Permission
	if err := global.DB.Find(&permissions).Error; err != nil {
		global.Logger.Error("获取权限列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取权限列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissions,
	})
}

// CreatePermission 创建权限
// @Summary 创建权限接口
// @Description 管理员创建新权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param permission body models.Permission true "权限信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/create [post]
func (p *PermissionApi) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		global.Logger.Error("创建权限参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查权限名是否已存在
	var count int64
	global.DB.Model(&models.Permission{}).Where("name = ?", permission.Name).Count(&count)
	if count > 0 {
		global.Logger.Error("权限名已存在: " + permission.Name)
		c.JSON(400, gin.H{"code": 400, "msg": "权限名已存在"})
		return
	}

	if err := global.DB.Create(&permission).Error; err != nil {
		global.Logger.Error("创建权限失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	global.Logger.Infof("管理员创建权限成功: %s", permission.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdatePermission 更新权限
// @Summary 更新权限接口
// @Description 管理员更新权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param permission body models.Permission true "更新的权限信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/update [put]
func (p *PermissionApi) UpdatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		global.Logger.Error("更新权限参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if permission.ID == 0 {
		global.Logger.Error("更新权限参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Save(&permission).Error; err != nil {
		global.Logger.Error("更新权限失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	global.Logger.Infof("管理员更新权限成功: %s", permission.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeletePermission 删除权限
// @Summary 删除权限接口
// @Description 管理员删除权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id query int true "权限ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/delete [delete]
func (p *PermissionApi) DeletePermission(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除权限参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查权限是否有子权限
	var childPermissions []models.Permission
	global.DB.Where("parent_id = ?", id).Find(&childPermissions)
	if len(childPermissions) > 0 {
		global.Logger.Error("删除权限失败: 权限有子权限")
		c.JSON(400, gin.H{"code": 400, "msg": "权限有子权限，无法删除"})
		return
	}

	// 检查权限是否有角色关联
	var rolePermissions []models.RolePermission
	global.DB.Where("permission_id = ?", id).Find(&rolePermissions)
	if len(rolePermissions) > 0 {
		// 先删除角色关联
		global.DB.Delete(&models.RolePermission{}, "permission_id = ?", id)
	}

	if err := global.DB.Delete(&models.Permission{}, id).Error; err != nil {
		global.Logger.Error("删除权限失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("管理员删除权限成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetPermissionTree 获取权限树结构
// @Summary 获取权限树结构接口
// @Description 查询系统中的权限树结构
// @Tags 权限管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]gin.H{"id":int, "name":string, "children":[]gin.H}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/tree [get]
func (p *PermissionApi) GetPermissionTree(c *gin.Context) {
	var permissions []models.Permission
	if err := global.DB.Order("sort").Find(&permissions).Error; err != nil {
		global.Logger.Error("获取权限树失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	// 构建权限树
	tree := buildPermissionTree(permissions, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}

// buildPermissionTree 构建权限树结构
func buildPermissionTree(permissions []models.Permission, parentID uint) []gin.H {
	tree := make([]gin.H, 0)

	for _, permission := range permissions {
		if permission.ParentID == parentID {
			children := buildPermissionTree(permissions, permission.ID)
			node := gin.H{
				"id":       permission.ID,
				"name":     permission.Name,
				"path":     permission.Path,
				"component": permission.Component,
				"icon":     permission.Icon,
				"type":     permission.Type,
				"sort":     permission.Sort,
				"children": children,
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// GetRolePermissions 获取角色权限
// @Summary 获取角色权限接口
// @Description 查询指定角色的权限ID列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role_id query int true "角色ID"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]int}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/permission/role-permissions [get]
func (p *PermissionApi) GetRolePermissions(c *gin.Context) {
	roleID := c.Query("role_id")
	if roleID == "" {
		global.Logger.Error("获取角色权限参数错误: 角色ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var rolePermissions []models.RolePermission
	if err := global.DB.Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		global.Logger.Error("获取角色权限失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	permissionIDs := make([]int, 0)
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, int(rp.PermissionID))
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": permissionIDs,
	})
}