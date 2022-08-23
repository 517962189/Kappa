package inits

import (
	"fmt"
	"github.com/517962189/Kappa"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"time"
)

const (
	DbPoolNoFound = "db connect pool not found"
)

var GormStorage map[string]*gorm.DB

func loadGormStorage() {
	GormStorage = make(map[string]*gorm.DB, 0)

	//初始化database配置数据
	dbConfig := kappa.KappaInstance.RegisterConfigStorage("database")
	//多数据库连接
	for group := range dbConfig.AllSettings() {
		dns := joinDns(group, dbConfig)

		fmt.Println(group, "---", dns)
		dbConnect, err := gorm.Open(mysql.Open(dns), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: dbConfig.GetBool(strings.Join([]string{group, "DbDriver.SingularTable"}, ".")), // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
		})
		if err != nil {
			panic(err)
		}

		sqlDb, _ := dbConnect.DB()
		//设置连接池
		sqlDb.SetMaxIdleConns(dbConfig.GetInt(strings.Join([]string{group, "DbDriver.SetMaxIdleConns"}, ".")))
		//打开
		sqlDb.SetMaxOpenConns(dbConfig.GetInt(strings.Join([]string{group, "DbDriver.SetMaxOpenConns"}, ".")))
		//超时
		sqlDb.SetConnMaxLifetime(time.Second * time.Duration(dbConfig.GetInt(strings.Join([]string{group, "DbDriver.SetConnMaxLifetime"}, "."))))
		GormStorage[group] = dbConnect
	}
	log.Println("DB connect Successful!")
}

/**
 * key string  多数据库连接 database 数据库库名称
 */
func joinDns(group string, dbConfig *viper.Viper) string {
	dbConfig.GetString(strings.Join([]string{group, "Password"}, "."))
	user := dbConfig.GetString(strings.Join([]string{group, "User"}, "."))
	pwd := dbConfig.GetString(strings.Join([]string{group, "Password"}, "."))
	host := dbConfig.GetString(strings.Join([]string{group, "Host"}, "."))
	port := dbConfig.GetString(strings.Join([]string{group, "Port"}, "."))
	dbName := dbConfig.GetString(strings.Join([]string{group, "Dbname"}, "."))
	charset := dbConfig.GetString(strings.Join([]string{group, "Charset"}, "."))
	dsn := strings.Join([]string{user, ":", pwd, "@tcp(", host, ":", port, ")/", dbName, "?charset=", charset, "&parseTime=True&loc=Local"}, "")
	return dsn
}

//提供接口暴露接口
type GormStorageInterface interface {
	RegisterDbStorage(group string) *gorm.DB
}
