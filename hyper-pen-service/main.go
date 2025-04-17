package main

import (
	"hyper-pen-service/handlers"
	"hyper-pen-service/middleware"
	"hyper-pen-service/models"

	"github.com/kataras/iris/v12"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := iris.New()

	// 连接数据库
	db, err := gorm.Open(sqlite.Open("hyper-pen.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 自动迁移数据库表
	db.AutoMigrate(&models.User{}, &models.Note{}, &models.Category{}, &models.Tag{}, &models.ShareLink{})

	// 创建处理器
	authHandler := handlers.NewAuthHandler(db)
	noteHandler := handlers.NewNoteHandler(db)
	shareHandler := handlers.NewShareHandler(db)
	tagHandler := handlers.NewTagHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)

	// 注册路由
	api := app.Party("/api")
	{
		// 认证相关路由
		auth := api.Party("/auth")
		{
			auth.Post("/register", authHandler.Register)
			auth.Post("/login", authHandler.Login)
			auth.Get("/github", authHandler.GitHubOAuthLogin)
			auth.Get("/github/callback", authHandler.GitHubOAuthCallback)
			auth.Get("/wechat", authHandler.WechatLogin)
			auth.Get("/wechat/callback", authHandler.WechatCallback)
		}

		// 笔记相关路由
		notes := api.Party("/notes")
		notes.Use(middleware.AuthRequired)
		{
			notes.Get("", noteHandler.GetNotes)
			notes.Post("", noteHandler.CreateNote)
			notes.Get("/{id:string}", noteHandler.GetNote)
			notes.Put("/{id:string}", noteHandler.UpdateNote)
			notes.Delete("/{id:string}", noteHandler.DeleteNote)
			notes.Get("/search", noteHandler.SearchNotes)

			// 分享相关路由
			notes.Get("/{id:string}/share-links", shareHandler.GetShareLinks)
			notes.Post("/{id:string}/share-links", shareHandler.CreateShareLink)
		}

		// 分享链接相关路由
		shareLinks := api.Party("/share-links")
		shareLinks.Use(middleware.AuthRequired)
		{
			shareLinks.Delete("/{id:string}", shareHandler.DeleteShareLink)
		}

		// 标签相关路由
		tags := api.Party("/tags")
		tags.Use(middleware.AuthRequired)
		{
			tags.Get("", tagHandler.GetTags)
			tags.Post("", tagHandler.CreateTag)
			tags.Put("/{id:string}", tagHandler.UpdateTag)
			tags.Delete("/{id:string}", tagHandler.DeleteTag)
		}

		// 分类相关路由
		categories := api.Party("/categories")
		categories.Use(middleware.AuthRequired)
		{
			categories.Get("", categoryHandler.GetCategories)
			categories.Post("", categoryHandler.CreateCategory)
			categories.Put("/{id:string}", categoryHandler.UpdateCategory)
			categories.Delete("/{id:string}", categoryHandler.DeleteCategory)
		}

		// 共享笔记路由（不需要认证）
		api.Get("/shared/{token:string}", shareHandler.GetSharedNote)
	}

	app.Listen(":8080")
}
