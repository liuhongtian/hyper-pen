package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"hyper-pen-service/models"
	"time"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// ShareHandler 处理共享笔记相关的请求
type ShareHandler struct {
	db *gorm.DB
}

// NewShareHandler 创建新的共享处理器
func NewShareHandler(db *gorm.DB) *ShareHandler {
	return &ShareHandler{db: db}
}

// GetSharedNote 获取共享的笔记
func (h *ShareHandler) GetSharedNote(ctx iris.Context) {
	token := ctx.Params().Get("token")

	var shareLink models.ShareLink
	if err := h.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&shareLink).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "分享链接不存在或已过期",
		})
		return
	}

	var note models.Note
	if err := h.db.Preload("Category").Preload("Tags").First(&note, shareLink.NoteID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "笔记不存在",
		})
		return
	}

	ctx.JSON(note)
}

// CreateShareLinkRequest 创建分享链接的请求
type CreateShareLinkRequest struct {
	ExpiresIn int `json:"expires_in"` // 过期时间（小时），0表示永久
}

// CreateShareLink 创建分享链接
func (h *ShareHandler) CreateShareLink(ctx iris.Context) {
	noteID := ctx.Params().Get("id")
	userID := ctx.Values().Get("userID").(uint)

	// 解析请求
	var req CreateShareLinkRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error": "无效的请求参数",
		})
		return
	}

	// 检查笔记是否存在且属于当前用户
	var note models.Note
	if err := h.db.First(&note, noteID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "笔记不存在",
		})
		return
	}

	if note.UserID != userID {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"error": "无权操作此笔记",
		})
		return
	}

	// 计算过期时间
	var expiresAt time.Time
	if req.ExpiresIn > 0 {
		expiresAt = time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
	} else {
		// 设置为100年后过期，相当于永久
		expiresAt = time.Now().Add(100 * 365 * 24 * time.Hour)
	}

	// 生成分享链接
	shareLink := models.ShareLink{
		NoteID:    noteID,
		Token:     generateToken(),
		ExpiresAt: expiresAt,
	}

	if err := h.db.Create(&shareLink).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "创建分享链接失败",
		})
		return
	}

	ctx.JSON(shareLink)
}

// GetShareLinks 获取笔记的所有分享链接
func (h *ShareHandler) GetShareLinks(ctx iris.Context) {
	noteID := ctx.Params().GetUintDefault("id", 0)
	userID := ctx.Values().Get("userID").(uint)

	// 检查笔记是否存在且属于当前用户
	var note models.Note
	if err := h.db.First(&note, noteID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "笔记不存在",
		})
		return
	}

	if note.UserID != userID {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"error": "无权操作此笔记",
		})
		return
	}

	var shareLinks []models.ShareLink
	if err := h.db.Where("note_id = ?", noteID).Find(&shareLinks).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "获取分享链接失败",
		})
		return
	}

	ctx.JSON(shareLinks)
}

// DeleteShareLink 删除分享链接
func (h *ShareHandler) DeleteShareLink(ctx iris.Context) {
	shareLinkID := ctx.Params().GetUintDefault("id", 0)
	userID := ctx.Values().Get("userID").(uint)

	// 检查分享链接是否存在且属于当前用户
	var shareLink models.ShareLink
	if err := h.db.First(&shareLink, shareLinkID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "分享链接不存在",
		})
		return
	}

	var note models.Note
	if err := h.db.First(&note, shareLink.NoteID).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{
			"error": "笔记不存在",
		})
		return
	}

	if note.UserID != userID {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{
			"error": "无权操作此分享链接",
		})
		return
	}

	if err := h.db.Delete(&shareLink).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{
			"error": "删除分享链接失败",
		})
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

// generateToken 生成随机的分享链接token
func generateToken() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(token)
}
