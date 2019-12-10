package sbbattery

type Routine struct {
	charge int
}

func New() *Routine {
	return new(Routine)
}

func (r *Routine) Update() error {
	return nil
}

func (r *Routine) String() string {
	return "battery"
}
