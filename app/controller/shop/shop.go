package shop

import (
	"gin_demo/app/controller/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BizError struct {
	message string
}

func (e BizError) Error() string {
	return e.message
}

func NewBizError(msg string) BizError {
	return BizError{
		message: msg,
	}
}

func goodsHandler(context *gin.Context) {
	currentUser := auth.GetCurrentUser(context)
	panic(NewBizError("测试异常"))
	context.JSON(http.StatusOK, gin.H{
		"message": currentUser.Username + ":goodsHandler",
	})
}
func checkoutHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "checkoutHandler",
	})
}
