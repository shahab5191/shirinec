package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"shirinec.com/internal/handlers"
)

type Router interface {
	SetupRouter()
	setupCategoryRouter()
	setupAuthRouter()
	setupUserRouter()
	setupItemRouter()
	setupAccountRouter()
	setupMediaRouter()
	setupFinancialGroupRouter()
    setupTransactionRouter()
    setupIncomeRouter()
}

type router struct {
	GinEngine *gin.Engine
	Deps      *handler.Dependencies
	db        *pgxpool.Pool
}

func NewRouter(ginEngine *gin.Engine, deps *handler.Dependencies, db *pgxpool.Pool) Router {
	return &router{
		GinEngine: ginEngine,
		Deps:      deps,
		db:        db,
	}
}

func (r *router) SetupRouter() {
	r.setupAuthRouter()
	r.setupCategoryRouter()
	r.setupUserRouter()
	r.setupItemRouter()
	r.setupAccountRouter()
	r.setupMediaRouter()
	r.setupFinancialGroupRouter()
    r.setupTransactionRouter()
    r.setupIncomeRouter()
}
