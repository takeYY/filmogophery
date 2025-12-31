package hasher

import (
	"filmogophery/internal/pkg/constant"

	"golang.org/x/crypto/bcrypt"
)

type (
	IPasswordHasher interface {
		Hash(password string) (constant.PasswordHasher, error)
		Compare(hashedPassword constant.PasswordHasher, password string) error
	}
	bcryptHasher struct {
		cost int
	}
)

func NewBcryptHasher() IPasswordHasher {
	return &bcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *bcryptHasher) Hash(password string) (constant.PasswordHasher, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return constant.PasswordHasher(string(bytes)), err
}

func (h *bcryptHasher) Compare(hashedPassword constant.PasswordHasher, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
