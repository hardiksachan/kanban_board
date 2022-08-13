package native

import (
	"github.com/hardiksachan/kanban_board/backend/internal/users"
	"github.com/hardiksachan/kanban_board/backend/internal/users/core/domain"
	"strconv"
)

type UserStore struct {
	users []*domain.User
	id    int
}

func NewUserStore() *UserStore {
	return &UserStore{id: 0}
}

func (s *UserStore) Insert(user *domain.User) (*domain.User, error) {
	s.id++

	user.ID = strconv.Itoa(s.id)
	s.users = append(s.users, user)

	return user, nil
}

func (s *UserStore) Update(user *domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *UserStore) Remove(user *domain.User) error {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}

	return &users.Error{
		Code:    users.ENOTFOUND,
		Message: "user does not exist",
		Op:      "UserStore.Remove",
	}
}

func (s *UserStore) FindById(userId string) (*domain.User, error) {
	for _, user := range s.users {
		if userId == user.ID {
			return user, nil
		}
	}

	return nil, &users.Error{
		Code:    users.ENOTFOUND,
		Message: "user does not exist",
		Op:      "UserStore.FindById",
	}
}

func (s *UserStore) FindByEmail(email string) (*domain.User, error) {
	for _, user := range s.users {
		if email == user.Email {
			return user, nil
		}
	}

	return nil, &users.Error{
		Code:    users.ENOTFOUND,
		Message: "user does not exist",
		Op:      "UserStore.FindByEmail",
	}
}

func (s *UserStore) CheckByEmail(email string) (bool, error) {
	for _, user := range s.users {
		if email == user.Email {
			return true, nil
		}
	}

	return false, nil
}
