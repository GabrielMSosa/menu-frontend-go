package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GabrielMSosa/menu-frontend-go/internal/domain"
	"github.com/GabrielMSosa/menu-frontend-go/internal/repository"
)

var (
	ErrNotFound         = errors.New("Menu not found")
	ErrIdcEqual         = errors.New("Conflict Cid equal")
	ErrNoSaveOrInternal = errors.New("Internal Server Error")
	ErrGeneric          = errors.New("Internal Server Error")
	ErrIdLocality       = errors.New("Conflict locality_od not exist")
	ErrRestrictFK       = errors.New("Error Locality not exist")
)

type Service interface {
	// GetAll retorna un slice de Menu si esta todo ok
	//
	// Retorna un error generico si falla algo en la db.
	GetAll(ctx context.Context) (ret []domain.Menu, err error)
	//Save retorna el valor guardado en la db si fue exitoso.
	//
	// Retorna ErrIdcEqual en caso de queres escribir un Menu con CID existente.
	//
	// Retorna ErrGeneric si no puede guardar en la db o si busca el valor con el id.
	Save(ctx context.Context, s domain.Menu) (ret domain.Menu, err error)
	//GetByID retorna el Menu si todo esta ok
	//
	//Retorna ErrNotFound en caso de no encontra con el valor
	GetByID(ctx context.Context, id int) (ret domain.Menu, err error)
	//Delete borra item y rotorna nil si esta todo bien.
	//
	//Retorna ErrGeneric si tiene un error al llamar al repository.
	//
	// Retorna ErrNotFound si no encuentra con la id
	Delete(ctx context.Context, id int) (err error)

	//Update retorna el valor actualizado si todo esta ok.
	//
	//Retorna ErrNotFound si no encuentra con la id
	//
	//Retorna ErrGeneric en caso de que no pueda actualizar el Menu.
	//
	// Retorna ErrConflict si intenta escrivir el campo CID con un valor existente.
	Update(ctx context.Context, data map[string]interface{}, id int) (ret domain.Menu, err error)
}

type service struct {
	rep repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{
		rep: repo,
	}
}

func (sv *service) GetAll(ctx context.Context) (ret []domain.Menu, err error) {
	ret, err = sv.rep.GetAll(ctx)
	if err != nil {
		fmt.Println("Service GetAll", err.Error())
		err = ErrGeneric
		return
	}
	return
}
func (sv *service) Save(ctx context.Context, s domain.Menu) (ret domain.Menu, err error) {
	b := sv.rep.Exists(ctx, s.ID)
	if b {
		err = ErrIdcEqual
		fmt.Println("Service Save CID ", err.Error())
		return
	}
	s.ID = 0
	id, err := sv.rep.Save(ctx, s)
	if err != nil {
		fmt.Println("Service Save", err.Error())
		switch {
		case errors.Is(err, repository.ErrRestrictFK):
			err = ErrIdLocality
			return
		default:
			err = ErrNoSaveOrInternal
			return
		}
	}
	//vamos a ahacer un get para ver si se guardo realmente el valor
	ret, err = sv.rep.GetById(ctx, id)
	if err != nil {
		fmt.Println("Service save data in get by id for return: ", err.Error())
		err = ErrNoSaveOrInternal
	}

	return
}
func (sv *service) GetByID(ctx context.Context, id int) (ret domain.Menu, err error) {
	ret, err = sv.rep.GetById(ctx, id)
	if err != nil {
		fmt.Println("ERROR IN GET BY ID SERVICE GetByID", err.Error())
		err = ErrNotFound
		return
	}
	return
}
func (sv *service) Delete(ctx context.Context, id int) (err error) {
	_, err = sv.rep.GetById(ctx, id)
	if err != nil {
		fmt.Println("ERROR IN GET BY ID SERVICE delete", err.Error())
		err = ErrNotFound
		return
	}
	err = sv.rep.Delete(ctx, id)
	if err != nil {
		fmt.Println("ERROR IN DELETE ID SERVICE delete", err.Error())
		err = ErrGeneric
		return
	}
	return
}
func (sv *service) Update(ctx context.Context, data map[string]interface{}, id int) (ret domain.Menu, err error) {
	//TOY RE PAJA PARA ESTE METODO
	return
}
