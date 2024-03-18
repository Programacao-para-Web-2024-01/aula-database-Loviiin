package db

import (
	"errors"
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
	m  map[int]Student
	mu *sync.RWMutex
}

func NewStudentRepository() *StudentRepository {
	return &StudentRepository{
		m:  make(map[int]Student),
		mu: &sync.RWMutex{},
	}
}

func (sr *StudentRepository) List() ([]Student, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	students := make([]Student, len(sr.m))
	i := 0
	for _, student := range sr.m {
		students[i] = student
	}
	return students, nil
}

func (sr *StudentRepository) Get(id int) (*Student, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	student, ok := sr.m[id]
	if !ok {
		return nil, errors.New("student not found")
	}

	return &student, nil
}

func (sr *StudentRepository) Create(student Student) (int, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	id := len(sr.m) + 1
	sr.m[id] = student

	return id, nil
}

func (sr *StudentRepository) Update(id int, student Student) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.m[id] = student

	return nil
}

func (sr *StudentRepository) Delete(id int) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	delete(sr.m, id)

	return nil
}
