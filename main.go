package main

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "GoRestServer/docs"
	"GoRestServer/model"
	"GoRestServer/pkg/cache"
	"GoRestServer/pkg/config"
	"GoRestServer/router"
	"GoRestServer/runner"
)

var server struct {
	Config *config.Configuration
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.App.RunMode == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func init() {
	config.Setup()
	initLog()
	model.Setup()
	cache.Setup()
}

// @title Train System RestAPI Server
// @version 1.0
// @description This is a RestAPI server for train system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email zhengsl@nicetech-video.com

// @license.name TOPZEN
// @license.url

// @host localhost:8080
// @BasePath /api

// @tag.name QueryCondition
// @tag.description 支持多级表查询，查询条件参考SQL Where中的操作符，包括等于({key: value})、大于({key:{"gt": value})、小于({key:{"lt": value})、大于等于({key:{"gte": value})、小于等于({key:{"lte": value})、between({key:{"between": [value1, value2]})、like({key:{"like": regex})和in({key:{"in": [...]})。
// @tag.docs.url https://www.w3schools.com/sql/sql_where.asp

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @tokenUrl https://localhost/api/login
func main() {
	gin.SetMode(config.App.RunMode)
	log.Info().Fields(structs.Map(config.Config)).Msg("dump configuration")

	model.AutoMigrate()

	//demo.JsonTest()
	//demo.TestSql()

	engine := router.InitRouter()

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//log.Fatal().Err(engine.Run(addr))
	runner.Run(engine)

}
