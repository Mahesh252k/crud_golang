package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Mahesh252k/students-api/internal/config"
	"github.com/Mahesh252k/students-api/internal/types"
)

type Mysql struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Mysql, error) {
	db, err := sql.Open("mysql", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Verify connection immediately
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// MySQL table creation
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        age INT
    );`)

	if err != nil {
		return nil, err
	}

	return &Mysql{Db: db}, nil
}

// CreateStudent implements the storage.Storage interface
func (s *Mysql) CreateStudent(name string, email string, age int) (int64, error) {
	query := "INSERT INTO students (name, email, age) VALUES (?, ?, ?)"

	result, err := s.Db.Exec(query, name, email, age)
	if err != nil {
		return 0, err
	}

	// Get the ID of the student just inserted
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetById implements the storage.Storage interface
func (s *Mysql) GetStudentById(id int64) (types.Student, error) {
	query, err := s.Db.Prepare("SELECT ID, name, email, age FROM students WHERE id = ?")
	if err != nil {
		return types.Student{}, err
	}
	defer query.Close()

	var student types.Student
	query.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student not found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, err
	}
	return student, nil
}

// GetAllStudents implents the storage.Storage interface
func (s *Mysql) GetAllStudents() ([]types.Student, error) {
	rows, err := s.Db.Query("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := []types.Student{}
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}
