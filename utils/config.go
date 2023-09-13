package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"zinx/ziface"
)

type Config struct {
	// Server
	Server ziface.IServer `mapstructure:"server"` // 当前Zinx全局的Server对象
	Host   string         `mapstructure:"host"`   // 当前服务器监听IP
	Port   int            `mapstructure:"port"`   // 监听端口
	Name   string         `mapstructure:"name"`   // 服务器名称

	// Zinx
	Version          string `mapstructure:"version"`             // Zinx版本号
	MaxConn          int    `mapstructure:"max_conn"`            // 服务器允许最大连接数
	MaxPackageSize   uint32 `mapstructure:"max_package_size"`    // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 `mapstructure:"worker_pool_size"`    // Zinx工作池数量
	MaxWorkerTaskLen uint32 `mapstructure:"max_worker_task_len"` //  允许开辟最大worker数量
}

var GlobalConfig *Config

func init() {
	GlobalConfig = &Config{
		Host:             "0.0.0.0",
		Port:             8999,
		Name:             "Zinx-Server-App",
		Version:          "v-default",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// Set up Viper
	viper.SetConfigFile("demo/zinx_v0.10/config/config.yaml") // Specify the configuration file name and location
	viper.SetConfigType("yaml")                               // Set the configuration file type (YAML in this case)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("ReadInConfig: %s \n", err)
	}

	// Create a Config struct and populate it with values from the configuration file
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		fmt.Printf("Unable to unmarshal config: %s \n", err)
	}

}
