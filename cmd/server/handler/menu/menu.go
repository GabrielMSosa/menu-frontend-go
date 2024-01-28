package menu

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/GabrielMSosa/menu-frontend-go/internal/domain"
	"github.com/GabrielMSosa/menu-frontend-go/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	ErrDataEmpty         = errors.New("error data empty")
	ErrInvalidNumber     = errors.New("Invalid telephone numbers")
	ErrInvalidData       = errors.New("Invalid address or Company name")
	ErrCidInvalid        = errors.New("invalid Cid")
	ErrInvalidName       = errors.New("error invalid name data")
	ErrInvalidType       = errors.New("error invalid type data")
	ErrINvalidLocalityId = errors.New("invalid LocalityId")
	ErrLocalityId        = errors.New("Locality not exist")
	ErrInvalidId         = errors.New("Invalid ID")
)

type Menu struct {
	srv service.Service
}

func NewMenu(s service.Service) *Menu {
	return &Menu{srv: s}
}

func (h *Menu) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		list, err := h.srv.GetAll(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, list)

	}
}
func (h *Menu) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := ValidateAndGetId(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		list, err := h.srv.GetByID(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
		ctx.JSON(http.StatusOK, list)

	}
}

func (h *Menu) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data domain.Menu
		err := ctx.ShouldBindJSON(&data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = ValidateCid(data.ID)
		switch {
		case errors.Is(err, ErrDataEmpty):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		case errors.Is(err, ErrCidInvalid):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = ValidaDatas(data.Icon)
		switch {
		case errors.Is(err, ErrDataEmpty):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		case errors.Is(err, ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		err = ValidaDatas(data.RouterLink)
		switch {
		case errors.Is(err, ErrDataEmpty):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		case errors.Is(err, ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return

		}
		err = ValidaDatas(data.Text)
		switch {
		case errors.Is(err, ErrDataEmpty):
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		case errors.Is(err, ErrInvalidData):
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		out, err := h.srv.Save(ctx, data)
		switch {
		case errors.Is(err, service.ErrIdcEqual):
			ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		case errors.Is(err, service.ErrGeneric):
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		case errors.Is(err, service.ErrIdLocality):
			ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, out)

	}

}

func ValidateAndGetId(c *gin.Context) (id int, err error) {
	idstring := c.Param("id")
	id, err = strconv.Atoi(idstring)
	if err != nil {
		return
	}
	if id < 0 {
		err = ErrInvalidId
		return
	}
	return
}

func ValidateCid(cid int) (err error) {
	if cid == 0 {
		err = ErrDataEmpty
		return
	}
	if cid < 0 {
		err = ErrCidInvalid
		return
	}
	return
}

func ValidaDatas(data string) (err error) {
	if data == "" {
		err = ErrDataEmpty
		return
	}
	if len(data) < 4 {
		err = ErrInvalidData
		return
	}
	return
}
