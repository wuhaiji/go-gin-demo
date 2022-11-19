package blog

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func postHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "postHandler",
	})
}
