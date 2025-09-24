package dept_api

import "github.com/gin-gonic/gin"

// DepartmentApi 部门API结构体
type DepartmentApi struct{}

// NewDepartmentApi 创建部门API实例
func NewDepartmentApi() *DepartmentApi {
	return &DepartmentApi{}
}

// RegisterRoutes 注册部门API路由
func (d *DepartmentApi) RegisterRoutes(router *gin.RouterGroup) {
	deptRouter := router.Group("/dept")
	{
		deptRouter.GET("/list", d.GetDepartmentList)
		deptRouter.POST("/create", d.CreateDepartment)
		deptRouter.PUT("/update", d.UpdateDepartment)
		deptRouter.DELETE("/delete", d.DeleteDepartment)
		deptRouter.GET("/tree", d.GetDepartmentTree)
		deptRouter.GET("/users", d.GetDepartmentUsers)
	}
}