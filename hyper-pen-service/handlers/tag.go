package handlers

import (
	"hyper-pen-service/models"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// TagHandler 处理标签相关的请求
type TagHandler struct {
	db *gorm.DB
}

// NewTagHandler 创建新的标签处理器
func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

// GetTags 获取用户的所有标签
func (h *TagHandler) GetTags(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	if userID == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "未授权访问"})
		return
	}

	var tags []models.Tag
	if err := h.db.Where("user_id = ?", userID).Find(&tags).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取标签失败"})
		return
	}

	ctx.JSON(tags)
}

// CreateTag 创建新标签
func (h *TagHandler) CreateTag(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	if userID == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "未授权访问"})
		return
	}

	var tag models.Tag
	if err := ctx.ReadJSON(&tag); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "无效的请求数据"})
		return
	}

	tag.ID = uuid.New().String()
	tag.UserID = userID

	if err := h.db.Create(&tag).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "创建标签失败"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(tag)
}

// UpdateTag 更新标签
func (h *TagHandler) UpdateTag(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	if userID == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "未授权访问"})
		return
	}

	tagID := ctx.Params().Get("id")
	var tag models.Tag
	if err := h.db.Where("id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "标签不存在"})
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "获取标签失败"})
		}
		return
	}

	if err := ctx.ReadJSON(&tag); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "无效的请求数据"})
		return
	}

	if err := h.db.Save(&tag).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "更新标签失败"})
		return
	}

	ctx.JSON(tag)
}

// DeleteTag 删除标签
func (h *TagHandler) DeleteTag(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	if userID == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(iris.Map{"error": "未授权访问"})
		return
	}

	tagID := ctx.Params().Get("id")
	var tag models.Tag
	if err := h.db.Where("id = ? AND user_id = ?", tagID, userID).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "标签不存在"})
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "获取标签失败"})
		}
		return
	}

	if err := h.db.Delete(&tag).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "删除标签失败"})
		return
	}

	ctx.JSON(iris.Map{"message": "标签已删除"})
}
