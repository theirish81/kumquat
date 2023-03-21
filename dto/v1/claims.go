package v1

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	Access    string   `json:"access"`
	Sequences []string `json:"sequences"`
	Exp       int64    `json:"exp"`
}

func (c *Claims) HasSequence(seq string) bool {
	for _, availableSequence := range c.Sequences {
		if seq == availableSequence {
			return true
		}
	}
	return false
}

func (c *Claims) IsAccessAll() bool {
	return c.Access == "all"
}

func (c *Claims) CanAccess(seq string) bool {
	return c.IsAccessAll() || c.HasSequence(seq)
}
