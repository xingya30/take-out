package initialize

import (
	"take-out/config"
	"take-out/global"

	"github.com/Meng-Xin/logger"
	"github.com/gin-gonic/gin"
)

func GlobalInit() *gin.Engine {
	// 配置文件初始化
	// 也就是说全局变量config等于configData *AllConfig这个结构体咯
	global.Config = config.InitLoadConfig()
	// Log初始化
	// 也就是说global.log是一个logtool,下面这个代码是作者写了个库
	// 然后套了三层函数，把配置交给库，库就创建出一个可以记录日志的对象，保存在 logTool 中。
	global.Log = logger.NewZapLogCenter(logger.NewZapConfig(
		logger.WithServiceName("takeout"),
	))
	// Gorm初始化
	// 返回了一个db结构体
	global.DB = InitDatabase(global.Config.DataSource.Dsn())
	// Redis初始化
	global.Redis = initRedis()
	//Router初始化
	router := routerInit()
	return router
}
