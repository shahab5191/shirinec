package routes

import (
	handler "shirinec.com/internal/handlers"
	"shirinec.com/internal/middlewares"
	"shirinec.com/internal/services"
)

func (r *router) setupMediaRouter() {
	mediaService := services.NewMediaService(r.Deps.MediaRepo, r.Deps.ItemRepo, r.Deps.CategoryRepo)
	mediaHandler := handler.NewMediaHandler(mediaService)

	flags := middlewares.AuthMiddleWareFlags{
		ShouldBeActive: true,
	}

	authMiddleware := middlewares.AuthMiddleWare(flags, r.db)

	r.GinEngine.POST("/media/upload", authMiddleware, mediaHandler.Upload)
}
