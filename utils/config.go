package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"zinx/ziface"
)

type Config struct {
	// Server
	Server ziface.IServer // 当前Zinx全局的Server对象
	Host   string         // 当前服务器监听IP
	Port   int            // 监听端口
	Name   string         // 服务器名称

	// Zinx
	Version        string // Zinx版本号
	MaxConn        int    // 服务器允许最大连接数
	MaxPackageSize uint32 // 当前Zinx框架数据包的最大值
}

var GlobalConfig *Config

func init() {
	GlobalConfig = &Config{
		Host:           "0.0.0.0",
		Port:           8999,
		Name:           "Zinx-Server-App",
		Version:        "",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	// Set up Viper
	viper.SetConfigFile("demo/zinx_v0.7/config/config.yaml") // Specify the configuration file name and location
	viper.SetConfigType("yaml")                              // Set the configuration file type (YAML in this case)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("ReadInConfig: %s \n", err)
	}

	// Create a Config struct and populate it with values from the configuration file
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		fmt.Printf("Unable to unmarshal config: %s \n", err)
	}

}
