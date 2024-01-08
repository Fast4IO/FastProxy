package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	// 生成全局配置对象
	GlobalConfig Config
)

type Config struct {
	// Define your configuration fields here
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	ProxyCnf ProxyConfig    `yaml:"proxy"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	Host         string `yaml:"host"`
	SupportProxy bool   `yaml:"support_proxy"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ProxyConfig struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

func init() {
	configFile := flag.String("c", "", "配置文件路径")
	// 解析命令行参数
	flag.Parse()

	fmt.Println("配置文件路径：", *configFile)
	// Read the conf.yml file
	file, err := os.Open(*configFile)
	if err != nil {
		log.Fatalf("Failed to read conf.yml: %v", err)
	}
	defer file.Close()

	// 创建解析器
	decoder := yaml.NewDecoder(file)

	// 解析 YAML 数据
	err = decoder.Decode(&GlobalConfig)
	if err != nil {
		fmt.Println("Error decoding YAML:", err)
		return
	}
}
