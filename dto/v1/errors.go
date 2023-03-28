package v1

// Error is the error DTO
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
