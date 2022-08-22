package inits

import "os"

//系统默认环境常量值
const OsEnv = "KAPPA_ENV"

func GetOsEnv() string {
	if os.Getenv(OsEnv) == "" {
		os.Setenv(OsEnv, "dev")
	}
	return  os.Getenv(OsEnv)
}

func SetOsEnv(value string) error{
	 return os.Setenv(OsEnv,value)
}
