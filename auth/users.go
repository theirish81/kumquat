package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
	"os"
)

type Users map[string]User

func (u Users) Authenticate(username string, password string) *User {
	if ux, ok := u[username]; ok && bcrypt.CompareHashAndPassword([]byte(ux.Password), []byte(password)) == nil {
		return &ux
	}
	return nil
}

type User struct {
	Password  string   `yaml:"password"`
	Access    string   `yaml:"access"`
	Sequences []string `yaml:"sequences"`
}

func LoadUsers(filePath string) (Users, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var users Users
	err = yaml.Unmarshal(data, &users)
	return users, err
}
