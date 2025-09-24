package user_api

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