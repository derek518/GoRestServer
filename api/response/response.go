package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"GoRestServer/helper"
	"GoRestServer/model/base"
)

func Succeed(context *gin.Context, statusCode int, content interface{}) {
	message := ""
	if statusCode > 0 {
		message = helper.StatusText(statusCode)
	}

	context.JSON(http.StatusOK,
		&base.JsonObject{
			Code:    0,
			Message: message,
			Content: content,
		})
}

func Failed(context *gin.Context, httpCode int, statusCode int, err error) {
	context.JSON(httpCode, base.JsonObject{
		Code:    1,
		Message: helper.StatusText(statusCode),
		Content: err,
	})
}

func FailedWithOK(context *gin.Context, statusCode int, err error) {
	context.JSON(http.StatusOK, base.JsonObject{
		Code:    1,
		Message: helper.StatusText(statusCode),
		Content: err,
	})
}
