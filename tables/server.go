package tables

import (
	"time"
	"xorm.io/xorm"
)

type Server struct {
	ID            uint32    `xorm:"int unique pk not null autoincr 'Server_ID'"`
	Address       string    `xorm:"varchar(1020) unique not null 'Server_Address'"`
	LastCheckTime time.Time `xorm:"timestamp not null"`
}

func (Server) TableName() string {
	return "Servers"
}

func (Server) GetConstraints() []Constraint {
	return nil
}

func (Server) GetChildTables() []Table {
	return []Table{&Game{}}
}

func (s Server) GetID() interface{} {
	return s.ID
}

func (s Server) GetIDObject() Table {
	return &Server{ID: s.ID}
}

func (Server) GetNullableColumns() []string {
	return nil
}

func (s Server) GetChildrenGames(e *xorm.Engine) ([]Game, error) {
	var children []Game
	err := e.Where("Server_ID = ?", s.ID).Find(&children)
	if err == nil {
		return children, nil
	}
	return nil, err
}
