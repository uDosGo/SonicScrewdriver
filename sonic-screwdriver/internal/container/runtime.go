package container

type Runtime interface {
	Start(name string) error
	Stop(name string) error
}
