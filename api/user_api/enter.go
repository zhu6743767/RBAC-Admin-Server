package user_api

import "github.com/gin-gonic/gin"

// UserApi 用户API结构体
// 提供用户相关的所有API接口
// 包括登录、注册、用户管理等功能
// 采用结构体方式组织，便于依赖注入和测试

type UserApi struct {}

// NewUserApi 创建用户API实例
// 返回一个新的UserApi结构体指针
func NewUserApi() *UserApi {
	return &UserApi{}
}

// RegisterRoutes 注册用户API路由
func (u *UserApi) RegisterRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/user")
	{
		userRouter.GET("/list", u.GetUserList)
		userRouter.POST("/create", u.CreateUser)
		userRouter.PUT("/update", u.UpdateUser)
		userRouter.DELETE("/delete", u.DeleteUser)
	}
}