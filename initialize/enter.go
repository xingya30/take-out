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

// 好了，那我现在算是把enter.go文件看完了，做的全是初始化操作，把四个全局变量分别赋值成结构体，最后的router被复制成一个gin框架的默认路由，
// 顺便还注册了admin/的各种post put get 请求时的对应函数，我说的对吗

// enter.go 初始化配置、日志、MySQL 和 Redis，保存到四个全局变量；然后创建并配置 Gin 引擎，
// 注册中间件和 /admin 下的各模块接口，最后返回引擎指针，供 main 启动服务。
