package kappa
//河童框架
import (
	"github.com/517962189/Kappa/inits"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"
)

const (
	GinAddrEmpty = "register gin addr empty"
)
var (
	KappaInstance *KappaServer
	Once sync.Once
)

type KappaServer struct {
	inits.GormStorageInterface
	inits.ConfigStorageInterface
	engine *gin.Engine
}

func NewKappaServer() *KappaServer {
	if KappaInstance == nil {
		Once.Do(func() {
			KappaInstance = &KappaServer{
				engine:  gin.Default(),
			}
		})
	}
	return KappaInstance
}

//获取所有config 列表
func (s *KappaServer)RegisterConfigStorage(key string) *viper.Viper {
	if _, ok := inits.ConfigStorage[key]; !ok {
		panic(key+"::"+ inits.FileNoExist)
	}
	return inits.ConfigStorage[key]
}

func (s *KappaServer)RegisterDbStorage(group string) *gorm.DB {
	if _, ok := inits.GormStorage[group]; !ok {
		panic(inits.DbPoolNoFound)
	}
	return inits.GormStorage[group]
}



// 获取原生gin
func (s *KappaServer) Gin() *gin.Engine {
	return s.engine
}

func (s *KappaServer) registerMiddleWare() {
	globalMiddleWare := inits.NewMiddleWare().GetGlobal()
	//路由加载中间件
	if len(globalMiddleWare) > 0 {
		s.engine.Use(globalMiddleWare...)
	}
}

func (s *KappaServer) Run(addr ...string){

	//初始化DB
	inits.loadGormStorage()

	//初始化中间件
	s.registerMiddleWare()
	//注册listern 监听服务

	if len(addr) == 0{
		panic(GinAddrEmpty)
	}
	err := s.engine.Run(addr...)
	if err != nil{
		panic(err)
	}
}


