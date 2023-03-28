package auth

import "testing"

func TestUsers_Authenticate(t *testing.T) {
	users, _ := LoadUsers("../etc/users.yaml")
	user := users.Authenticate("theirish81", "foobar")
	if user == nil {
		t.Error("authentication process did not work")
	}
	user = users.Authenticate("theirish81", "bananas")
	if user != nil {
		t.Error("authentication returned a user with wrong password")
	}
}
