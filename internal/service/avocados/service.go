package avocados

import (
	"entrega/internal/config"
	"github.com/jmoiron/sqlx"
)

//Avocado ...
type Avocado struct {
	ID        int64
	Name      string
	Image     string
	Price     int64
	Stock     int64
	Quantity  int64
}

//Service ...
type Service interface {
	AddAvocado(Avocado) string
	GetAll() []*Avocado
	GetByID(int64) (*Avocado, string)
	ChangeAvoQuantity(int64, int64) (string, error)
	ChangeAvoPrice(int64, int64) (string, error)
	DeleteByID(int64) (string, error)
}

type service struct {
	db   *sqlx.DB
	conf *config.Config
}

//New ...
func New(db *sqlx.DB, c *config.Config) (Service, error) {
	return service{db, c}, nil
}

func (s service) AddAvocado(a Avocado) string {
	insertAvocado := `INSERT INTO avocados (name, image, price, stock, quantity) VALUES(?,?,?,?,?)`
	data := s.db.MustExec(insertAvocado, a.Name, a.Image, a.Price, a.Stock, a.Quantity)
	
	if data != nil {
		return "Se agrego correctamente " + a.Name
	}

	return "Error al guardar"
}

func (s service) GetAll() []*Avocado {
	var list []*Avocado

	if err := s.db.Select(&list, "SELECT * FROM avocados"); err != nil {
		panic(err)
	}

	return list
}

func (s service) GetByID(ID int64) (*Avocado, string) {
	var avocado Avocado
	err := s.db.QueryRowx("SELECT * FROM avocados WHERE id = ?", ID).StructScan(&avocado)

	if err != nil {
		return nil, "No se avoencontro nada"
	}

	return &avocado, "Se trajo correctamente"
}

func (s service) ChangeAvoQuantity(id int64, newQuant int64)(string, error){
	var avo Avocado
	err := s.db.QueryRowx("SELECT * FROM avocados WHERE id = $1", id).StructScan(&avo)

	if avo.Name==""{
		return "No se encontro nada", err
	}

	_, err = s.db.Exec("UPDATE avocados SET quantity = $1", newQuant)
	if err != nil{
		return "Error al cambiar la cantidad", err
	}

	return "La cantidad del avocado se cambio correctamente", nil
}

func (s service) ChangeAvoPrice(id int64, newPrice int64)(string, error){
	var avo Avocado
	err := s.db.QueryRowx("SELECT * FROM avocados WHERE id = $1", id).StructScan(&avo)

	if avo.Name==""{
		return "No se encontro nada", err
	}

	_, err = s.db.Exec("UPDATE avocados SET price = $1", newPrice)
	if err != nil{
		return "Error al cambiar el precio", err
	}

	return "El precio del avocado se cambio correctamente", nil
}

func (s service) DeleteByID(ID int64) (string,error) {
	_, err := s.db.Exec("DELETE FROM avocados WHERE id = $1", ID)

	if err != nil {
		return "Error al borrar", err
	}

	return "Se borro correctamente", err
}