package auth

import (
	"os"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

// Users is a map of User where the key is the username, and the value represent the details
type Users map[string]User

// Authenticate will return a user if username and password match our users file
func (u Users) Authenticate(username string, password string) *User {
	if ux, ok := u[username]; ok && bcrypt.CompareHashAndPassword([]byte(ux.Password), []byte(password)) == nil {
		return &ux
	}
	return nil
}

// User is one user configuration
type User struct {
	Password  string   `yaml:"password"`
	Access    string   `yaml:"access"`
	Sequences []string `yaml:"sequences"`
}

// LoadUsers loads the users from file
func LoadUsers(filePath string) (Users, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var users Users
	err = yaml.Unmarshal(data, &users)
	return users, err
}
