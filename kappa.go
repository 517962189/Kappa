package kappa

//河童框架
import (
	"errors"
	"github.com/517962189/Kappa/inits"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"
)

const (

)

var (
	GinAddrEmpty = errors.New("register gin addr empty")

	KappaInstance *Kappa
	once          sync.Once
)

type Kappa struct {
	inits.HooksStorageInterface
	inits.GormStorageInterface
	inits.ConfigStorageInterface
	engine *gin.Engine
}

func Default() *Kappa {
	if KappaInstance == nil {
		once.Do(func() {
			KappaInstance = &Kappa{
				engine: gin.Default(),
			}
		})
	}
	return KappaInstance
}

//获取所有config 列表
func (s *Kappa) RegisterConfigStorage(key string) *viper.Viper {
	if _, ok := inits.ConfigStorage[key]; !ok {
		panic(key + "::" + inits.FileNoExist)
	}
	return inits.ConfigStorage[key]
}

func (s *Kappa) RegisterDbStorage(group string) *gorm.DB {
	if _, ok := inits.GormStorage[group]; !ok {
		panic(inits.DbPoolNoFound)
	}
	return inits.GormStorage[group]
}

func (s *Kappa) RegisterHooksStorage(userFunc inits.UserFunc) {
	inits.UserFuncSlice = append(inits.UserFuncSlice, userFunc)
}

// 获取原生gin
func (s *Kappa) Gin() *gin.Engine {
	return s.engine
}

func (s *Kappa) loadMiddleWare() {
	globalMiddleWare := inits.NewMiddleWare().GetGlobal()
	//路由加载中间件
	if len(globalMiddleWare) > 0 {
		s.engine.Use(globalMiddleWare...)
	}
}

func (s *Kappa) loadHookStorage() {

	if len(inits.UserFuncSlice) == 0 {
		return
	}

	for _, function := range inits.UserFuncSlice {
		err := function()
		if err != nil {
			panic(err)
		}
	}
}

func (s *Kappa) Run(addr ...string) {
	//注册listern 监听服务
	if len(addr) == 0 {
		panic(GinAddrEmpty)
	}

	//初始化DB
	inits.LoadGormStorage()

	//初始化中间件
	s.loadMiddleWare()

	//加载钩子数据
	s.loadHookStorage()

	err := s.engine.Run(addr...)
	if err != nil {
		panic(err)
	}
}
