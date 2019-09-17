package config

import (
	"flag"
	"github.com/jinzhu/configor"
	"path/filepath"
	"strings"
	"time"
)

// Configuration is stuff that can be configured externally per env variables or config file (config.yml,config.json).
type Configuration struct {
	App struct {
		RunMode   string `default:"debug"`
		JwtSecret string `default:"17K9347KKRB41ppP8buWYY5g4Jb2Z34DUc"`
		NeedAuth  *bool  `default:"true"`

		ImageSavePath string `default:"upload/images/"`
		//MB
		ImageMaxSize   int      `default:"5"`
		ImageAllowExts []string `default:"[.jpg,.jpeg,.png]"`

		QrCodeSavePath string `default:"qrcode/"`

		LogFilePath string `default:"logs/"`
		LogFileName string `default:"server"`
		LogFileExt  string `default:"log"`
		TimeFormat  string `default:"20190909"`
	}
	Server struct {
		ListenAddr string `default:""`
		Port       int    `default:"8080"`
		SSL        struct {
			Enabled         *bool  `default:"false"`
			RedirectToHTTPS *bool  `default:"true"`
			ListenAddr      string `default:""`
			Port            int    `default:"443"`
			CertFile        string `default:""`
			CertKey         string `default:""`
			LetsEncrypt     struct {
				Enabled   *bool  `default:"false"`
				AcceptTOS *bool  `default:"false"`
				Cache     string `default:"data/certs"`
				Hosts     []string
			}
		}
		ResponseHeaders map[string]string
		Stream          struct {
			AllowedOrigins []string
		}
	}
	Database struct {
		Type        string `default:"mysql"`
		User        string `default:"root"`
		Password    string `default:"123456"`
		Url         string `default:"@tcp(127.0.0.1:3306)/trainsystem?charset=utf8mb4&parseTime=True&loc=Local"`
		TablePrefix string `default:"tss_"`
		ShowSql     bool   `default:"false"`
	}
	Redis struct {
		Enabled    *bool  `default:"false"`
		Host       string `default:"127.0.0.1:6379"`
		Password   string
		Expiration time.Duration `default:"3600"`
	}
}

var Config = &Configuration{}
var App = &Config.App
var Server = &Config.Server
var Database = &Config.Database
var Redis = &Config.Redis

func addTrailingSlashToPaths(conf *Configuration) {
	if !strings.HasSuffix(conf.App.ImageSavePath, "/") && !strings.HasSuffix(conf.App.ImageSavePath, "\\") {
		conf.App.ImageSavePath = conf.App.ImageSavePath + string(filepath.Separator)
	}
	if !strings.HasSuffix(conf.App.QrCodeSavePath, "/") && !strings.HasSuffix(conf.App.QrCodeSavePath, "\\") {
		conf.App.QrCodeSavePath = conf.App.QrCodeSavePath + string(filepath.Separator)
	}
	if !strings.HasSuffix(conf.App.LogFilePath, "/") && !strings.HasSuffix(conf.App.LogFilePath, "\\") {
		conf.App.LogFilePath = conf.App.LogFilePath + string(filepath.Separator)
	}

	conf.Database.TablePrefix = strings.TrimSpace(conf.Database.TablePrefix)
}

// Initialize the configuration instance
func Setup() {
	//  获取命令行参数
	var (
		configPath = flag.String("conf", "", "configuration file path")
	)
	flag.Parse()

	files := []string{"config.yml", "conf/config.yml"}

	// Overwrite the config file path from command line.
	if len(*configPath) > 0 {
		files = []string{*configPath}
	}
	err := configor.New(&configor.Config{ENVPrefix: "ENV"}).Load(Config, files...)
	if err != nil {
		panic(err)
	}

	addTrailingSlashToPaths(Config)
}
