package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GabrielMSosa/menu-frontend-go/internal/domain"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrRestrictFK    = errors.New("error Locality no exist")
	ErrGenericDriver = errors.New("error bad query")
	ErrNotFound      = errors.New("Menu not found")
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	GetById(ctx context.Context, id int) (domain.Menu, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Menu) (int, error)
	Update(ctx context.Context, s domain.Menu) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// Retorna todos las filas de la db
func (r *repository) GetAll(ctx context.Context) ([]domain.Menu, error) {
	query := "SELECT * FROM menu"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var menu []domain.Menu
	for rows.Next() {
		s := domain.Menu{}
		_ = rows.Scan(&s.ID, &s.Icon, &s.RouterLink, &s.Text)
		menu = append(menu, s)
	}
	return menu, nil
}

func (r *repository) GetById(ctx context.Context, id int) (domain.Menu, error) {
	query := "SELECT * FROM menu WHERE id=?;"
	row := r.db.QueryRow(query, id)
	s := domain.Menu{}
	err := row.Scan(&s.ID, &s.Icon, &s.RouterLink, &s.Text)
	if err != nil {
		return domain.Menu{}, err
	}

	return s, nil
}
func (r *repository) Exists(ctx context.Context, cid int) bool {
	query := "SELECT id FROM menu WHERE id=?;"
	row := r.db.QueryRow(query, cid)
	err := row.Scan(&cid)
	return err == nil
}
func (r *repository) Save(ctx context.Context, s domain.Menu) (int, error) {
	query := "INSERT INTO menu (icon,router_link,Text) VALUES(?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(s.Icon, s.RouterLink, s.Text)
	//capturamos otros errores que SQL a veces no nos muestra
	driverErr, ok := err.(*mysql.MySQLError)
	if ok {
		fmt.Println("Error method save exec query:", driverErr.Number, " || ", driverErr.Message)
		switch {
		case driverErr.Number == 1452:
			return 0, ErrRestrictFK
		default:
			return 0, ErrGenericDriver
		}
	}
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
func (r *repository) Update(ctx context.Context, s domain.Menu) error {
	query := "UPDATE menu SET icon=?,router_link=?,Text=? WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(s.Icon, s.RouterLink, s.Text, s.ID)
	//capturamos otros errores que SQL a veces no nos muestra
	driverErr, ok := err.(*mysql.MySQLError)
	if ok {
		fmt.Println("Error method save exec query:", driverErr.Number, " || ", driverErr.Message)
		switch {
		case driverErr.Number == 1452:
			return ErrRestrictFK
		default:
			return ErrGenericDriver
		}
	}
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM menu WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect < 1 {
		return ErrNotFound
	}
	return nil
}
