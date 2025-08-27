// config/config.go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	ServerPort string
	JWTSecret  string
	LogLevel   string
	LogFile    string
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "123456"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		LogFile:    getEnv("LOG_FILE", "logs/app.log"),
	}
}

// getEnv 获取环境变量，如果不存在则使用默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取整型环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return value
}
