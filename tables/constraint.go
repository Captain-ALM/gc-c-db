package tables

type Constraint struct {
	Name            string
	Command         string
	CascadeOnDelete bool
}

func (cs Constraint) GetAddCommand() string {
	theCommand := "ADD CONSTRAINT " + cs.Name + " " + cs.Command
	if cs.CascadeOnDelete {
		theCommand += " ON DELETE CASCADE"
	}
	return theCommand
}

func (cs Constraint) GetDropCommand() string {
	return "DROP CONSTRAINT " + cs.Name
}
