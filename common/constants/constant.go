/*
@Time : 2019/4/29 10:48 AM
@Author : yangping
@File : init
@Desc :
*/
package constants

const (
	// 命令行提示语句
	StartMessage = "server conf file , default is ./server.conf"

	StartRunMode = "RUN_MODE"

	PRO = "pro"

	ProConfigFile = "./config-pro.yml"

	DevConfigFile = "./config.yml"
)

const (
	ZERO = iota
	ONE
	TWO
	THREE
	// 空字符串
	EmptyStr = ""
	// 下划线
	UnderLine = "_"
	// 默认页数
	DefaultPage = 1
	// 默认每页数量
	DefaultPageSize = 10
	// Mysql 最大空闲连接数
	MysqlMaxIdleConn = 1000
	// Mysql 最大连接数
	MysqlMaxOpenConn = 2000
)

const (
	//
	Domain = "http://172.16.1.50:9069/v1/api/go?tinyUrl="
	// 短链生成方式   默认
	ConvertDefault = "default"
	// 短链生成方式	自定义
	ConvertCustom = "custom"
	// sign
	URL = "url"
	// tiny key
	TinyUrl = "tiny"
	// long key
	LongUrl = "long"
	// 过期时间 s
	ExpireTime = 60 * 1
)

// 数据库表  mongodb
const (
	TinyGroup = "tiny_group"
	TinyInfo  = "tiny_info"
	User      = "user"
	JwtToken  = "jwt_token"
)

// 默认配置参数
const (
	JwtSecret         = "scncysyp"
	JwtExpireTime     = 3600
	Issuer            = "YaPi"
	RedirectUrlHeader = "http://"
)
