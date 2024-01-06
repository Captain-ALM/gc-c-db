package tables

import (
	"errors"
	"strconv"
	"time"
	"xorm.io/xorm"
)

type Game struct {
	ID            uint32    `xorm:"integer(10) unique pk not null autoincr 'Game_ID'"`
	QuizID        uint32    `xorm:"integer(10) not null 'Quiz_ID'"`
	ServerID      uint32    `xorm:"integer(10) not null 'Server_ID'"`
	State         byte      `xorm:"tinyint(3) not null 'Game_State'"`
	Expiry        time.Time `xorm:"timestamp not null 'Game_Expiry'"`
	CountdownMax  uint32    `xorm:"integer(10) not null 'Game_Countdown_Max'"`
	StreakEnabled bool      `xorm:"bit(1) not null 'Game_Streak'"`
	QuestionNo    uint32    `xorm:"integer(10) not null 'Game_Question'"`
}

func (Game) TableName() string {
	return "Game"
}

func (Game) GetConstraints() []Constraint {
	return []Constraint{{
		Name:            "FK_Quiz_ID",
		Command:         "FOREIGN KEY (Quiz_ID) REFERENCES Quiz(Quiz_ID)",
		CascadeOnDelete: true,
	}, {
		Name:            "FK_Server_ID",
		Command:         "FOREIGN KEY (Server_ID) REFERENCES Servers(Server_ID)",
		CascadeOnDelete: true,
	}}
}

func (Game) GetChildTables() []Table {
	return []Table{&Guest{}}
}

func (g Game) GetID() interface{} {
	return g.ID
}

func (g Game) GetIDObject() Table {
	return &Game{ID: g.ID}
}

func (Game) GetNullableColumns() []string {
	return nil
}

func (g Game) GetParentQuiz(e *xorm.Engine) (Quiz, error) {
	parent := Quiz{ID: g.QuizID}
	has, err := e.Get(&parent)
	if err != nil {
		return parent, err
	}
	if !has {
		return parent, errors.New("parent quiz not found for ID " + strconv.Itoa(int(g.ID)))
	}
	return parent, nil
}

func (g Game) GetParentServer(e *xorm.Engine) (Server, error) {
	parent := Server{ID: g.ServerID}
	has, err := e.Get(&parent)
	if err != nil {
		return parent, err
	}
	if !has {
		return parent, errors.New("parent server not found for ID " + strconv.Itoa(int(g.ID)))
	}
	return parent, nil
}

func (g Game) GetChildrenGuests(e *xorm.Engine) ([]Guest, error) {
	var children []Guest
	err := e.Where("Game_ID = ?", g.ID).Find(&children)
	if err == nil {
		return children, nil
	}
	return nil, err
}
