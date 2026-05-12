package hevy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/swrm-io/go-hevy"
)

func TestBodyMeasurementImperial(t *testing.T) {
	m := hevy.BodyMeasurement{
		Date:      "2026-02-26",
		WeightKg:  ptr(100.0),
		LeanMassKg: ptr(80.0),
		FatPercent: ptr(20.0),
		NeckCm:    ptr(40.0),
		WaistCm:   ptr(90.0),
	}

	imp := m.Imperial()

	assert.Equal(t, "2026-02-26", imp.Date)

	// FatPercent is unit-less — should pass through unchanged
	require.NotNil(t, imp.FatPercent)
	assert.Equal(t, 20.0, *imp.FatPercent)

	// 100 kg → 220.462 lbs
	require.NotNil(t, imp.WeightLbs)
	assert.InDelta(t, 220.462, *imp.WeightLbs, 0.001)

	// 80 kg → 176.370 lbs
	require.NotNil(t, imp.LeanMassLbs)
	assert.InDelta(t, 176.370, *imp.LeanMassLbs, 0.001)

	// 40 cm → 15.748 in
	require.NotNil(t, imp.NeckIn)
	assert.InDelta(t, 15.748, *imp.NeckIn, 0.001)

	// 90 cm → 35.433 in
	require.NotNil(t, imp.WaistIn)
	assert.InDelta(t, 35.433, *imp.WaistIn, 0.001)

	// nil fields should remain nil
	assert.Nil(t, imp.ChestIn)
	assert.Nil(t, imp.HipsIn)
}
