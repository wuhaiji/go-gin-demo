package shop

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func goodsHandler(context *gin.Context) {
	time.Sleep(time.Second * 5)
	context.JSON(http.StatusOK, gin.H{
		"message": "goodsHandler",
	})
}
func checkoutHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "checkoutHandler",
	})
}
