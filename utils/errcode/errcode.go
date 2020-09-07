// @Author: Perry
// @Date  : 2020/4/21
// @Desc  : 

package errcode

import "errors"

const (
	OK                 = 200
	RUNERROR           = 1001
	CFGLOADERROR       = 4001
	INCOMPLETECFGERROR = 4002
	DBINITERROR        = 5001
	DBEMPTYQUERYSQL    = 5002
	DBTBLINITERROR     = 5003
	AUTHNOTALLOWED     = 6001
	UNKNOWNERROR       = 9999
)

var recodeText = map[int]string{
	RUNERROR:           "进程启动失败",
	CFGLOADERROR:       "配置文件读取失败",
	INCOMPLETECFGERROR: "配置文件项缺失",
	DBINITERROR:        "数据库初始化失败",
	DBEMPTYQUERYSQL:    "SQL语句内容为空",
	DBTBLINITERROR:     "数据库表初始化失败",
	AUTHNOTALLOWED:     "用户权限无法进行该操作",
	UNKNOWNERROR:       "未知错误",
}

func Text(code int) string {
	str, ok := recodeText[code]
	if ok {
		return str
	}
	return recodeText[UNKNOWNERROR]
}

func Error(code int) error {
	str, ok := recodeText[code]
	if ok {
		return errors.New(str)
	}
	return errors.New(recodeText[UNKNOWNERROR])
}
