package mysql

import (
	"sync"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModMysqlName,
}

type callerStore struct {
	Name   string
	caller sync.Map
	cfg    Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller(name string) *gorm.DB {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	return obj.(*gorm.DB)
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Mysql {
		db, err := provider(cfg)
		if err != nil {
			return err
		}
		defaultCaller.caller.Store(name, db)
	}
	return nil
}

func provider(cfg CallerCfg) (resp *gorm.DB, err error) {
	var db *gorm.DB
	// dsn = "username:password@tcp(addr)/stt_config?charset=utf8&parseTime=True&loc=Local&readTimeout=1s&timeout=1s&writeTimeout=1s"
	db, err = gorm.Open(cfg.Dialect, cfg.Username+":"+cfg.Password+"@"+cfg.Network+"("+cfg.Addr+")/"+cfg.Db+
		"?charset="+cfg.Charset+"&parseTime="+cfg.ParseTime+"&loc="+cfg.Loc+
		"&timeout="+cfg.Timeout.Duration.String()+"&readTimeout="+cfg.ReadTimeout.Duration.String()+"&writeTimeout="+cfg.WriteTimeout.Duration.String())
	if err != nil {
		return
	}
	db.LogMode(cfg.Debug)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetConnMaxLifetime(cfg.ConnMaxLifetime.Duration)
	err = db.DB().Ping()

	if err != nil {
		return
	}
	resp = db
	return
}
