package routes

import (
	"github.com/gin-gonic/gin"
	"shirinec.com/internal/handlers"
)

type Router interface {
	SetupRouter()
	setupIncomeRouter()
	setupAuthRouter()
	setupExpenseRouter()
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
	r.setupIncomeRouter()
	r.setupExpenseRouter()
}
