package inits

//注册钩子 例如: pprof  , prometheus , swagger 等
type UserFunc func() error

var UserFuncSlice []UserFunc

type HooksStorageInterface interface {
	RegisterHooksStorage(userFunc UserFunc)
}
