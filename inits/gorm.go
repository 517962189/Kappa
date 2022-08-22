package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"time"
)

var DB map[string]*gorm.DB

func InitGorm() {

	DB = make(map[string]*gorm.DB, 0)
	//初始化DB内存
	if _, ok := Configs["database"]; !ok {
		return
	}
	//初始化database配置数据
	dbConfig := Configs["database"]
	//多数据库连接
	for group := range dbConfig.AllSettings() {
		dns := JoinDns(group, dbConfig)

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
		DB[group] = dbConnect
	}
	log.Println("DB connect Successful!")
}

/**
 * key string  多数据库连接 database 数据库库名称
 */
func JoinDns(group string, dbConfig *viper.Viper) string {
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
