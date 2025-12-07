package utm

import "github.com/google/uuid"

type OperationalIntent struct {
	Reference OperationalIntentReference `json:"reference"`
	Details   OperationalIntentDetails   `json:"details"`
}

type OperationalIntentReference struct {
	ID uuid.UUID `json:"id"`
}

type OperationalIntentDetails struct {
	Volumes           []Volume4d `json:"volumes"`
	OffNominalVolumes []Volume4d `json:"off_nominal_volumes"`
	Priority          uint16     `json:"priority"`
}

type Volume4d struct {
}
