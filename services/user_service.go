package services

import (
	"errors"
	"sync"
	"sync/atomic"

	"go-microservice-highload/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	mu    sync.RWMutex
	seq   atomic.Int64
	users map[int]models.User
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]models.User),
	}
}

func (s *UserService) Create(u models.User) (models.User, error) {
	if err := u.Validate(); err != nil {
		return models.User{}, err
	}

	id := int(s.seq.Add(1))
	u.ID = id

	s.mu.Lock()
	s.users[id] = u
	s.mu.Unlock()

	return u, nil
}

func (s *UserService) GetAll() []models.User {
	s.mu.RLock()
	out := make([]models.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	s.mu.RUnlock()
	return out
}

func (s *UserService) GetByID(id int) (models.User, bool) {
	s.mu.RLock()
	u, ok := s.users[id]
	s.mu.RUnlock()
	return u, ok
}

func (s *UserService) Update(id int, u models.User) (models.User, error) {
	if err := u.Validate(); err != nil {
		return models.User{}, err
	}

	s.mu.Lock()
	_, ok := s.users[id]
	if !ok {
		s.mu.Unlock()
		return models.User{}, ErrUserNotFound
	}
	u.ID = id
	s.users[id] = u
	s.mu.Unlock()

	return u, nil
}

func (s *UserService) Delete(id int) bool {
	s.mu.Lock()
	_, ok := s.users[id]
	if ok {
		delete(s.users, id)
	}
	s.mu.Unlock()
	return ok
}
