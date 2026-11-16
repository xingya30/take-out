package initialize

import (
	"errors"
	"take-out/global"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	GormToManyRequestError = errors.New("gorm: to many request")
)

func InitDatabase(dsn string) *gorm.DB {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	// mysql.config{}创建一个mysql结构体，然后使用mysql.New()函数把mysql结构体转换成dialector对象
	// 为什么要转换成dialector对象？
	// 因为想要用gorm.Open()函数把mysql结构体转换成gorm结构体，需要先用new函数转换成dialector对象
	// gorm.open()函数的接口的第一个参数要求dialecter结构体，不认识未转换过的mysql结构体
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	// 慢日志中间件
	SlowQueryLog(db)

	return db
}

// SlowQueryLog 慢查询日志
func SlowQueryLog(db *gorm.DB) {
	err := db.Callback().Query().Before("*").Register("slow_query_start", func(d *gorm.DB) {
		now := time.Now()
		d.Set("start_time", now)
	})
	if err != nil {
		panic(err)
	}

	err = db.Callback().Query().After("*").Register("slow_query_end", func(d *gorm.DB) {
		now := time.Now()
		start, ok := d.Get("start_time")
		if ok {
			duration := now.Sub(start.(time.Time))
			// 一般认为 200 Ms 为Sql慢查询
			if duration > time.Millisecond*200 {
				global.Log.ErrContext(d.Statement.Context, "慢查询", "SQL:", d.Statement.SQL.String())
			}
		}
	})
	if err != nil {
		panic(err)
	}
}
