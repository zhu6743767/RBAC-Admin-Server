package dept_api

import (
	"rbac_admin_server/global"
	"rbac_admin_server/models"

	"github.com/gin-gonic/gin"
)

// GetDepartmentList 获取部门列表
// @Summary 获取部门列表接口
// @Description 查询系统中的部门列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.Department}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/list [get]
func (d *DepartmentApi) GetDepartmentList(c *gin.Context) {
	var departments []models.Department
	if err := global.DB.Find(&departments).Error; err != nil {
		global.Logger.Error("获取部门列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取部门列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": departments,
	})
}

// CreateDepartment 创建部门
// @Summary 创建部门接口
// @Description 管理员创建新部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param department body models.Department true "部门信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/create [post]
func (d *DepartmentApi) CreateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		global.Logger.Error("创建部门参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查部门名称是否已存在
	var count int64
	global.DB.Model(&models.Department{}).Where("name = ?", department.Name).Count(&count)
	if count > 0 {
		global.Logger.Error("部门名称已存在: " + department.Name)
		c.JSON(400, gin.H{"code": 400, "msg": "部门名称已存在"})
		return
	}

	if err := global.DB.Create(&department).Error; err != nil {
		global.Logger.Error("创建部门失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	global.Logger.Infof("管理员创建部门成功: %s", department.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// UpdateDepartment 更新部门
// @Summary 更新部门接口
// @Description 管理员更新部门信息
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param department body models.Department true "更新的部门信息"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/update [put]
func (d *DepartmentApi) UpdateDepartment(c *gin.Context) {
	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		global.Logger.Error("更新部门参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if department.ID == 0 {
		global.Logger.Error("更新部门参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Save(&department).Error; err != nil {
		global.Logger.Error("更新部门失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	global.Logger.Infof("管理员更新部门成功: %s", department.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteDepartment 删除部门
// @Summary 删除部门接口
// @Description 管理员删除部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param id query int true "部门ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/delete [delete]
func (d *DepartmentApi) DeleteDepartment(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除部门参数错误: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 检查部门是否有子部门
	var childDepartments []models.Department
	global.DB.Where("parent_id = ?", id).Find(&childDepartments)
	if len(childDepartments) > 0 {
		global.Logger.Error("删除部门失败: 部门有子部门")
		c.JSON(400, gin.H{"code": 400, "msg": "部门有子部门，无法删除"})
		return
	}

	// 检查部门是否有用户
	var users []models.User
	global.DB.Where("department_id = ?", id).Find(&users)
	if len(users) > 0 {
		global.Logger.Error("删除部门失败: 部门有用户")
		c.JSON(400, gin.H{"code": 400, "msg": "部门有用户，无法删除"})
		return
	}

	if err := global.DB.Delete(&models.Department{}, id).Error; err != nil {
		global.Logger.Error("删除部门失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("管理员删除部门成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetDepartmentTree 获取部门树结构
// @Summary 获取部门树结构接口
// @Description 查询系统中的部门树结构
// @Tags 部门管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]gin.H{"id":int, "name":string, "children":[]gin.H}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/tree [get]
func (d *DepartmentApi) GetDepartmentTree(c *gin.Context) {
	var departments []models.Department
	if err := global.DB.Order("sort").Find(&departments).Error; err != nil {
		global.Logger.Error("获取部门树失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	// 构建部门树
	tree := buildDepartmentTree(departments, 0)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": tree,
	})
}

// buildDepartmentTree 构建部门树结构
func buildDepartmentTree(departments []models.Department, parentID uint) []gin.H {
	tree := make([]gin.H, 0)

	for _, dept := range departments {
		if dept.ParentID == parentID {
			children := buildDepartmentTree(departments, dept.ID)
			node := gin.H{
				"id":       dept.ID,
				"name":     dept.Name,
				"parent_id": dept.ParentID,
				"sort":     dept.Sort,
				"children": children,
			}
			tree = append(tree, node)
		}
	}

	return tree
}

// GetDepartmentUsers 获取部门用户
// @Summary 获取部门用户接口
// @Description 查询指定部门下的用户列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Param dept_id query int true "部门ID"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.User}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /api/admin/dept/users [get]
func (d *DepartmentApi) GetDepartmentUsers(c *gin.Context) {
	deptID := c.Query("dept_id")
	if deptID == "" {
		global.Logger.Error("获取部门用户参数错误: 部门ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	var users []models.User
	if err := global.DB.Where("department_id = ?", deptID).Find(&users).Error; err != nil {
		global.Logger.Error("获取部门用户失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": users,
	})
}