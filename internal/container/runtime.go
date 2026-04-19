package container

type Runtime interface {
	Start(name string) error
	Stop(name string) error
	List() ([]string, error)
	Remove(name string) error
}
