package uuidx

import "github.com/google/uuid"

func Parse(id string) (uuid.UUID, error) { return uuid.Parse(id) }
