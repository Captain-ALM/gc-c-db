package tables

import (
	"errors"
	"strconv"
	"xorm.io/xorm"
)

type Guest struct {
	ID     uint32 `xorm:"integer(10) unique pk not null autoincr 'Guest_ID'"`
	Name   string `xorm:"char(64) not null 'Guest_Name'"`
	GameID uint32 `xorm:"integer(10) not null 'Game_ID'"`
	Score  uint32 `xorm:"integer(10) not null 'Score'"`
}

func (Guest) TableName() string {
	return "Guest"
}

func (Guest) GetConstraints() []Constraint {
	return []Constraint{{
		Name:            "FK_Game_ID",
		Command:         "FOREIGN KEY (Game_ID) REFERENCES Game(Game_ID)",
		CascadeOnDelete: true,
	}}
}

func (Guest) GetChildTables() []Table {
	return nil
}

func (g Guest) GetParentGame(e *xorm.Engine) (Game, error) {
	parent := Game{ID: g.GameID}
	has, err := e.Get(&parent)
	if err != nil {
		return parent, err
	}
	if !has {
		return parent, errors.New("parent game not found for ID " + strconv.Itoa(int(g.ID)))
	}
	return parent, nil
}
