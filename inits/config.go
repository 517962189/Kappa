package inits

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	FileNoExist = "files no found"
)

var ConfigStorage map[string]*viper.Viper

func init() {
	//初始化内存空间
	ConfigStorage = make(map[string]*viper.Viper, 0)

	dir, _ := os.Getwd()
	dirRoot := dir + "/config/" + GetOsEnv()
	//验证目录是否存在 ，否则创建
	if err := dirExists(dirRoot); err != nil {
		os.MkdirAll(dirRoot, os.ModePerm)
	}
	walkDir(dirRoot)
}

//获取目录下所有文件
func walkDir(dir string) {
	filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			key := strings.Split(fi.Name(), ".")
			ConfigStorage[key[0]] = viperInstance(filename)
		}
		return nil
	})
}

//实例 viper 初始话所有配置文件
func viperInstance(path string) *viper.Viper {
	service := viper.New()
	service.SetConfigFile(path)
	if err := service.ReadInConfig(); err != nil {
		panic(err)
	}
	return service
}

//验证目录是否存在
func dirExists(dir string) error {
	_, err := os.Stat(dir)
	return err
}

//提供配置接口暴露接口
type ConfigStorageInterface interface {
	RegisterConfigStorage(key string) *viper.Viper
}
