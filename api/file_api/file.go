package file_api

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rbac_admin_server/global"
	"rbac_admin_server/models"
)

// UploadFile 上传单个文件
// @Summary 上传单个文件接口
// @Description 上传单个文件到服务器
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type query string false "文件类型"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":models.File}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/file/upload [post]
func (f *FileApi) UploadFile(c *gin.Context) {
	// 单个文件
	file, _ := c.FormFile("file")
	if file == nil {
		global.Logger.Error("上传文件失败: 文件为空")
		c.JSON(400, gin.H{"code": 400, "msg": "文件不能为空"})
		return
	}

	// 创建上传目录
	uploadDir := "uploads/"
	dir := uploadDir + time.Now().Format("2006-01-02") + "/"
	if err := os.MkdirAll(dir, 0755); err != nil {
		global.Logger.Error("创建上传目录失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建目录失败"})
		return
	}

	// 生成新的文件名
	ext := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dst := dir + newFileName

	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		global.Logger.Error("保存文件失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "保存文件失败"})
		return
	}

	// 保存文件信息到数据库
	fileModel := models.File{
		Name:      file.Filename,
		Path:      dst,
		Size:      file.Size,
		Type:      c.Query("type"),
		Extension: ext,
	}

	if err := global.DB.Create(&fileModel).Error; err != nil {
		global.Logger.Error("保存文件信息失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "保存文件信息失败"})
		return
	}

	global.Logger.Infof("文件上传成功: %s", file.Filename)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "上传成功",
		"data": fileModel,
	})
}

// UploadMultipleFiles 上传多个文件
// @Summary 上传多个文件接口
// @Description 上传多个文件到服务器
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "多个文件"
// @Param type query string false "文件类型"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":[]models.File}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/file/upload-multiple [post]
func (f *FileApi) UploadMultipleFiles(c *gin.Context) {
	// 多文件
	form, _ := c.MultipartForm()
	files := form.File["files"]
	if len(files) == 0 {
		global.Logger.Error("上传文件失败: 文件为空")
		c.JSON(400, gin.H{"code": 400, "msg": "文件不能为空"})
		return
	}

	// 创建上传目录
	uploadDir := "uploads/"
	dir := uploadDir + time.Now().Format("2006-01-02") + "/"
	if err := os.MkdirAll(dir, 0755); err != nil {
		global.Logger.Error("创建上传目录失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "创建目录失败"})
		return
	}

	// 保存多个文件
	fileModels := make([]models.File, 0)
	for _, file := range files {
		// 生成新的文件名
		ext := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + ext
		dst := dir + newFileName

		// 保存文件
		if err := c.SaveUploadedFile(file, dst); err != nil {
			global.Logger.Error("保存文件失败: " + err.Error())
			c.JSON(500, gin.H{"code": 500, "msg": "保存文件失败"})
			return
		}

		// 保存文件信息到数据库
		fileModel := models.File{
			Name:      file.Filename,
			Path:      dst,
			Size:      file.Size,
			Type:      c.Query("type"),
			Extension: ext,
		}
		fileModels = append(fileModels, fileModel)
	}

	// 批量保存文件信息
	if err := global.DB.CreateInBatches(fileModels, len(fileModels)).Error; err != nil {
		global.Logger.Error("保存文件信息失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "保存文件信息失败"})
		return
	}

	global.Logger.Infof("多文件上传成功: 共%d个文件", len(files))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("上传成功，共%d个文件", len(files)),
		"data": fileModels,
	})
}

// DownloadFile 下载文件
// @Summary 下载文件接口
// @Description 根据文件ID下载文件
// @Tags 文件管理
// @Accept json
// @Produce application/octet-stream
// @Param id path int true "文件ID"
// @Success 200 {file} file
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 404 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/file/download/{id} [get]
func (f *FileApi) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		global.Logger.Error("下载文件失败: ID格式错误")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 查询文件信息
	var fileModel models.File
	if err := global.DB.First(&fileModel, id).Error; err != nil {
		global.Logger.Error("下载文件失败: 文件不存在")
		c.JSON(404, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fileModel.Path); os.IsNotExist(err) {
		global.Logger.Error("下载文件失败: 文件不存在")
		c.JSON(404, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	global.Logger.Infof("文件下载成功: %s", fileModel.Name)
	c.File(fileModel.Path)
}

// GetFileList 获取文件列表
// @Summary 获取文件列表接口
// @Description 查询系统中的文件列表
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param type query string false "文件类型"
// @Success 200 {object} gin.H{"code":int, "msg":string, "data":gin.H{"list":[]models.File, "total":int}}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/file/list [get]
func (f *FileApi) GetFileList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	fileType := c.Query("type")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	// 构建查询条件
	query := global.DB.Model(&models.File{})
	if fileType != "" {
		query = query.Where("type = ?", fileType)
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 查询列表
	var files []models.File
	if err := query.
		Order("created_at DESC").
		Offset((pageInt - 1) * pageSizeInt).
		Limit(pageSizeInt).
		Find(&files).Error; err != nil {
		global.Logger.Error("获取文件列表失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "获取文件列表失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  files,
			"total": total,
		},
	})
}

// DeleteFile 删除文件
// @Summary 删除文件接口
// @Description 根据文件ID删除文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param id path int true "文件ID"
// @Success 200 {object} gin.H{"code":int, "msg":string}
// @Failure 400 {object} gin.H{"code":int, "msg":string}
// @Failure 404 {object} gin.H{"code":int, "msg":string}
// @Failure 500 {object} gin.H{"code":int, "msg":string}
// @Router /admin/file/delete/{id} [delete]
func (f *FileApi) DeleteFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		global.Logger.Error("删除文件失败: ID格式错误")
		c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 查询文件信息
	var fileModel models.File
	if err := global.DB.First(&fileModel, id).Error; err != nil {
		global.Logger.Error("删除文件失败: 文件不存在")
		c.JSON(404, gin.H{"code": 404, "msg": "文件不存在"})
		return
	}

	// 删除物理文件
	if err := os.Remove(fileModel.Path); err != nil {
		global.Logger.Error("删除物理文件失败: " + err.Error())
		// 继续执行，删除数据库记录
	}

	// 删除数据库记录
	if err := global.DB.Delete(&fileModel).Error; err != nil {
		global.Logger.Error("删除文件记录失败: " + err.Error())
		c.JSON(500, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	global.Logger.Infof("文件删除成功: %s", fileModel.Name)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}