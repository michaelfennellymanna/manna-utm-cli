package utm

import (
	"time"

	"github.com/google/uuid"
)

type OperationalIntentTelemetry struct {
	OperationalIntentId      uuid.UUID `json:"operational_intent_id"`
	Telemetry                Telemetry `json:"telemetry"`
	NextTelemetryOpportunity Time      `json:"next_telemetry_opportunity"`
}

type OperationalIntentTelemetryJson struct {
	OperationalIntentId      string    `json:"operational_intent_id"`
	Telemetry                Telemetry `json:"telemetry"`
	NextTelemetryOpportunity timeJson  `json:"next_telemetry_opportunity"`
}

func (oit OperationalIntentTelemetry) ToJson() OperationalIntentTelemetryJson {
	// conver the uuid to string
	return OperationalIntentTelemetryJson{
		OperationalIntentId:      oit.OperationalIntentId.String(),
		Telemetry:                oit.Telemetry,
		NextTelemetryOpportunity: oit.NextTelemetryOpportunity.toJson(),
	}
}

type Telemetry struct {
	TimeMeasured Time     `json:"time_measured"`
	Position     Position `json:"position"`
	Velocity     Velocity `json:"velocity"`
}

type Position struct {
	Longitude    float64  `json:"longitude"`
	Latitude     float64  `json:"latitude"`
	AccuracyH    string   `json:"accuracy_h"`
	AccuracyV    string   `json:"accuracy_v"`
	Extrapolated bool     `json:"extrapolated"`
	Altitude     Altitude `json:"altitude"`
}

type Altitude struct {
	Value     int64  `json:"value"`
	Reference string `json:"reference"`
	Units     string `json:"units"`
}

type Velocity struct {
	Speed      float64 `json:"speed"`
	UnitsSpeed string  `json:"units_speed"`
	Track      int     `json:"track"`
}

type Time struct {
	Value  time.Time `json:"value"`
	Format string    `json:"format"`
}

func (t Time) toJson() timeJson {
	return timeJson{
		Value:  t.Value.String(),
		Format: t.Format,
	}
}

type timeJson struct {
	Value  string `json:"value"`
	Format string `json:"format"`
}
