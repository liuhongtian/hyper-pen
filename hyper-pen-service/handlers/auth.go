package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"hyper-pen-service/config"
	"hyper-pen-service/models"
	"hyper-pen-service/utils"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login 处理用户登录
func (h *AuthHandler) Login(ctx iris.Context) {
	var req LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(LoginResponse{
		Token: token,
		User:  user,
	})
}

// Register 处理用户注册
func (h *AuthHandler) Register(ctx iris.Context) {
	var req RegisterRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := h.db.Create(&user).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create user"})
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(LoginResponse{
		Token: token,
		User:  user,
	})
}

// GitHubOAuthLogin 处理GitHub OAuth登录请求
func (h *AuthHandler) GitHubOAuthLogin(ctx iris.Context) {
	w := ctx.ResponseWriter()
	r := ctx.Request()
	// 构建GitHub OAuth授权URL
	authURL := "https://github.com/login/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", config.AppConfig.GitHubClientID)
	params.Add("redirect_uri", config.AppConfig.GitHubRedirectURI)
	params.Add("scope", "user:email")

	// 重定向到GitHub授权页面
	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
}

// GitHubOAuthCallback 处理GitHub OAuth回调
func (h *AuthHandler) GitHubOAuthCallback(ctx iris.Context) {
	w := ctx.ResponseWriter()
	//r := ctx.Request()
	code := ctx.URLParam("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	// 获取访问令牌
	accessToken, err := h.getGitHubAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// 获取用户信息
	userInfo, err := h.getGitHubUserInfo(accessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// 查找或创建用户
	var user models.User
	result := h.db.Where("github_id = ?", userInfo.ID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = models.User{
			Username:    userInfo.Login,
			Email:       userInfo.Email,
			GithubID:    userInfo.ID,
			AvatarURL:   userInfo.AvatarURL,
			GitHubToken: accessToken,
		}
		if err := h.db.Create(&user).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if result.Error != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	} else {
		// 更新现有用户信息
		user.AvatarURL = userInfo.AvatarURL
		user.GitHubToken = accessToken
		if err := h.db.Save(&user).Error; err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getGitHubAccessToken 获取GitHub访问令牌
func (h *AuthHandler) getGitHubAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", config.AppConfig.GitHubClientID)
	data.Set("client_secret", config.AppConfig.GitHubClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.AppConfig.GitHubRedirectURI)

	resp, err := http.Post("https://github.com/login/oauth/access_token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	return values.Get("access_token"), nil
}

// GitHubUserInfo GitHub用户信息结构
type GitHubUserInfo struct {
	ID        string `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// getGitHubUserInfo 获取GitHub用户信息
func (h *AuthHandler) getGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// WechatLogin 处理微信登录请求
func (h *AuthHandler) WechatLogin(ctx iris.Context) {
	w := ctx.ResponseWriter()
	r := ctx.Request()
	// 生成随机state参数，用于防止CSRF攻击
	state := generateRandomString(16)

	// 构建微信授权URL
	authURL := "https://open.weixin.qq.com/connect/qrconnect"
	params := url.Values{}
	params.Add("appid", config.AppConfig.WechatAppID)
	params.Add("redirect_uri", config.AppConfig.WechatRedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_login")
	params.Add("state", state)

	// 重定向到微信授权页面
	http.Redirect(w, r, authURL+"?"+params.Encode(), http.StatusFound)
}

// WechatCallback 处理微信回调
func (h *AuthHandler) WechatCallback(ctx iris.Context) {
	w := ctx.ResponseWriter()
	//r := ctx.Request()
	code := ctx.URLParam("code")
	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	// 获取访问令牌
	accessToken, openID, err := h.getWechatAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	// 获取用户信息
	userInfo, err := h.getWechatUserInfo(accessToken, openID)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// 查找或创建用户
	var user models.User
	result := h.db.Where("wechat_id = ?", openID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = models.User{
			WechatID:  openID,
			Username:  userInfo.Nickname,
			AvatarURL: userInfo.HeadImgURL,
		}
		if err := h.db.Create(&user).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	} else if result.Error != nil {
		http.Error(w, "Failed to find user", http.StatusInternalServerError)
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(&user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 返回成功响应
	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"avatar_url": user.AvatarURL,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// WechatUserInfo 微信用户信息结构
type WechatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// getWechatAccessToken 获取微信访问令牌
func (h *AuthHandler) getWechatAccessToken(code string) (string, string, error) {
	params := url.Values{}
	params.Add("appid", config.AppConfig.WechatAppID)
	params.Add("secret", config.AppConfig.WechatAppSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")

	resp, err := http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?" + params.Encode())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result.AccessToken, result.OpenID, nil
}

// getWechatUserInfo 获取微信用户信息
func (h *AuthHandler) getWechatUserInfo(accessToken string, openID string) (*WechatUserInfo, error) {
	params := url.Values{}
	params.Add("access_token", accessToken)
	params.Add("openid", openID)
	params.Add("lang", "zh_CN")

	resp, err := http.Get("https://api.weixin.qq.com/sns/userinfo?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo WechatUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
