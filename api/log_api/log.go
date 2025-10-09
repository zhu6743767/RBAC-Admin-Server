package log_api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"rbac_admin_server/global"
	"rbac_admin_server/models"
)

// GetLogList 获取日志列表
// @Summary 获取日志列表接口
// @Description 查询系统中的日志列表
// @Tags 日志管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param level query string false "日志级别"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"list":[]models.Log, "total":int}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/log/list [get]
func (l *LogApi) GetLogList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	level := c.Query("level")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// 构建查询条件
	query := global.DB.Model(&models.Log{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if startTime != "" && endTime != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var logs []models.Log
	if err := query.
		Order("created_at DESC").
		Offset((pageInt - 1) * pageSizeInt).
		Limit(pageSizeInt).
		Find(&logs).Error; err != nil {
		global.Logger.Error("获取日志列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取日志列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  logs,
			"total": total,
		},
	})
}

// GetUserLogs 获取用户日志
// @Summary 获取用户日志接口
// @Description 查询指定用户的操作日志
// @Tags 日志管理
// @Accept json
// @Produce json
// @Param user_id query int true "用户ID"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"list":[]models.Log, "total":int}}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/log/user-logs [get]
func (l *LogApi) GetUserLogs(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		global.Logger.Error("获取用户日志失败: 用户ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// 构建查询条件
	query := global.DB.Model(&models.Log{}).Where("user_id = ?", userID)

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var logs []models.Log
	if err := query.
		Order("created_at DESC").
		Offset((pageInt - 1) * pageSizeInt).
		Limit(pageSizeInt).
		Find(&logs).Error; err != nil {
		global.Logger.Error("获取用户日志失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取用户日志失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  logs,
			"total": total,
		},
	})
}

// DeleteLog 删除日志
// @Summary 删除日志接口
// @Description 根据日志ID删除日志
// @Tags 日志管理
// @Accept json
// @Produce json
// @Param id query int true "日志ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/log/delete [delete]
func (l *LogApi) DeleteLog(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		global.Logger.Error("删除日志失败: ID为空")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.Log{}, id).Error; err != nil {
		global.Logger.Error("删除日志失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("删除日志成功: ID=%s", id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// DeleteMultipleLogs 批量删除日志
// @Summary 批量删除日志接口
// @Description 批量删除指定的日志
// @Tags 日志管理
// @Accept json
// @Produce json
// @Param ids body []int true "日志ID列表"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/log/delete-multiple [delete]
func (l *LogApi) DeleteMultipleLogs(c *gin.Context) {
	var req struct {
		Ids []int `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		global.Logger.Error("批量删除日志参数错误: " + err.Error())
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := global.DB.Delete(&models.Log{}, req.Ids).Error; err != nil {
		global.Logger.Error("批量删除日志失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("批量删除日志成功: 共%d条", len(req.Ids))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// GetLogDashboard 获取日志仪表盘数据
// @Summary 获取日志仪表盘数据接口
// @Description 获取日志统计数据
// @Tags 日志管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"total":int, "levels":map[string]int, "recent":[]gin.H}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/log/dashboard [get]
func (l *LogApi) GetLogDashboard(c *gin.Context) {
	// 获取总日志数
	var total int64
	global.DB.Model(&models.Log{}).Count(&total)

	// 获取各级别日志数量
	var levelStats []struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	global.DB.Model(&models.Log{}).
		Select("level, count(*) as count").
		Group("level").
		Scan(&levelStats)

	levelMap := make(map[string]int64)
	for _, stat := range levelStats {
		levelMap[stat.Level] = stat.Count
	}

	// 获取最近7天的日志数量
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var dailyStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	global.DB.Model(&models.Log{}).
		Select("date(created_at) as date, count(*) as count").
		Where("created_at >= ?", sevenDaysAgo).
		Group("date").
		Order("date").
		Scan(&dailyStats)

	// 获取最近的10条日志
	var recentLogs []models.Log
	global.DB.Model(&models.Log{}).
		Order("created_at DESC").
		Limit(10).
		Find(&recentLogs)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"total":  total,
			"levels": levelMap,
			"daily":  dailyStats,
			"recent": recentLogs,
		},
	})
}