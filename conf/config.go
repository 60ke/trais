package conf

import (
	"time"

	"github.com/60ke/trais/log"

	"github.com/go-ini/ini"
)

// APP
type APP struct {
	LogLevel string
	LogName  string
}

var (
	APPSetting      = &APP{}
	ServerSetting   = &Server{}
	DatabaseSetting = &Mysql{}
	TaskSetting     = &Task{}
	BscHosts        = &[]string{}
)

type Server struct {
	// 后端服务运行模式release or debug
	RunMode string
	// 后端服务端口
	HttpPort int

	// server超时设置
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 数据库相关配置
type Mysql struct {
	User     string
	Password string
	Host     string
	Bsc      string
	TM       string
}

// 任务定时配置
type Task struct {
	BscSync string
	TMSync  string
}

var cfg *ini.File

func ConfInit(path string) {
	var err error
	cfg, err = ini.Load(path)
	if err != nil {
		log.Logger.Fatalf("setting.Setup, fail to parse 'conf.ini': %v", err)
	}

	mapTo("app", APPSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("task", TaskSetting)

	rpcs := cfg.Section("cluster")

	*BscHosts = rpcs.Key("bschost").Strings(",")

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Logger.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
