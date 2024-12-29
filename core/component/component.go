package component

type Component interface {
	Init() error
	Start() error
	Stop() error
	Register() error
}
