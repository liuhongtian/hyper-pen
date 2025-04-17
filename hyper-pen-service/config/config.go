package config

import (
	"os"
)

type Config struct {
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURI  string
	JWTSecret          string
	WechatAppID        string
	WechatAppSecret    string
	WechatRedirectURI  string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURI:  getEnv("GITHUB_REDIRECT_URI", "http://localhost:3000/auth/github/callback"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		WechatAppID:        getEnv("WECHAT_APP_ID", ""),
		WechatAppSecret:    getEnv("WECHAT_APP_SECRET", ""),
		WechatRedirectURI:  getEnv("WECHAT_REDIRECT_URI", "http://localhost:3000/auth/wechat/callback"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
