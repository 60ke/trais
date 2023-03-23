package conf

import (
	"time"

	"github.com/60ke/trais/log"

	"github.com/go-ini/ini"
)

var (
	APPSetting        = &APPConf{}
	ServerSetting     = &ServerConf{}
	DatabaseSetting   = &MysqlConf{}
	TaskSetting       = &TaskConf{}
	BscHosts          = &[]string{}
	DownloaderSetting = &DownloaderConf{}
)

// APPConf
type APPConf struct {
	LogLevel string
	LogName  string
}

type DownloaderConf struct {
	BscStep int64
	TMStep  int64
}

type ServerConf struct {
	// 后端服务运行模式release or debug
	RunMode string
	// 后端服务端口
	HttpPort int

	// server超时设置
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 数据库相关配置
type MysqlConf struct {
	User     string
	Password string
	Host     string
	Bsc      string
	TM       string
}

// 任务定时配置
type TaskConf struct {
	BscSync          string
	TMSync           string
	UpdateBscBalance string
	HandleFailBlock  string
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
	mapTo("downloader", DownloaderSetting)

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
