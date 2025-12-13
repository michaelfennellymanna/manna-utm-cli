package virtual_uspace

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"manna.aero/manna.utm.cli/pkg/config"
)

func TestVirtualSeriesToLinearSeries(t *testing.T) {
	// create a subseries from two points

	// get the series of GeoJson points from the subseries

}

func TestTelemetrySubSeries(t *testing.T) {
	// create a subseries from two points

	// get the series of GeoJson points from the subseries

}

func TestVirtualTelemetrySeries_IncreaseDetail(t *testing.T) {
	appCnf, err := config.LoadConfig("/Users/michael.fennelly/projects/manna-utm-cli/config.yaml")
	assert.NoError(t, err)

	oicnf, err := appCnf.GetOperationalIntentConfigByName("SWITZERLAND1")
	assert.NoError(t, err)

	voi := NewOperationalIntentManager(oicnf, 10)

	fc := voi.GeoJson(false, false, true)

	jsonData, err := json.MarshalIndent(fc, "", "   ")
	assert.NoError(t, err)

	err = os.WriteFile("./test.geojson", jsonData, os.ModePerm)
}
