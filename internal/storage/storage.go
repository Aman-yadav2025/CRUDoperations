package storage

import "CRUDapp/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student,error)
	GetAll()([]types.Student,error)
	UpdateStudentById(name string, email string, age int, id int64) (types.Student,error)
	DeleteStudentById(id int64) (bool,error)
}
