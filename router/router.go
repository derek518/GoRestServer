package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
	"GoRestServer/api/auth"
	"GoRestServer/middleware"
	"GoRestServer/pkg/config"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// Add a logger middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	r.Use(logger.SetLogger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Authorization", "Content-Type", "Access-Control-Allow-Origin", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	authMiddleware, err := middleware.JWT()
	if err != nil {
		log.Fatal().Err(err).Msg("Init JWT")
		panic(err)
	}
	apiRoute := r.Group("/api")

	apiRoute.POST("/login", authMiddleware.LoginHandler)
	apiRoute.GET("/refresh_token", authMiddleware.RefreshHandler)

	RegisterAuthRouter(apiRoute, authMiddleware.MiddlewareFunc())

	return r
}

func RegisterAuthRouter(r *gin.RouterGroup, auth gin.HandlerFunc) {
	authRouter := r.Group("/auth")
	if *config.App.NeedAuth {
		authRouter.Use(auth)
	}

	// Function
	{
		authRouter.POST("/save_function", api_auth.SaveFunction)
		authRouter.POST("/delete_function", api_auth.DeleteFunction)
		authRouter.GET("/get_function", api_auth.GetFunction)
		authRouter.GET("/get_function_all", api_auth.GetFunctionAll)
		authRouter.POST("/get_function_page", api_auth.GetFunctionPage)
	}

	// Role
	{
		authRouter.POST("/save_role", api_auth.SaveRole)
		authRouter.POST("/delete_role", api_auth.DeleteRole)
		authRouter.GET("/get_role", api_auth.GetRole)
		authRouter.GET("/get_role_all", api_auth.GetRoleAll)
		authRouter.POST("/get_role_page", api_auth.GetRolePage)
	}

	// User
	{
		authRouter.POST("/save_user", api_auth.SaveUser)
		authRouter.POST("/delete_user", api_auth.DeleteUser)
		authRouter.GET("/get_user", api_auth.GetUser)
		authRouter.GET("/get_user_all", api_auth.GetUserAll)
		authRouter.POST("/get_user_page", api_auth.GetUserPage)
	}
}
