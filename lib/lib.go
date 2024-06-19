package lib

import (
	"challenge-goapi/entity"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// func hash password with md5 digest
func HashMD5(password string) string {
	hash := md5.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateToken(employee *entity.Employee) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["uid"] = employee.Id
	claim["name"] = employee.Name
	claim["email"] = employee.Email
	claim["permission"] = employee.Department
	claim["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte("@Enigma2024"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// validate price
func ValidatePrice(n int) error {
	if n < 0 {
		return errors.New("price must be greater than or equal to 0")
	}
	return nil
}

func ValidateDepartment(s string) error {
	switch s {
	case "admin", "cashier", "manager":
		return nil
	default:
		return errors.New("invalid department")
	}
}

func ValidateEmail(s string) error {
	if !strings.Contains(s, "@") {
		return errors.New("invalid email")
	}
	return nil
}

func ValidatePhoneNumber(s string) error {
	if len(s) != 12 {
		return errors.New("invalid phone number")
	}
	return nil
}

func ValidateUnit(s string) error {
	switch s {
	case "unit", "kg", "buah", "pasang":
		return nil
	default:
		return errors.New("invalid unit")
	}
}

func ValidateString(s string) error {
	if len(s) < 3 {
		return errors.New("invalid name " + s + "")
	}
	return nil
}

func ValidatePassword(s string) error {
	if len(s) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !HasUppercase(s) {
		return errors.New("password must contains uppercase letter")
	}
	if !HasLowercase(s) {
		return errors.New("password must contains lowercase letter")
	}
	if !HasDigit(s) {
		return errors.New("password must contain digit")
	}
	if !HasSpecialChar(s) {
		return errors.New("password must contain special character")
	}
	return nil
}

// / helper functions to check password criteria
func HasUppercase(s string) bool {
	for _, r := range s {
		if 'A' <= r && r <= 'Z' {
			return true
		}
	}
	return false
}

func HasLowercase(s string) bool {
	for _, r := range s {
		if 'a' <= r && r <= 'z' {
			return true
		}
	}
	return false
}

func HasDigit(s string) bool {
	for _, r := range s {
		if '0' <= r && r <= '9' {
			return true
		}
	}
	return false
}

func HasSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()-_=+{}[]|:;<>,.?/~`"
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}
