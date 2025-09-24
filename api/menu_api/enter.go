package menu_api

import "github.com/gin-gonic/gin"

// MenuApi 菜单API结构体
type MenuApi struct{}

// NewMenuApi 创建菜单API实例
func NewMenuApi() *MenuApi {
	return &MenuApi{}
}

// RegisterRoutes 注册菜单API路由
func (m *MenuApi) RegisterRoutes(router *gin.RouterGroup) {
	menuRouter := router.Group("/menu")
	{
		menuRouter.GET("/list", m.GetMenuList)
		menuRouter.POST("/create", m.CreateMenu)
		menuRouter.PUT("/update", m.UpdateMenu)
		menuRouter.DELETE("/delete", m.DeleteMenu)
		menuRouter.GET("/tree", m.GetMenuTree)
		menuRouter.GET("/user-menus", m.GetUserMenus)
	}
}