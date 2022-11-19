package main

import (
	"bytes"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"time"
)

// Login Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	r.Use(ginBodyLogMiddleware())
	r.Use(statCostMiddleware())
	r.Use(CustomExceptionHandler())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},                         // 允许跨域发来请求的网站
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool { // 自定义过滤源站的方法
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.POST("/login/json", func(c *gin.Context) {
		var login Login
		err := c.ShouldBind(&login)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Printf("login info:%#v\n", login)
		c.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"password": login.Password,
		})
	})

	r.POST("/login/form", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	// 绑定QueryString示例 (/loginQuery?user=q1mi&password=123456)
	r.GET("/login/querystring", func(c *gin.Context) {
		var login Login
		// ShouldBind()会根据请求的Content-Type自行选择绑定器
		if err := c.ShouldBind(&login); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"user":     login.User,
				"password": login.Password,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
	r.POST("/upload", func(ctx *gin.Context) {
		var form, err = ctx.MultipartForm()
		if err != nil {
			panic(err)
		}
		files := form.File["file"]
		hashMap := make(map[string]string)
		for _, file := range files {
			log.Println(file.Filename)
			var directory = "c:/tmp/"
			dst := fmt.Sprintf(directory+"%s", file.Filename)
			_, err := os.Stat(directory)
			if err != nil {
				if os.IsNotExist(err) {
					err := os.Mkdir(directory, 777)
					if err != nil {
						panic(err)
					}
					err = os.Chmod(directory, 777)
					if err != nil {
						panic(err)
					}
				}
			}

			err = ctx.SaveUploadedFile(file, dst)
			if err != nil {
				panic(err)
			}
			hashMap[file.Filename] = dst
		}
		ctx.JSON(http.StatusOK, hashMap)
	})

	//重定向
	r.GET("/test", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com/")
	})

	// 路由重定向
	r.GET("/hello", func(c *gin.Context) {
		panic(gin.H{
			"message": "测试异常",
		})
		c.String(http.StatusOK, "hello world")
	})
	r.GET("/router/redirect", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/hello"

		r.HandleContext(ctx)
	})
	shopGroup := r.Group("/shop")
	{
		shopGroup.GET("/cart", func(context *gin.Context) {
			StringOk(context, "shop cart")
		})
	}

	err := r.Run("localhost:8080")
	_, err = fmt.Fprint(gin.DefaultWriter, "app start on http://localhost:8080")
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
}
func statCostMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()
		context.Set("name", "123")
		context.Next()
		cost := time.Since(startTime)
		log.Printf("花费时间：%s\n", cost)
	}
}

var g errgroup.Group

func ginBodyLogMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		blw := &bodyLogWriter{
			body:           bytes.NewBuffer([]byte{}),
			ResponseWriter: context.Writer,
		}
		context.Writer = blw
		context.Next()
		fmt.Println("Response body:" + blw.body.String())
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter               // 嵌入gin框架ResponseWriter
	body               *bytes.Buffer // 我们记录用的response
}

// Write 写入响应体数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // 我们记录一份
	return w.ResponseWriter.Write(b) // 真正写入响应
}
