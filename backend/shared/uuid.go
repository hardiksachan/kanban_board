package shared

import (
	gofrsUUID "github.com/gofrs/uuid"
	googleUUID "github.com/google/uuid"
)

func GetUUIDFromString(uuid string) (*googleUUID.UUID, error) {
	fromString, err := gofrsUUID.FromString(uuid)
	if err != nil {
		return nil, err
	}

	validUUID := googleUUID.UUID(fromString)

	return &validUUID, nil
}
