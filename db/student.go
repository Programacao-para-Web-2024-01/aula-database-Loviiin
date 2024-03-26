package db

import (
	"database/sql"
	"sync"
)

type Student struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type StudentRepository struct {
	db *sql.DB
	m  map[int]Student
	mu *sync.RWMutex
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{
		db: db,
	}
}

func (sr *StudentRepository) List() ([]Student, error) {
	rows, err := sr.db.Query(`SELECT id, name, age, email, phone FROM students`)
	if err != nil {
		return nil, err
	}

	var students []Student

	for rows.Next() {
		var student Student
		err = rows.Scan(&student.Id, &student.Name, &student.Age, &student.Email, &student.Phone)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	rows.Close()

	return students, nil
}

func (sr *StudentRepository) Get(id int) (*Student, error) {
	row := sr.db.QueryRow(`
		SELECT id, name, age, email, phone
		FROM students
		WHERE id = ?`, id)

	var student Student
	err := row.Scan(&student.Id, &student.Name, &student.Age, &student.Email, &student.Phone)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (sr *StudentRepository) Create(student Student) (int, error) {
	result, err := sr.db.Exec(`
		INSERT INTO students (name, age, email, phone)
		VALUES (?, ?, ?, ?)`,
		student.Name, student.Age, student.Email, student.Phone)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sr *StudentRepository) Update(id int, student Student) error {
	_, err := sr.db.Exec(`
        UPDATE students
        SET name=?, age=?, email=?, phone=?
        WHERE id=?`,
		student.Name, student.Age, student.Email, student.Phone, id)
	return err
}
func (sr *StudentRepository) Delete(id int) error {
	_, err := sr.db.Exec(`
	DELETE from students
	WHERE id=?`, id)
	return err
}
