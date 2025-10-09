package menu_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/gin-gonic/gin"
)

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表接口
// @Description 查询系统中的菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.Permission}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/list [get]
func (m *MenuApi) GetMenuList(c *gin.Context) {
	var menus []models.Permission
	if err := global.DB.Where("type in (1, 2)").Order("sort").Find(&menus).Error; err != nil {
		global.Logger.Error("获取菜单列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取菜单列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": menus,
	})
}

// CreateMenu 创建菜单
// @Summary 创建菜单接口
// @Description 管理员创建新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body models.Permission true "菜单信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/create [post]
func (m *MenuApi) CreateMenu(c *gin.Context) {
	var menu models.Permission
	if err := c.ShouldBindJSON(&menu); err != nil {
		global.Logger.Error("创建菜单参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 菜单类型只能是dir(目录)或menu(菜单)
	if menu.Type != "dir" && menu.Type != "menu" {
		global.Logger.Error("创建菜单参数错误: 菜单类型错误")
		c.JSON(400, gin.H{"code": 400, "msg": "菜单类型只能是目录或菜单"})
		return
	}

	// 检查菜单名称是否已存在
	var count int64
	global.DB.Model(&models.Permission{}).Where("name = ?", menu.Name).Count(&count)
	if count > 0 {
		global.Logger.Error("菜单名称已存在: " + menu.Name)
		c.JSON(400, gin.H{"code": 400, "msg": "菜单名称已存在"})
		return
	}

	if err := global.DB.Create(&menu).Error; err != nil {
		global.Logger.Error("创建菜单失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	global.Logger.Infof("管理员创建菜单成功: %s", menu.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateMenu 更新菜单
// @Summary 更新菜单接口
// @Description 管理员更新菜单信息
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param menu body models.Permission true "更新的菜单信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/update [put]
func (m *MenuApi) UpdateMenu(c *gin.Context) {
	var menu models.Permission
	if err := c.ShouldBindJSON(&menu); err != nil {
		global.Logger.Error("更新菜单参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if menu.ID == 0 {
		global.Logger.Error("更新菜单参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 菜单类型只能是dir(目录)或menu(菜单)
	if menu.Type != "dir" && menu.Type != "menu" {
		global.Logger.Error("更新菜单参数错误: 菜单类型错误")
		c.JSON(400, gin.H{"code": 400, "msg": "菜单类型只能是目录或菜单"})
		return
	}

	if err := global.DB.Save(&menu).Error; err != nil {
		global.Logger.Error("更新菜单失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	global.Logger.Infof("管理员更新菜单成功: %s", menu.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteMenu 删除菜单
// @Summary 删除菜单接口
// @Description 管理员删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Param id query int true "菜单ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/delete [delete]
func (m *MenuApi) DeleteMenu(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除菜单参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查菜单是否有子菜单
	var childMenus []models.Permission
	global.DB.Where("parent_id = ?", id).Find(&childMenus)
	if len(childMenus) > 0 {
		global.Logger.Error("删除菜单失败: 菜单有子菜单")
		c.JSON(400, gin.H{"code": 400, "msg": "菜单有子菜单，无法删除"})
		return
	}

	// 检查菜单是否有权限关联
	var rolePermissions []models.RolePermission
	global.DB.Where("permission_id = ?", id).Find(&rolePermissions)
	if len(rolePermissions) > 0 {
		// 先删除角色关联
		global.DB.Delete(&models.RolePermission{}, "permission_id = ?", id)
	}

	if err := global.DB.Delete(&models.Permission{}, id).Error; err != nil {
		global.Logger.Error("删除菜单失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("管理员删除菜单成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetMenuTree 获取菜单树结构
// @Summary 获取菜单树结构接口
// @Description 查询系统中的菜单树结构
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]gin.H{"id":int, "name":string, "children":[]gin.H}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/tree [get]
func (m *MenuApi) GetMenuTree(c *gin.Context) {
	var menus []models.Permission
	if err := global.DB.Where("type in (1, 2)").Order("sort").Find(&menus).Error; err != nil {
		global.Logger.Error("获取菜单树失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	// 构建菜单树
	tree := buildMenuTree(menus, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}

// buildMenuTree 构建菜单树结构
func buildMenuTree(menus []models.Permission, parentID uint) []gin.H {
	tree := make([]gin.H, 0)

	for _, menu := range menus {
		if menu.ParentID == parentID {
			children := buildMenuTree(menus, menu.ID)
			node := gin.H{
				"id":       menu.ID,
				"name":     menu.Name,
				"path":     menu.Path,
				"component": menu.Component,
				"icon":     menu.Icon,
				"type":     menu.Type,
				"sort":     menu.Sort,
				"children": children,
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// GetUserMenus 获取用户菜单
// @Summary 获取用户菜单接口
// @Description 查询指定用户的菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]gin.H{"id":int, "name":string, "children":[]gin.H}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/menu/user-menus [get]
func (m *MenuApi) GetUserMenus(c *gin.Context) {
	// 从token中获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Error("获取用户菜单失败: 用户未登录")
		c.JSON(401, gin.H{"code": 401, "msg": "用户未登录"})
		return
	}

	// 如果是超级管理员，返回所有菜单
	var user models.User
	global.DB.First(&user, userID)
	if user.IsAdmin {
		var menus []models.Permission
		if err := global.DB.Where("type in (1, 2)").Order("sort").Find(&menus).Error; err != nil {
			global.Logger.Error("获取用户菜单失败: " + err.Error())
			c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
			return
		}

		tree := buildMenuTree(menus, 0)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "获取成功",
			"data": tree,
		})
		return
	}

	// 普通用户，根据权限获取菜单
	var permissions []models.Permission
	if err := global.DB.Table("permissions").
		Select("permissions.*").
		Joins("join role_permissions on permissions.id = role_permissions.permission_id").
		Joins("join user_roles on role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? and permissions.type in (1, 2)", userID).
		Order("permissions.sort").
		Distinct().
		Find(&permissions).Error; err != nil {
		global.Logger.Error("获取用户菜单失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	// 构建菜单树
	tree := buildMenuTree(permissions, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}