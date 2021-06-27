package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Encrypt(s string) (string, error)
	Compare(hashed string, plain string) bool
}

type service struct {
	bcryptCost int
}

func NewService(bcryptCost int) Service {
	return &service{bcryptCost: bcryptCost}
}

func (svc *service) Encrypt(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), svc.bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (svc *service) Compare(hashed string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}
