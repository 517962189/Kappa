package inits

import "github.com/gin-gonic/gin"


type MiddleWare struct {
	global        []gin.HandlerFunc
	handleFuncMap map[string][]gin.HandlerFunc
}

func NewMiddleWare() *MiddleWare {
	return &MiddleWare{
		global:        make([]gin.HandlerFunc, 0),
		handleFuncMap: make(map[string][]gin.HandlerFunc, 0),
	}
}

//设置全局中间件
func (m *MiddleWare) SetGlobal(handleFuncSlice []gin.HandlerFunc) {
	m.global = append(m.global, handleFuncSlice...)
}

//设置全局中间件
func (m *MiddleWare) GetGlobal() []gin.HandlerFunc {
	return m.global
}

//增加组中间件
func (m *MiddleWare) SetGroup(groupName string, handleFuncSlice []gin.HandlerFunc) {
	if groupName == "" {
		return
	}

	m.handleFuncMap[groupName] = append(m.handleFuncMap[groupName], handleFuncSlice...)
}

//增加组中间件
func (m *MiddleWare) GetGroup(groupName string) map[string][]gin.HandlerFunc {
	if groupName == "" {
		return nil
	}
	return m.handleFuncMap
}
