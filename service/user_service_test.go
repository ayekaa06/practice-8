package service

import (
	"errors"
	"practice-8/repository"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", msg)
		return
	}
	if msg != "" && err.Error() != msg {
		t.Errorf("expected error %q, got %q", msg, err.Error())
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestGetUserByID(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.GetUserByIDFunc = func(id int) (*repository.User, error) {
		if id == 1 {
			return user, nil
		}
		return nil, errors.New("not found")
	}

	svc := NewUserService(mockRepo)
	result, err := svc.GetUserByID(1)
	assertNoError(t, err)
	if result != user {
		t.Errorf("expected user %v, got %v", user, result)
	}
}

func TestCreateUser(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.CreateUserFunc = func(u *repository.User) error { return nil }

	svc := NewUserService(mockRepo)
	err := svc.CreateUser(user)
	assertNoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	existing := &repository.User{ID: 2, Name: "Alice", Email: "alice@example.com"}
	mockRepo.GetByEmailFunc = func(email string) (*repository.User, error) {
		return existing, nil
	}

	svc := NewUserService(mockRepo)
	err := svc.RegisterUser(&repository.User{Name: "Bob"}, "alice@example.com")
	assertError(t, err, "user with this email already exists")
}

func TestRegisterUser_NewUserSuccess(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.GetByEmailFunc = func(email string) (*repository.User, error) {
		return nil, nil
	}
	mockRepo.CreateUserFunc = func(u *repository.User) error { return nil }

	svc := NewUserService(mockRepo)
	newUser := &repository.User{Name: "Bob", Email: "bob@example.com"}
	err := svc.RegisterUser(newUser, "bob@example.com")
	assertNoError(t, err)
	if len(mockRepo.CreateUserCalls) != 1 {
		t.Errorf("expected CreateUser to be called once, got %d", len(mockRepo.CreateUserCalls))
	}
}

func TestRegisterUser_RepoErrorOnCreate(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.GetByEmailFunc = func(email string) (*repository.User, error) {
		return nil, nil
	}
	mockRepo.CreateUserFunc = func(u *repository.User) error {
		return errors.New("db connection lost")
	}

	svc := NewUserService(mockRepo)
	err := svc.RegisterUser(&repository.User{Name: "Carol"}, "carol@example.com")
	assertError(t, err, "db connection lost")
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	err := svc.UpdateUserName(1, "")
	assertError(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.GetUserByIDFunc = func(id int) (*repository.User, error) {
		return nil, errors.New("user not found")
	}

	svc := NewUserService(mockRepo)
	err := svc.UpdateUserName(99, "NewName")
	assertError(t, err, "user not found")
}

func TestUpdateUserName_Success(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	user := &repository.User{ID: 2, Name: "OldName"}
	mockRepo.GetUserByIDFunc = func(id int) (*repository.User, error) {
		return user, nil
	}
	mockRepo.UpdateUserFunc = func(u *repository.User) error { return nil }

	svc := NewUserService(mockRepo)
	err := svc.UpdateUserName(2, "NewName")
	assertNoError(t, err)
	if len(mockRepo.UpdateUserCalls) != 1 {
		t.Fatal("UpdateUser was not called")
	}
	assertEqual(t, mockRepo.UpdateUserCalls[0].Name, "NewName")
}

func TestUpdateUserName_UpdateUserFails(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	user := &repository.User{ID: 3, Name: "OldName"}
	mockRepo.GetUserByIDFunc = func(id int) (*repository.User, error) {
		return user, nil
	}
	mockRepo.UpdateUserFunc = func(u *repository.User) error {
		return errors.New("update failed")
	}

	svc := NewUserService(mockRepo)
	err := svc.UpdateUserName(3, "NewName")
	assertError(t, err, "update failed")
}

func TestDeleteUser_AttemptDeleteAdmin(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	svc := NewUserService(mockRepo)
	err := svc.DeleteUser(1)
	assertError(t, err, "it is not allowed to delete admin user")
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.DeleteUserFunc = func(id int) error { return nil }

	svc := NewUserService(mockRepo)
	err := svc.DeleteUser(5)
	assertNoError(t, err)
	if len(mockRepo.DeleteUserCalls) != 1 {
		t.Fatal("DeleteUser was not called")
	}
	assertEqual(t, mockRepo.DeleteUserCalls[0], 5)
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.DeleteUserFunc = func(id int) error {
		return errors.New("repo error on delete")
	}

	svc := NewUserService(mockRepo)
	err := svc.DeleteUser(3)
	assertError(t, err, "repo error on delete")
}
