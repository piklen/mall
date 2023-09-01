package conf

import (
	"gopkg.in/ini.v1"
	"mall/dao"
	"strings"
)

var (
	AppMode  string
	HttpPort string

	DB         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string

	ValidEmail string
	SmtpHost   string
	SmtpEmail  string
	SmtpPass   string

	Host        string
	ProductPath string
	AvatarPath  string
)

func Init() {
	//本地读取环境变量
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		panic(err)
	}
	LoadingServer(file)
	LoadingMySql(file)
	LoadingRedis(file)
	LoadingEmail(file)
	LoadingPhotoPath(file)
	//mysql读
	pathRead := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	//MySQL写
	pathWrite := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	dao.Database(pathRead, pathWrite)

}
func LoadingServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	AppMode = file.Section("service").Key("HttpPort").String()
}
func LoadingMySql(file *ini.File) {
	AppMode = file.Section("mysql").Key("DB").String()
	AppMode = file.Section("mysql").Key("DbHost").String()
	AppMode = file.Section("mysql").Key("DbPort").String()
	AppMode = file.Section("mysql").Key("DbUser").String()
	AppMode = file.Section("mysql").Key("DbPassword").String()
	AppMode = file.Section("mysql").Key("DbName").String()
}
func LoadingRedis(file *ini.File) {
	AppMode = file.Section("redis").Key("RedisDb").String()
	AppMode = file.Section("redis").Key("RedisAddr").String()
	AppMode = file.Section("redis").Key("RedisPw").String()
	AppMode = file.Section("redis").Key("RedisDbName").String()
}
func LoadingEmail(file *ini.File) {
	AppMode = file.Section("email").Key("ValidEmail").String()
	AppMode = file.Section("email").Key("SmtpHost").String()
	AppMode = file.Section("email").Key("SmtpEmail").String()
	AppMode = file.Section("email").Key("SmtpPass").String()
}
func LoadingPhotoPath(file *ini.File) {
	AppMode = file.Section("path").Key("Host").String()
	AppMode = file.Section("path").Key("ProductPath").String()
	AppMode = file.Section("path").Key("AvatarPath").String()
}
