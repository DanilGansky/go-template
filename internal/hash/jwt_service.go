package hash

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	Generate(id int) (string, error)
	Validate(token string) (int, bool)
}

type jwtService struct {
	secret, issuer string
}

func NewTokenService(secret, issuer string) (TokenService, error) {
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}
	return &jwtService{secret: secret, issuer: issuer}, nil
}

func (s *jwtService) Generate(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // TODO: make it customizable
		"user_id": id,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iss":     s.issuer,
	})

	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) Validate(token string) (int, bool) {
	defer func() {
		recover() //nolint
	}()

	validatedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, errors.New("token does not valid")
		}
		return []byte(s.secret), nil
	})

	if validatedToken.Claims != nil {
		claims := validatedToken.Claims.(jwt.MapClaims)
		if !claims.VerifyIssuer(s.issuer, true) {
			return 0, false
		}
		return int(claims["user_id"].(float64)), true
	}
	return 0, false
}
