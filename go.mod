module github.com/60ke/trais

go 1.19

replace github.com/60ke/trais => ../trias

require (
	github.com/ethereum/go-ethereum v1.11.4
	github.com/go-ini/ini v1.67.0
	github.com/hnlq715/golang-lru v0.3.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/robfig/cron/v3 v3.0.1
	go.uber.org/zap v1.24.0
	gorm.io/driver/mysql v1.4.7
	gorm.io/gorm v1.24.6
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
