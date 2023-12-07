package tables

import "xorm.io/xorm"

type User struct {
	Email string `xorm:"char(255) unique pk not null 'User_Email'"`
	Token string `xorm:"char(255) unique index 'User_Token'"`
}

func (User) GetConstraints() []Constraint {
	return nil
}

func (User) GetChildTables() []Table {
	return []Table{&Quiz{}}
}

func (User) TableName() string {
	return "User"
}

func (u User) GetChildrenQuizzes(e *xorm.Engine) ([]Quiz, error) {
	var children []Quiz
	err := e.Where("Owner_Email = ?", u.Email).Find(&children)
	if err == nil {
		return children, nil
	}
	return nil, err
}
