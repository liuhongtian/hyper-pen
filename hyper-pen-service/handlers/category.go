package handlers

import (
	"hyper-pen-service/models"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// CategoryHandler 处理分类相关的请求
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler 创建新的分类处理器
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetCategories 获取所有分类
func (h *CategoryHandler) GetCategories(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	var categories []models.Category
	if err := h.db.Where("user_id = ?", userID).Preload("Notes").Find(&categories).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "获取分类失败"})
		return
	}
	ctx.JSON(categories)
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	var category models.Category
	if err := ctx.ReadJSON(&category); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "无效的请求数据"})
		return
	}

	category.ID = uuid.New().String()
	category.UserID = userID

	if err := h.db.Create(&category).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "创建分类失败"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(category)
}

// UpdateCategory 更新分类
func (h *CategoryHandler) UpdateCategory(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	id := ctx.Params().Get("id")
	var category models.Category
	if err := ctx.ReadJSON(&category); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "无效的请求数据"})
		return
	}

	var existingCategory models.Category
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&existingCategory).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "分类不存在"})
		return
	}

	if err := h.db.Model(&existingCategory).Updates(category).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "更新分类失败"})
		return
	}

	ctx.JSON(existingCategory)
}

// DeleteCategory 删除分类
func (h *CategoryHandler) DeleteCategory(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)
	id := ctx.Params().Get("id")

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "分类不存在"})
		return
	}

	if err := h.db.Delete(&category).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "删除分类失败"})
		return
	}

	ctx.JSON(iris.Map{"message": "分类已删除"})
}
