package tables

import (
	"encoding/json"
	"errors"
	"strconv"
	"xorm.io/xorm"
)

type Quiz struct {
	ID         uint32          `xorm:"integer(10) unique pk not null autoincr 'Quiz_ID'"`
	OwnerEmail string          `xorm:"char(255) not null 'Owner_Email'"`
	Name       string          `xorm:"varchar(1020) not null 'Quiz_Name'"`
	IsPublic   bool            `xorm:"bit(1) not null 'Quiz_Public'"`
	Questions  json.RawMessage `xorm:"blob not null 'Quiz_Questions'"`
	Answers    json.RawMessage `xorm:"blob not null 'Quiz_Answers'"`
}

func (Quiz) TableName() string {
	return "Quiz"
}

func (Quiz) GetConstraints() []Constraint {
	return []Constraint{{
		Name:            "FK_Owner_Email",
		Command:         "FOREIGN KEY (Owner_Email) REFERENCES User(User_Email)",
		CascadeOnDelete: true,
	}}
}

func (Quiz) GetChildTables() []Table {
	return []Table{&Game{}}
}

func (q Quiz) GetID() interface{} {
	return q.ID
}

func (q Quiz) GetIDObject() Table {
	return &Quiz{ID: q.ID}
}

func (Quiz) GetNullableColumns() []string {
	return nil
}

func (q Quiz) GetParentUser(e *xorm.Engine) (User, error) {
	parent := User{Email: q.OwnerEmail}
	has, err := e.Get(&parent)
	if err != nil {
		return parent, err
	}
	if !has {
		return parent, errors.New("parent user not found for ID " + strconv.Itoa(int(q.ID)))
	}
	return parent, nil
}

func (q Quiz) GetChildrenGames(e *xorm.Engine) ([]Game, error) {
	var children []Game
	err := e.Where("Quiz_ID = ?", q.ID).Find(&children)
	if err == nil {
		return children, nil
	}
	return nil, err
}
