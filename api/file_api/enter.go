package file_api

import "github.com/gin-gonic/gin"

// FileApi 文件API结构体
type FileApi struct{}

// NewFileApi 创建文件API实例
func NewFileApi() *FileApi {
	return &FileApi{}
}

// RegisterRoutes 注册文件API路由
func (f *FileApi) RegisterRoutes(router *gin.RouterGroup) {
	fileRouter := router.Group("/file")
	{
		fileRouter.POST("/upload", f.UploadFile)
		fileRouter.POST("/upload-multiple", f.UploadMultipleFiles)
		fileRouter.GET("/download/:id", f.DownloadFile)
		fileRouter.GET("/list", f.GetFileList)
		fileRouter.DELETE("/delete/:id", f.DeleteFile)
	}
}