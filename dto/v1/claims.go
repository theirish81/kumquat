package v1

import "github.com/golang-jwt/jwt/v5"

// Claims are our JWT claims
type Claims struct {
	jwt.RegisteredClaims
	Access    string   `json:"access"`
	Sequences []string `json:"sequences"`
	Exp       int64    `json:"exp"`
}

// HasSequence returns true if the user has specific access to a sequence
func (c *Claims) HasSequence(seq string) bool {
	for _, availableSequence := range c.Sequences {
		if seq == availableSequence {
			return true
		}
	}
	return false
}

// IsAccessAll returns true if the user has access to all sequences
func (c *Claims) IsAccessAll() bool {
	return c.Access == "all"
}

// CanAccess returns true if the user has access to a specific sequence, whether it's "access=all" or has an explicit
// sequence set
func (c *Claims) CanAccess(seq string) bool {
	return c.IsAccessAll() || c.HasSequence(seq)
}
