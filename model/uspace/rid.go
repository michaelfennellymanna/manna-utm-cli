package uspace

// Telemetry is equivalent to https://github.com/m4a3/manna-utm/blob/persistence/src/main/java/manna/aero/utm/model/manna/MannaUspaceTelemetry.java
type Telemetry struct {
	Altitude      float64 `json:"altitude"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lng"`
	Heading       float64 `json:"heading"`
	Speed         float64 `json:"speed"`
	VerticalSpeed float64 `json:"vertical_speed"`
	TimeMeasured  int64   `json:"time_measured"`
	Mode          float64 `json:"mode"`
	Armed         bool    `json:"armed"`
}
