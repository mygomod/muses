package muses

import (
	"github.com/goecology/muses/pkg/common"
	"io/ioutil"
	"reflect"
)

// 通过反射取包里面的值
var orderCallerList = []callerAttr{
	{common.ModAppName},
	{common.ModEtcdName},
	{common.ModLoggerName},
	{common.ModPromName},
	{common.ModMysqlName},
	{common.ModOssName},
	{common.ModRedisName},
	{common.ModMongoName},
	{common.ModMixCacheName},
	{common.ModGinSessionName},
	{common.ModEchoSessionName},
	{common.ModTplBeegoName},
	{common.ModStatName},
	{common.ModGinName},
	{common.ModCmdName},
}

type callerAttr struct {
	Name string
}

// Container from file.
func parseFile(path string) ([]byte, error) {
	// read file to []byte
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return b, err
	}
	return b, nil
}

func sortCallers(callers []common.CallerFunc) (callerSort []common.Caller, err error) {
	callerMap := make(map[string]common.Caller)
	callerSort = make([]common.Caller, 0)
	for _, caller := range callers {
		obj := caller()
		name := getCallerName(obj)
		callerMap[name] = obj
	}

	for _, value := range orderCallerList {
		caller, ok := callerMap[value.Name]
		if ok {
			// 如果存在于map，加入到排序里的caller sort
			callerSort = append(callerSort, caller)
		}
	}
	return
}

func getCallerName(caller common.Caller) string {
	return reflect.ValueOf(caller).Elem().FieldByName("Name").String()
}
