package internal

import (
	"errors"
	"os"
	"path"
	"reflect"
)

type Proto map[string]string

const TYPE_STRING = "string"
const TYPE_MAP = "map[string]interface {}"
const TYPE_INT = "int"
const TYPE_ARRAY = "[]string"
const TYPE_BOOL = "bool"

// PrototypeCheck verifies that the provided config map matches the proposed protocol
func PrototypeCheck(config map[string]any, proto Proto) error {
	for k, oType := range proto {
		if val, ok := config[k]; ok {
			if reflect.TypeOf(val).String() != oType {
				return errors.New("field " + k + " is not the right type")
			}
		} else {
			return errors.New("required field " + k + " was not provided")
		}
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

func GetSequencePath(seq string) string {
	etcPath := os.Getenv("SEQUENCES_PATH")
	if len(etcPath) == 0 {
		etcPath = "etc/sequences"
	}
	return path.Join(etcPath, seq+".yaml")
}

func GetUsersPath() string {
	usersPath := os.Getenv("USERS_PATH")
	if len(usersPath) == 0 {
		usersPath = "etc/users.yaml"
	}
	return usersPath
}
