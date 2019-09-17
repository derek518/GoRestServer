package middleware

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoas/go-funk"
	"net/http"
	"time"
	"GoRestServer/model"
	"GoRestServer/model/auth"
	"GoRestServer/pkg/config"
	"GoRestServer/service/auth"
)

var (
	IdentityKey = "username"
	RoleKey     = "role_id"
)

func JWT() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "train system",
		Key:         []byte(config.App.JwtSecret),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model_auth.User); ok {
				return jwt.MapClaims{
					IdentityKey: v.Username,
					RoleKey:     v.RoleId,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model_auth.User{
				Username: claims[IdentityKey].(string),
				RoleId:   claims[RoleKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginObj model_auth.Login
			if err := c.ShouldBindJSON(&loginObj); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			userService := service_auth.UserServiceInstance(model_auth.UserModelInstance(model.SQL))
			return userService.Login(&loginObj)
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Debug().Interface("token", data).Str("URL", c.Request.URL.String()).Msg("Authorize")
			//if token, ok := c.Get("JWT_TOKEN"); ok {
			//	var expires time.Time
			//	if err := cache.Get("token_"+token.(string), &expires); err == nil {
			//		if expires.Sub(time.Now()) <= 0 {
			//			return false
			//		}
			//	}
			//} else {
			//	return false
			//}

			if user, ok := data.(*model_auth.User); ok {
				roleService := service_auth.RoleServiceInstance(model_auth.RoleModelInstance(model.SQL))
				role := roleService.GetByID(user.RoleId)
				funk.Find(role.Functions, func(function *model_auth.Function) bool {
					return function.Url == c.Request.URL.String()
				})
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Info().Str("URL", c.Request.URL.String()).Int("code", code).Msgf("Authorize failed: %s", message)
			c.JSON(code, gin.H{
				"code":    1,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			//cache.AddWithTimeout("token_"+token, expire, expire.Sub(time.Now()))
			c.JSON(http.StatusOK, gin.H{
				"code":   0,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   0,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
