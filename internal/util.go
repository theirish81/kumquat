package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type Proto map[string]string

func PrototypeCheck(config map[string]any, schema string) error {
	res, err := gojsonschema.Validate(gojsonschema.NewStringLoader(schema), gojsonschema.NewRawLoader(config))
	if err != nil {
		return err
	}
	if !res.Valid() {
		fmt.Println(res.Errors())
		return errors.New("invalid sequence configuration")
	}
	return nil
}

// SetDefault checks whether the provided map contain the provided key. If it does not, it will set it with the
// provided default value
func SetDefault(config map[string]any, key string, value any) map[string]any {
	if _, ok := config[key]; !ok {
		config[key] = value
	}
	return config
}

func GetSequencePath(seq string) (string, error) {
	if !IsSequenceAllowed(seq) {
		return "", errors.New("sequence name not allowed")
	}
	sequencePath := os.Getenv("SEQUENCES_PATH")
	if len(sequencePath) == 0 {
		sequencePath = "etc/sequences"
	}
	return path.Join(sequencePath, seq+".yaml"), nil
}

func GetUsersPath() string {
	usersPath := os.Getenv("USERS_PATH")
	if len(usersPath) == 0 {
		usersPath = "etc/users.yaml"
	}
	return usersPath
}

func IsSequenceAllowed(seq string) bool {
	for _, char := range []string{".", "\\", "/", "|", "{", "[", "$", "\"", "%", "(", "=", "?"} {
		if strings.Contains(seq, char) {
			return false
		}
	}
	return true
}
