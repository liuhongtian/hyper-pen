package handlers

import (
	"hyper-pen-service/models"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

// NoteHandler 处理笔记相关的请求
type NoteHandler struct {
	db *gorm.DB
}

// NewNoteHandler 创建新的笔记处理器
func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{db: db}
}

type NoteRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CategoryID string   `json:"category_id"`
	TagIDs     []string `json:"tag_ids"`
}

// CreateNote 创建笔记
func (h *NoteHandler) CreateNote(ctx iris.Context) {
	var req NoteRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	ctx.Application().Logger().Infof("创建笔记请求: %+v", req)

	userID := ctx.Values().Get("userID").(uint)

	note := models.Note{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
	}

	// 开始事务
	tx := h.db.Begin()

	// 创建笔记
	if err := tx.Create(&note).Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create note"})
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to associate tags"})
			return
		}
		if err := tx.Model(&note).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to associate tags"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to save note"})
		return
	}

	// 重新加载笔记以获取完整数据
	if err := h.db.Preload("Tags").Preload("Category").First(&note, "id=?", note.ID).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to load note"})
		return
	}

	ctx.JSON(iris.Map{
		"message": "Note created successfully",
		"note":    note,
	})
}

// GetNotes 获取笔记列表
func (h *NoteHandler) GetNotes(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)

	var notes []models.Note
	if err := h.db.Where("user_id = ?", userID).Order("created_at desc").Preload("Tags").Preload("Category").Find(&notes).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch notes"})
		return
	}

	ctx.JSON(notes)
}

// GetNote 获取单个笔记
func (h *NoteHandler) GetNote(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid note ID"})
		return
	}

	userID := ctx.Values().Get("userID").(uint)

	var note models.Note
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Preload("Tags").Preload("Category").First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "Note not found"})
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch note"})
		return
	}

	ctx.JSON(note)
}

// UpdateNote 更新笔记
func (h *NoteHandler) UpdateNote(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid note ID"})
		return
	}

	var req NoteRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	userID := ctx.Values().Get("userID").(uint)

	// 开始事务
	tx := h.db.Begin()

	var note models.Note
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "Note not found"})
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch note"})
		return
	}

	note.Title = req.Title
	note.Content = req.Content
	note.CategoryID = req.CategoryID

	if err := tx.Save(&note).Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update note"})
		return
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to update tags"})
			return
		}
		if err := tx.Model(&note).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to update tags"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update note"})
		return
	}

	// 重新加载笔记以获取完整数据
	if err := h.db.Where("id = ?", note.ID).Preload("Tags").Preload("Category").First(&note).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to load note"})
		return
	}

	ctx.JSON(iris.Map{
		"message": "Note updated successfully",
		"note":    note,
	})
}

// DeleteNote 删除笔记
func (h *NoteHandler) DeleteNote(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if id == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid note ID"})
		return
	}

	userID := ctx.Values().Get("userID").(uint)

	// 开始事务
	tx := h.db.Begin()

	// 先检查笔记是否存在
	var note models.Note
	if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "Note not found"})
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch note"})
		return
	}

	// 删除笔记
	if err := tx.Delete(&note).Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to delete note"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to delete note"})
		return
	}

	ctx.JSON(iris.Map{
		"message": "Note deleted successfully",
	})
}

// SearchNotes 搜索笔记
func (h *NoteHandler) SearchNotes(ctx iris.Context) {
	userID := ctx.Values().Get("userID").(uint)

	// 获取搜索参数
	query := ctx.URLParam("q")
	categoryID := ctx.URLParam("category_id")
	tagIDs := ctx.URLParamSlice("tag_ids")

	// 构建查询
	db := h.db.Where("user_id = ?", userID)

	// 关键词搜索
	if query != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// 分类筛选
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// 标签筛选
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			db = db.Joins("JOIN note_tags ON note_tags.note_id = notes.id").
				Where("note_tags.tag_id = ?", tagID)
		}
	}

	// 执行查询
	var notes []models.Note
	if err := db.Preload("Tags").Preload("Category").Find(&notes).Error; err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to search notes"})
		return
	}

	ctx.JSON(notes)
}
