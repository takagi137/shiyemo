package upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bizerr "github.com/shiyemo/server/internal/common/errors"
	"github.com/shiyemo/server/internal/common/response"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	service *UploadService
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		service: NewUploadService(),
	}
}

// NewUploadHandlerWithService 创建上传处理器（指定 service）
func NewUploadHandlerWithService(service *UploadService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	// 获取当前用户 ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "请选择文件")
		return
	}
	defer file.Close()

	// 获取分类（可选，默认为 general）
	category := c.DefaultPostForm("category", "general")

	// 上传文件
	result, err := h.service.Upload(currentUserID.(int64), file, header.Filename, category)
	if err != nil {
		bizErr, ok := err.(*bizerr.BizError)
		if ok {
			response.Error(c, http.StatusBadRequest, bizErr.Code, bizErr.Message)
		} else {
			response.Error(c, http.StatusInternalServerError, bizerr.ErrInternal.Code, bizerr.ErrInternal.Message)
		}
		return
	}

	response.Success(c, result)
}

// UploadMultiple 批量上传图片
func (h *UploadHandler) UploadMultiple(c *gin.Context) {
	// 获取当前用户 ID
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, bizerr.ErrUnauthorized.Code, bizerr.ErrUnauthorized.Message)
		return
	}

	// 获取多个文件
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "请上传文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "请选择文件")
		return
	}

	// 限制最多 9 张图片
	if len(files) > 9 {
		response.Error(c, http.StatusBadRequest, bizerr.ErrParamInvalid.Code, "最多上传 9 张图片")
		return
	}

	// 获取分类（可选，默认为 general）
	category := c.DefaultPostForm("category", "general")

	// 上传所有文件
	var results []UploadResult
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		result, err := h.service.Upload(currentUserID.(int64), file, header.Filename, category)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}

	response.Success(c, gin.H{
		"files": results,
		"count": len(results),
	})
}
