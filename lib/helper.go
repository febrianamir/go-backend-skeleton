package lib

import "github.com/google/uuid"

func GenerateUUID() string {
	return uuid.NewString()
}
