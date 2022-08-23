package inits

import "log"

//注册钩子 例如: pprof  , prometheus , swagger 等
type userFunc func() error

var userFuncSlice []userFunc

//注册
func RegisterHooks(userFunc userFunc) {
	userFuncSlice = append(userFuncSlice, userFunc)
}

//运行
func init() {
	if len(userFuncSlice) == 0 {
		return
	}

	for _, function := range userFuncSlice {
		err := function()
		if err != nil{
			log.Fatalln("hooks run fail",err)
		}
	}
}
