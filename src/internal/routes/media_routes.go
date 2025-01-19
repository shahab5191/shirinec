package routes

import (
	handler "shirinec.com/src/internal/handlers"
	"shirinec.com/src/internal/middlewares"
	"shirinec.com/src/internal/services"
)

func (r *router) setupMediaRouter() {
	mediaService := services.NewMediaService(r.Deps.MediaRepo, r.Deps.ItemRepo, r.Deps.CategoryRepo)
	mediaHandler := handler.NewMediaHandler(mediaService)

	flags := middlewares.AuthMiddleWareFlags{
		ShouldBeActive: true,
	}

	authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

	r.GinEngine.POST("/media/upload", authMiddleware, mediaHandler.Upload)
    r.GinEngine.GET("/file/:fileName", authMiddleware, mediaHandler.GetMedia)
    r.GinEngine.POST("/file/:fileName", authMiddleware, mediaHandler.UpdateMedia)
}
