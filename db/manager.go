package db

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"golang.local/gc-c-db/tables"
	"strings"
	"xorm.io/xorm"
)

type Manager struct {
	Engine *xorm.Engine
	Path   string
}

const ManagerEngineNil = "manager engine is nil"

func (m *Manager) Connect() error {
	var err error
	if strings.HasPrefix(m.Path, "mysql:") {
		m.Engine, err = xorm.NewEngine("mysql", m.Path[strings.Index(m.Path, ":")+1:])
	} else if strings.HasPrefix(m.Path, "sqlite:") {
		m.Engine, err = xorm.NewEngine("sqlite", m.Path[strings.Index(m.Path, ":")+1:])
	} else {
		err = errors.New("only sqlite and mysql are supported")
	}
	if err != nil {
		DebugPrintln("Connect Failure For: " + m.Path)
		m.Engine = nil
		return err
	}
	return nil
}

func (m *Manager) AssureAllTables() error {
	err := m.AssureTable(&tables.Server{}, false)
	if err != nil {
		return err
	}
	err = m.AssureTable(&tables.User{}, true)
	if err != nil {
		DebugPrintln("Tables Assured")
	}
	return err
}

func (m *Manager) AssureTable(t tables.Table, assureChildren bool) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	exists, err := m.Engine.IsTableExist(t)
	if err == nil {
		err = m.Engine.Sync(t)
		if err == nil && !exists {
			err = m.addConstraints(t)
			if assureChildren && err == nil {
				for _, ctbl := range t.GetChildTables() {
					err := m.AssureTable(ctbl, true)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if err != nil {
		DebugPrintln("Table Assured: " + t.TableName())
	}
	return err
}

func (m *Manager) DropTable(t tables.Table, dropChildren bool) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	exists, err := m.Engine.IsTableExist(t)
	if err == nil {
		if exists {
			if dropChildren {
				for _, ctbl := range t.GetChildTables() {
					err := m.DropTable(ctbl, true)
					if err != nil {
						return err
					}
				}
			}
			err = m.dropConstraints(t)
			if err == nil {
				err = m.Engine.DropTables(t)
			}
			/*} else {
			err = errors.New("table " + t.TableName() + " does not exist")*/
		} else {
			DebugPrintln("Drop ~ Table Non Existent: " + t.TableName())
		}
		if err != nil {
			DebugPrintln("Table Dropped: " + t.TableName())
		}
	}
	return err
}

func (m *Manager) DropAllTables() error {
	err := m.DropTable(&tables.User{}, true)
	if err != nil {
		return err
	}
	err = m.DropTable(&tables.Server{}, false)
	if err != nil {
		DebugPrintln("Tables Dropped")
	}
	return err
}

func (m *Manager) ClearTable(t tables.Table, clearChildren bool) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	exists, err := m.Engine.IsTableExist(t)
	if err == nil {
		if exists {
			if clearChildren {
				for _, ctbl := range t.GetChildTables() {
					err := m.ClearTable(ctbl, true)
					if err != nil {
						return err
					}
				}
			}
			err = m.dropConstraints(t)
			if err == nil {
				_, err = m.Engine.Exec("DELETE FROM " + t.TableName())
				if err == nil {
					err = m.addConstraints(t)
				}
			}
		} else {
			err = errors.New("table " + t.TableName() + " does not exist")
			DebugPrintln("Clear ~ Table Non Existent: " + t.TableName())
		}
	}
	if err != nil {
		DebugPrintln("Table Cleared: " + t.TableName())
	}
	return err
}

func (m *Manager) ClearAllTables() error {
	err := m.ClearTable(&tables.User{}, true)
	if err != nil {
		return err
	}
	err = m.ClearTable(&tables.Server{}, false)
	if err != nil {
		DebugPrintln("Tables Cleared")
	}
	return err
}

func (m *Manager) addConstraints(t tables.Table) error {
	for _, cnstr := range t.GetConstraints() {
		_, err := m.Engine.Exec("ALTER TABLE " + t.TableName() + " " + cnstr.GetAddCommand())
		if err != nil {
			return err
		}
	}
	DebugPrintln("Table Constrained: " + t.TableName())
	return nil
}

func (m *Manager) dropConstraints(t tables.Table) error {
	for _, cnstr := range t.GetConstraints() {
		_, err := m.Engine.Exec("ALTER TABLE " + t.TableName() + " " + cnstr.GetDropCommand())
		if err != nil {
			return err
		}
	}
	DebugPrintln("Table Unconstrained: " + t.TableName())
	return nil
}

const TableRecordNonExistent = "table record does not exist"

func (m *Manager) Load(t tables.Table) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	dbSession := m.Engine.NewSession()
	defer func() { _ = dbSession.Close() }()
	if len(t.GetNullableColumns()) > 0 {
		dbSession = dbSession.Nullable(t.GetNullableColumns()...)
	}
	exists, err := dbSession.Get(t)
	if err != nil {
		return err
	} else if !exists {
		DebugPrintln("Load Entry Non Existent: " + t.TableName())
		return errors.New(TableRecordNonExistent)
	}
	DebugPrintln("Load Entry Exists: " + t.TableName())
	return nil
}

func (m *Manager) Save(t tables.Table) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	idObj := t.GetIDObject()
	exists, err := m.Engine.Exist(idObj)
	if err != nil {
		return err
	}
	dbSession := m.Engine.AllCols()
	if len(t.GetNullableColumns()) > 0 {
		dbSession = dbSession.Nullable(t.GetNullableColumns()...)
	}
	if exists {
		DebugPrintln("Save Entry Exists: " + t.TableName())
		_, err := dbSession.ID(t.GetID()).Update(t)
		return err
	} else {
		DebugPrintln("Save Entry Non Existent: " + t.TableName())
	}
	_, err = dbSession.Insert(t)
	return err
}

func (m *Manager) Insert(t tables.Table) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	dbSession := m.Engine.AllCols()
	if len(t.GetNullableColumns()) > 0 {
		dbSession = dbSession.Nullable(t.GetNullableColumns()...)
	}
	DebugPrintln("Insert Entry: " + t.TableName())
	_, err := dbSession.Insert(t)
	return err
}

func (m *Manager) Delete(t tables.Table) error {
	if m.Engine == nil {
		return errors.New(ManagerEngineNil)
	}
	DebugPrintln("Delete Entry: " + t.TableName())
	_, err := m.Engine.Delete(t)
	return err
}
