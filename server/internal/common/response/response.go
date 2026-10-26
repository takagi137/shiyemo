package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// PageData 分页数据结构
type PageData struct {
	List    interface{} `json:"list"`
	Page    int         `json:"page"`
	Size    int         `json:"size"`
	HasMore bool        `json:"has_more"`
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, page, size int, total int64) {
	hasMore := int64(page*size) < total
	Success(c, PageData{
		List:    list,
		Page:    page,
		Size:    size,
		HasMore: hasMore,
	})
}
