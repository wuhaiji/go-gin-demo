package login

import (
	"gin_demo/util/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserInfo struct {
	Password string `json:"password" form:"password"`
	Username string `json:"username" form:"username"`
}

func authHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.Username == "whj" && user.Password == "123qwe" {
		// 生成Token
		tokenString, err := jwt.GenToken(user.Username)
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "success",
			"data": gin.H{"token": tokenString},
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}
