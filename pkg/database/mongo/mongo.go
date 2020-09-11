package mongo

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mygomod/muses/pkg/common"
	"os"
	"sync"
)

var defaultCaller = &callerStore{
	Name: common.ModMongoName,
}

type callerStore struct {
	Name   string
	caller sync.Map
	cfg    Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller(name string) *mgo.Session {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	return obj.(*mgo.Session)
}

func CopyDb(name string) *mgo.Database {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	session := obj.(*mgo.Session)
	for cfgName, cfg := range defaultCaller.cfg.Muses.Mongo {
		if cfgName == name {
			return session.Copy().DB(cfg.Database)
		}
	}
	return nil
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Mongo {
		db, err := provider(cfg)
		if err != nil {
			return err
		}
		defaultCaller.caller.Store(name, db)
	}
	return nil
}

func provider(cfg CallerCfg) (resp *mgo.Session, err error) {
	session, err := mgo.Dial(cfg.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	mgo.SetLogger(&cLogger{})
	mgo.SetDebug(cfg.Debug)
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return session, err
}

type cLogger struct{}

func (c *cLogger) Output(calldepth int, s string) error {
	fmt.Println("calldepth: ", calldepth, ", s: ", s)
	return nil
}
