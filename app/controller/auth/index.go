package auth

import (
	"errors"
	"gin_demo/app/entity"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"
)

func authHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user entity.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	//

	if user.Username == "whj" && user.Password == "123qwe" {

		SetCurrentUser(c, entity.UserInfo{Password: "whj", Username: "123qwe"}) // 邮箱和密码正确则将当前用户信息写入session中
		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
		})
		// 生成Token
		//tokenString, err := jwt.GenToken(user.Username)
		//if err != nil {
		//	panic(err)
		//}
		//c.JSON(http.StatusOK, gin.H{
		//	"code": 2000,
		//	"msg":  "success",
		//	"data": gin.H{"token": tokenString},
		//})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "登录失败",
	})
	return
}

const currentUser = "currentUser"

func GetCurrentUser(c *gin.Context) (user entity.UserInfo) {

	session := sessions.Default(c)
	get := session.Get(currentUser)
	if get == nil {
		return entity.UserInfo{}
	}
	user = get.(entity.UserInfo)
	return user
}
func SetCurrentUser(c *gin.Context, userInfo entity.UserInfo) {
	session := sessions.Default(c)
	session.Set("currentUser", userInfo)
	// 一定要Save否则不生效，若未使用gob注册User结构体，调用Save时会返回一个Error
	err := session.Save()
	if err != nil {
		panic(err)
	}
}

func postHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "postHandler",
	})
}

// 用于签名的字符串
var mySigningKey = []byte("golang_jwt_djt7djdh")

func GenRegisteredClaims() (string, error) {
	// 创建claims
	var claims = &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		Issuer:    "whj",
	}
	// 生成token
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateRegisterClaims(tokeString string) bool {
	token, err := jwt.Parse(tokeString, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("夏天夏天悄悄过去")

type CustomClaims struct {
	UserName             string
	jwt.RegisteredClaims //内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		username, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.UserName)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
