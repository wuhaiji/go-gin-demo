package routers

import (
	"encoding/gob"
	"gin_demo/app/auth"
	"gin_demo/app/entity"
	"gin_demo/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Option func(engine *gin.Engine)

func (o Option) apply(logger *zap.Logger) {
	//TODO implement me
	panic("implement me")
}

var options []Option

func Include(optionParams ...Option) {
	options = optionParams
}
func Init() *gin.Engine {
	engine := gin.New()

	engine.Use(logger.GinLogger())
	// 异常处理
	engine.Use(logger.RecoveryWithZap(false))

	// 注册User结构体
	gob.Register(entity.UserInfo{})
	// 创建基于cookie的存储引擎，secret11111 参数是用于加密的密钥
	store := cookie.NewStore([]byte("fjfjffkd"))
	// 设置session中间件，参数gin-cookie，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	store.Options(sessions.Options{
		MaxAge: 20 * 60, //30min
		Path:   "/",
	})
	engine.Use(sessions.Sessions("SESSION_ID", store))
	engine.Use(authMiddleware())
	for _, option := range options {
		option(engine)
	}
	return engine
}

func authMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path == "/auth" {
			context.Next()
		}
		currentUser := auth.GetCurrentUser(context)
		if (currentUser == entity.UserInfo{}) {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			context.Abort()
			return
		} else {
			context.Next()
		}
	}
}
