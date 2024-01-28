package routes

import (
	"database/sql"

	"github.com/GabrielMSosa/menu-frontend-go/cmd/server/handler/menu"
	"github.com/GabrielMSosa/menu-frontend-go/cmd/server/middleware"
	"github.com/GabrielMSosa/menu-frontend-go/internal/repository"
	"github.com/GabrielMSosa/menu-frontend-go/internal/service"

	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{
		eng: eng,
		db:  db,
	}
}
func (r *router) MapRoutes() {
	r.setGroup()
	r.buildMenuRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildMenuRoutes() {
	repo := repository.NewRepository(r.db)
	service := service.NewService(repo)
	handler := menu.NewMenu(service)
	gr := r.rg.Group("/menu", middleware.LoggerMIddleware())
	gr.POST("", handler.Create())
	gr.GET("/:id", handler.GetById())
	gr.GET("", handler.GetAll())
}
