package config

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

const (
	PGSQL string = "postgres"
)

// Database 数据库配置
type Database struct {
	DBName              string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	UserName            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableLgLog         bool   `mapstructure:"enable_lg_log" json:"enable_lg_log" yaml:"enable_lg_log"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

func InitDb(conf *Configuration) *gorm.DB {
	switch conf.Database.Driver {
	case "postgres":
		return initPGSQL(&conf.Database, conf)
	default:
		return initPGSQL(&conf.Database, conf)
	}
	return nil
}
func initPGSQL(database *Database, conf *Configuration) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		database.Host,
		database.UserName,
		database.Password,
		database.Database,
		database.Port,
	)

	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
		Logger: logger.New(getGormLogWriter(database, conf), logger.Config{
			SlowThreshold: 200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:      logger.Info,            // 日志级别
		}),
	}

	// gorm将类名转换成数据库表名的逻辑
	if gormConfig.NamingStrategy == nil {
		gormConfig.NamingStrategy = schema.NamingStrategy{
			SingularTable: true,
		}
	}
	if database.EnableLgLog {
		gormConfig.Logger = getGormLogger(database, conf)
	}
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(database.MaxIdleConns)
	return db
}
func getGormLogger(dbConfig *Database, conf *Configuration) logger.Interface {
	var logMode logger.LogLevel

	switch dbConfig.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(dbConfig, conf), logger.Config{
		SlowThreshold:             200 * time.Millisecond,        // 慢 SQL 阈值
		LogLevel:                  logMode,                       // 日志级别
		IgnoreRecordNotFoundError: false,                         // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  !dbConfig.EnableFileLogWriter, // 禁用彩色打印
	})
}

// 自定义 接管gorm日志，打印到文件 or 控制台
func getGormLogWriter(dbConfig *Database, conf *Configuration) logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if dbConfig.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   conf.Log.Dir + "/" + dbConfig.LogFilename,
			MaxSize:    conf.Log.MaxSize,
			MaxBackups: conf.Log.MaxBackups,
			MaxAge:     conf.Log.MaxAge,
			Compress:   conf.Log.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}
