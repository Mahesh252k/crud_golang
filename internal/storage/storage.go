package storage

import "github.com/Mahesh252k/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, emai string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
}
