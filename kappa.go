package kappa
//河童框架
import (
	"github.com/gin-gonic/gin"
	"github.com/517962189/Kappa/inits"
)

type KappaServer struct {
	engine *gin.Engine
}

func NewKappaServer() *KappaServer {
	return &KappaServer{
		gin.Default(),
	}
}

func (s *KappaServer) RegisterMiddleWare() {
	globalMiddleWare := inits.NewMiddleWare().GetGlobal()
	//路由加载中间件
	if len(globalMiddleWare) > 0 {
		s.engine.Use(globalMiddleWare...)
	}
}

func (s *KappaServer) Run(){
	//初始化钩子函数
	inits.RunHooks()
	//初始化数据库
	inits.InitGorm()
	//初始化中间件
	s.RegisterMiddleWare()

	s.RegisterServer()
}

//注册listern 监听服务
func (s *KappaServer) RegisterServer() error{
	return s.engine.Run(":6615")
}


