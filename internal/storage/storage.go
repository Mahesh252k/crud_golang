package storage

type Storage interface {
	CreateStudent(name string, emai string, age int) (int64, error)
}
