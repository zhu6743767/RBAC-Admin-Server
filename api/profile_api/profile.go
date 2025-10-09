package profile_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"rbac_admin_server/global"
	"rbac_admin_server/models"
	"rbac_admin_server/utils"
)

// GetUserInfoRequest 获取用户信息请求参数
// swagger:model GetUserInfoRequest
type GetUserInfoRequest struct{}

// UserInfoResponse 用户信息响应
// swagger:model UserInfoResponse
type UserInfoResponse struct {
	// 用户ID
	ID uint `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 昵称
	Nickname string `json:"nickname"`
	// 邮箱
	Email string `json:"email"`
	// 手机号
	Phone string `json:"phone"`
	// 头像
	Avatar string `json:"avatar"`
	// 性别
	Gender int `json:"gender"`
	// 状态
	Status int `json:"status"`
	// 创建时间
	CreatedAt time.Time `json:"created_at"`
	// 角色列表
	Roles []string `json:"roles"`
	// 部门信息
	Department *struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"department,omitempty"`
}

// GetUserInfo 获取用户个人信息
// @Summary 获取用户个人信息
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} utils.Response{data=UserInfoResponse}
// @Router /profile/info [get]
func (p *ProfileApi) GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	var user models.User
	result := global.DB.Preload("Roles").First(&user, userID)
	if result.Error != nil {
		utils.Error(c, http.StatusInternalServerError, utils.ERROR_GET_USER, nil)
		return
	}

	// 构建响应数据
	resp := UserInfoResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	// 填充角色信息
	for _, role := range user.Roles {
		resp.Roles = append(resp.Roles, role.Name)
	}

	// 填充部门信息
	if user.DeptID > 0 {
		var dept models.Department
		if err := global.DB.First(&dept, user.DeptID).Error; err == nil {
			resp.Department = &struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{dept.ID, dept.Name}
		}
	}

	utils.Success(c, resp)
}

// UpdateUserInfoRequest 更新用户信息请求参数
// swagger:model UpdateUserInfoRequest
type UpdateUserInfoRequest struct {
	// 昵称
	Nickname string `json:"nickname" validate:"max=50"`
	// 邮箱
	Email string `json:"email" validate:"email,max=100"`
	// 手机号
	Phone string `json:"phone" validate:"max=20"`
	// 头像
	Avatar string `json:"avatar" validate:"max=255"`
	// 性别 (0:未知, 1:男, 2:女)
	Gender int `json:"gender" validate:"min=0,max=2"`
}

// UpdateUserInfo 更新用户个人信息
// @Summary 更新用户个人信息
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body UpdateUserInfoRequest true "用户信息"
// @Success 200 {object} utils.Response{data=string}
// @Router /profile/info [put]
func (p *ProfileApi) UpdateUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	// 绑定请求参数
	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, utils.GetValidationError(errs))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, nil)
		return
	}

	// 更新用户信息
	result := global.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"avatar":   req.Avatar,
		"gender":   req.Gender,
	})

	if result.Error != nil {
		utils.Error(c, http.StatusInternalServerError, utils.ERROR_UPDATE_USER, nil)
		return
	}

	utils.Success(c, "更新成功")
}

// UpdatePasswordRequest 修改密码请求参数
// swagger:model UpdatePasswordRequest
type UpdatePasswordRequest struct {
	// 当前密码
	CurrentPassword string `json:"current_password" validate:"required,min=6,max=32"`
	// 新密码
	NewPassword string `json:"new_password" validate:"required,min=6,max=32"`
	// 确认新密码
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UpdatePassword 修改密码
// @Summary 修改密码
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body UpdatePasswordRequest true "密码信息"
// @Success 200 {object} utils.Response{data=string}
// @Router /profile/password [put]
func (p *ProfileApi) UpdatePassword(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	// 绑定请求参数
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, utils.GetValidationError(errs))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, nil)
		return
	}

	// 查询用户
	var user models.User
	result := global.DB.First(&user, userID)
	if result.Error != nil {
		utils.Error(c, http.StatusInternalServerError, utils.ERROR_GET_USER, nil)
		return
	}

	// 验证当前密码
	if !utils.ComparePassword(req.CurrentPassword, user.Password) {
		utils.Error(c, http.StatusBadRequest, utils.ERROR_PASSWORD_WRONG, nil)
		return
	}

	// 加密新密码
	passwordHash := utils.HashedPassword(req.NewPassword)
	if passwordHash == "" {
		utils.Error(c, http.StatusInternalServerError, utils.ERROR_ENCRYPT_PASSWORD, nil)
		return
	}

	// 更新密码
	result = global.DB.Model(&user).Update("password", passwordHash)
	if result.Error != nil {
		utils.Error(c, http.StatusInternalServerError, utils.ERROR_UPDATE_USER, nil)
		return
	}

	utils.Success(c, "密码修改成功")
}

// DashboardDataResponse 仪表盘数据响应
// swagger:model DashboardDataResponse
type DashboardDataResponse struct {
	// 用户总数
	TotalUsers int `json:"total_users"`
	// 角色总数
	TotalRoles int `json:"total_roles"`
	// 部门总数
	TotalDepartments int `json:"total_departments"`
	// 菜单总数
	TotalMenus int `json:"total_menus"`
	// 权限总数
	TotalPermissions int `json:"total_permissions"`
	// 文件总数
	TotalFiles int `json:"total_files"`
	// 今日登录次数
	TodayLogins int `json:"today_logins"`
	// 系统运行时间（秒）
	SystemUptime int64 `json:"system_uptime"`
}

