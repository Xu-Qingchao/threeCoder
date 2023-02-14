package config

import (
	"github.com/spf13/viper"
	"os"
)

func InitConfig() {
	workDir, _ := os.Getwd()                 // 找到工作目录
	viper.SetConfigName("config")            // 配置文件的文件名
	viper.SetConfigType("yml")               // 配置文件的后缀
	viper.AddConfigPath(workDir + "/config") // 获取到配置文件的路径
	err := viper.ReadInConfig()              // 读取
	if err != nil {
		return
	}
}
