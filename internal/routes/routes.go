package routes

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/handlers"
)

type Router interface {
	SetupRouter()
	setupCategoryRouter()
	setupAuthRouter()
    setupUserRouter()
}

type router struct {
	GinEngine *gin.Engine
	Deps      *handler.Dependencies
}

func NewRouter(ginEngine *gin.Engine, deps *handler.Dependencies) Router {
	return &router{
		GinEngine: ginEngine,
		Deps:      deps,
	}
}

func (r *router) SetupRouter() {
	r.setupAuthRouter()
	r.setupCategoryRouter()
    r.setupUserRouter()
}