// GetDashboardData 获取仪表盘数据
// @Summary 获取仪表盘数据
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} utils.Response{data=DashboardDataResponse}
// @Router /profile/dashboard [get]
func (p *ProfileApi) GetDashboardData(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	var dashboardData DashboardDataResponse

	// 获取用户总数
	var totalUsers int64
	global.DB.Model(&models.User{}).Count(&totalUsers)
	dashboardData.TotalUsers = int(totalUsers)

	// 获取角色总数
	var totalRoles int64
	global.DB.Model(&models.Role{}).Count(&totalRoles)
	dashboardData.TotalRoles = int(totalRoles)

	// 获取部门总数
	var totalDepartments int64
	global.DB.Model(&models.Department{}).Count(&totalDepartments)
	dashboardData.TotalDepartments = int(totalDepartments)

	// 获取菜单总数
	var totalMenus int64
	global.DB.Model(&models.Menu{}).Count(&totalMenus)
	dashboardData.TotalMenus = int(totalMenus)

	// 获取权限总数
	var totalPermissions int64
	global.DB.Model(&models.Permission{}).Count(&totalPermissions)
	dashboardData.TotalPermissions = int(totalPermissions)

	// 获取文件总数
	var totalFiles int64
	global.DB.Model(&models.File{}).Count(&totalFiles)
	dashboardData.TotalFiles = int(totalFiles)

	// 获取今日登录次数
	var todayLogins int64
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	global.DB.Model(&models.Log{}).Where("created_at >= ? AND action = ?", todayStart, "login").Count(&todayLogins)
	dashboardData.TodayLogins = int(todayLogins)

	// 计算系统运行时间（简化处理，返回0）
	dashboardData.SystemUptime = 0

	utils.Success(c, dashboardData)
}

// UserSettingsResponse 用户设置响应
// swagger:model UserSettingsResponse
type UserSettingsResponse struct {
	// 主题设置 (light, dark, auto)
	Theme string `json:"theme"`
	// 语言设置 (zh-CN, en-US)
	Language string `json:"language"`
	// 布局设置
	Layout string `json:"layout"`
	// 侧边栏折叠状态
	SidebarCollapsed bool `json:"sidebar_collapsed"`
	// 通知设置
	Notifications map[string]bool `json:"notifications"`
}

// GetUserSettings 获取用户设置
// @Summary 获取用户设置
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} utils.Response{data=UserSettingsResponse}
// @Router /profile/settings [get]
func (p *ProfileApi) GetUserSettings(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	// 这里简化处理，实际应该从数据库或缓存获取用户设置
	// 为了演示，返回默认设置
	settings := UserSettingsResponse{
		Theme:            "light",
		Language:         "zh-CN",
		Layout:           "default",
		SidebarCollapsed: false,
		Notifications: map[string]bool{
			"email":  true,
			"sms":    false,
			"push":   true,
			"system": true,
		},
	}

	// 实际项目中，应该查询用户设置表
	// var userSettings models.UserSettings
	// err := global.DB.Where("user_id = ?", userID).First(&userSettings).Error
	// if err == nil {
	//     // 填充设置数据
	// }

	utils.Success(c, settings)
}

// UpdateUserSettingsRequest 更新用户设置请求参数
// swagger:model UpdateUserSettingsRequest
type UpdateUserSettingsRequest struct {
	// 主题设置 (light, dark, auto)
	Theme string `json:"theme" validate:"oneof=light dark auto"`
	// 语言设置 (zh-CN, en-US)
	Language string `json:"language" validate:"oneof=zh-CN en-US"`
	// 布局设置
	Layout string `json:"layout"`
	// 侧边栏折叠状态
	SidebarCollapsed bool `json:"sidebar_collapsed"`
	// 通知设置
	Notifications map[string]bool `json:"notifications"`
}

// UpdateUserSettings 更新用户设置
// @Summary 更新用户设置
// @Tags 个人信息管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param data body UpdateUserSettingsRequest true "设置信息"
// @Success 200 {object} utils.Response{data=string}
// @Router /profile/settings [put]
func (p *ProfileApi) UpdateUserSettings(c *gin.Context) {
	// 从上下文获取用户ID
	_, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, utils.ERROR_UNAUTHORIZED, nil)
		return
	}

	// 绑定请求参数
	var req UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, utils.GetValidationError(errs))
			return
		}
		utils.Error(c, http.StatusBadRequest, utils.ERROR_INVALID_PARAM, nil)
		return
	}

	// 这里简化处理，实际应该更新数据库或缓存中的用户设置
	// 为了演示，直接返回成功

	// 实际项目中，应该更新用户设置表
	// var userSettings models.UserSettings
	// err := global.DB.Where("user_id = ?", userID).First(&userSettings).Error
	// if err == gorm.ErrRecordNotFound {
	//     // 创建新的设置记录
	// } else if err == nil {
	//     // 更新现有设置
	// }

	utils.Success(c, "设置更新成功")
}