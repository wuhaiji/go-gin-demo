package shop

import (
	"gin_demo/app/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	//panic(NewBizError("测试异常"))
	zap.L().Info("测试log")
	context.JSON(http.StatusOK, gin.H{
		"message": currentUser.Username + ":goodsHandler",
	})
}
func checkoutHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "checkoutHandler",
	})
}
