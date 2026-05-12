package hevy

const (
	kgToLbs    = 2.20462262185
	cmToInches = 0.393700787402
)

func convertPtr(v *float64, factor float64) *float64 {
	if v == nil {
		return nil
	}
	result := *v * factor
	return &result
}

// BodyMeasurementImperial holds a BodyMeasurement with weights in pounds and
// lengths in inches.
type BodyMeasurementImperial struct {
	Date           string
	WeightLbs      *float64
	LeanMassLbs    *float64
	FatPercent     *float64
	NeckIn         *float64
	ShoulderIn     *float64
	ChestIn        *float64
	LeftBicepIn    *float64
	RightBicepIn   *float64
	LeftForearmIn  *float64
	RightForearmIn *float64
	AbdomenIn      *float64
	WaistIn        *float64
	HipsIn         *float64
	LeftThighIn    *float64
	RightThighIn   *float64
	LeftCalfIn     *float64
	RightCalfIn    *float64
}

// Imperial converts the measurement to imperial units (pounds and inches).
func (m BodyMeasurement) Imperial() BodyMeasurementImperial {
	return BodyMeasurementImperial{
		Date:           m.Date,
		WeightLbs:      convertPtr(m.WeightKg, kgToLbs),
		LeanMassLbs:    convertPtr(m.LeanMassKg, kgToLbs),
		FatPercent:     m.FatPercent,
		NeckIn:         convertPtr(m.NeckCm, cmToInches),
		ShoulderIn:     convertPtr(m.ShoulderCm, cmToInches),
		ChestIn:        convertPtr(m.ChestCm, cmToInches),
		LeftBicepIn:    convertPtr(m.LeftBicepCm, cmToInches),
		RightBicepIn:   convertPtr(m.RightBicepCm, cmToInches),
		LeftForearmIn:  convertPtr(m.LeftForearmCm, cmToInches),
		RightForearmIn: convertPtr(m.RightForearmCm, cmToInches),
		AbdomenIn:      convertPtr(m.AbdomenCm, cmToInches),
		WaistIn:        convertPtr(m.WaistCm, cmToInches),
		HipsIn:         convertPtr(m.HipsCm, cmToInches),
		LeftThighIn:    convertPtr(m.LeftThighCm, cmToInches),
		RightThighIn:   convertPtr(m.RightThighCm, cmToInches),
		LeftCalfIn:     convertPtr(m.LeftCalfCm, cmToInches),
		RightCalfIn:    convertPtr(m.RightCalfCm, cmToInches),
	}
}
