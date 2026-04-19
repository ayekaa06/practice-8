package repository

import "errors"

type MockUserRepository struct {
	GetUserByIDFunc func(id int) (*User, error)
	CreateUserFunc  func(user *User) error
	GetByEmailFunc  func(email string) (*User, error)
	UpdateUserFunc  func(user *User) error
	DeleteUserFunc  func(id int) error

	GetUserByIDCalls []int
	CreateUserCalls  []*User
	GetByEmailCalls  []string
	UpdateUserCalls  []*User
	DeleteUserCalls  []int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) GetUserByID(id int) (*User, error) {
	m.GetUserByIDCalls = append(m.GetUserByIDCalls, id)
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(id)
	}
	return nil, errors.New("GetUserByID not configured")
}

func (m *MockUserRepository) CreateUser(user *User) error {
	m.CreateUserCalls = append(m.CreateUserCalls, user)
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(user)
	}
	return errors.New("CreateUser not configured")
}

func (m *MockUserRepository) GetByEmail(email string) (*User, error) {
	m.GetByEmailCalls = append(m.GetByEmailCalls, email)
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, errors.New("GetByEmail not configured")
}

func (m *MockUserRepository) UpdateUser(user *User) error {
	m.UpdateUserCalls = append(m.UpdateUserCalls, user)
	if m.UpdateUserFunc != nil {
		return m.UpdateUserFunc(user)
	}
	return errors.New("UpdateUser not configured")
}

func (m *MockUserRepository) DeleteUser(id int) error {
	m.DeleteUserCalls = append(m.DeleteUserCalls, id)
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return errors.New("DeleteUser not configured")
}
