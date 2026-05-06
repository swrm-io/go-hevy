package hevy

import (
	"encoding/json"
	"math"
	"time"

	"github.com/google/uuid"
)

type SetType string
type EventType string

const (
	NormalSet  SetType = "normal"
	WarmupSet  SetType = "warmup"
	DropSet    SetType = "dropset"
	FailureSet SetType = "failure"

	UpdatedEvent EventType = "updated"
	DeletedEvent EventType = "deleted"
)

// Base Classes

type Workout struct {
	ID          uuid.UUID  `json:"id"`          // The workout ID.
	Title       string     `json:"title"`       // The workout title.
	Description string     `json:"description"` // The workout description.
	StartTime   time.Time  `json:"start_time"`  // ISO 8601 timestamp of when the workout was recorded to have started.
	EndTime     time.Time  `json:"end_time"`    // ISO 8601 timestamp of when the workout was recorded to have ended.
	CreatedAt   time.Time  `json:"created_at"`  // ISO 8601 timestamp of when the workout was created.
	UpdatedAt   time.Time  `json:"updated_at"`  // ISO 8601 timestamp of when the workout was last updated.
	Exercises   []Exercise `json:"exercises"`   // Exercise that belong to the workout.
	VolumeKG    float64    `json:"-"`           // Volume of workout in KG
	VolumeLB    float64    `json:"-"`           // Volume of the workout in LB
}

func (w *Workout) UnmarshalJSON(b []byte) error {
	type mask Workout
	var base mask

	err := json.Unmarshal(b, &base)
	if err != nil {
		return err
	}

	vkg := float64(0)
	vlb := float64(0)
	for _, s := range base.Exercises {
		vkg = vkg + s.VolumeKG
		vlb = vlb + s.VolumeLB
	}

	w.ID = base.ID
	w.Title = base.Title
	w.Description = base.Description
	w.StartTime = base.StartTime
	w.EndTime = base.EndTime
	w.UpdatedAt = base.UpdatedAt
	w.Exercises = base.Exercises
	w.VolumeKG = vkg
	w.VolumeLB = vlb

	return nil
}

type Exercise struct {
	Index               int     `json:"index"`                // Index indicating the order of the exercise in the workout / routine.
	Title               string  `json:"title"`                // Title of the exercise
	Notes               string  `json:"notes"`                // Notes on the exercise
	ExcersiseTemplateID string  `json:"exercise_template_id"` // The id of the exercise template. This can be used to fetch the exercise template.
	SupersetID          int     `json:"supersets_id"`         // The id of the superset that the exercise belongs to. A value of null indicates the exercise is not part of a superset.
	Sets                []Set   `json:"sets"`                 // List of sets for the exercise.
	VolumeKG            float64 `json:"-"`                    // Volume of exercise in KG
	VolumeLB            float64 `json:"-"`                    // Volume of the exercise in LB
}

func (e *Exercise) UnmarshalJSON(b []byte) error {
	type mask Exercise
	var base mask

	err := json.Unmarshal(b, &base)
	if err != nil {
		return err
	}

	vkg := float64(0)
	vlb := float64(0)
	for _, s := range base.Sets {
		vkg = vkg + s.VolumeKG
		vlb = vlb + s.VolumeLB
	}

	e.Index = base.Index
	e.Title = base.Title
	e.Notes = base.Notes
	e.ExcersiseTemplateID = base.ExcersiseTemplateID
	e.SupersetID = base.SupersetID
	e.Sets = base.Sets
	e.VolumeKG = vkg
	e.VolumeLB = vlb

	return nil
}

// Set of the specifc workout
type Set struct {
	Index           int     `json:"index"`            // Index indicating the order of the set in the workout.
	SetType         SetType `json:"set_type"`         // The type of set.
	WeightKG        float64 `json:"weight_kg"`        // Weight lifted in kilograms.
	WeightLB        float64 `json:"-"`                // Weight lifted in pounds (computed)
	VolumeKG        float64 `json:"-"`                // Total Volume of the set in KG
	VolumeLB        float64 `json:"-"`                // Total volume of the set in LB
	Reps            int     `json:"reps"`             // Number of reps logged for the set
	DistanceMeters  float64 `json:"distance_meters"`  // Number of meters logged for the set
	DurationSeconds int     `json:"duration_seconds"` // Number of seconds logged for the set
	RPE             float64 `json:"rpe"`              // RPE (Relative perceived exertion) value logged for the set
}

// UnmarshalJSON unmarshals the given struct, and also computes the
// following fields:
// WeightLB
// VolumeKG
// VolumeLB
func (s *Set) UnmarshalJSON(b []byte) error {
	type mask Set
	var base mask

	err := json.Unmarshal(b, &base)
	if err != nil {
		return err
	}

	s.Index = base.Index
	s.SetType = base.SetType
	s.WeightKG = base.WeightKG
	s.WeightLB = math.Round(s.WeightKG * 2.20462262185)
	s.VolumeKG = s.WeightKG * float64(base.Reps)
	s.VolumeLB = s.WeightLB * float64(base.Reps)
	s.Reps = base.Reps
	s.DistanceMeters = base.DistanceMeters
	s.DurationSeconds = base.DurationSeconds
	s.RPE = base.RPE

	return nil
}

type Routine struct {
	ID        uuid.UUID  `json:"id"`         // The routine ID.
	Title     string     `json:"title"`      // The routine title.
	CreatedAt time.Time  `json:"created_at"` // ISO 8601 timestamp of when the routine was created.
	UpdatedAt time.Time  `json:"updated_at"` // ISO 8601 timestamp of when the routine was last updated.
	Exercises []Exercise `json:"exercises"`  // Exercise that belong to the workout.
}

type RoutineFolder struct {
	ID        int       `json:"id"`         // The routine folder ID.
	Index     int       `json:"index"`      // Index indicating the order of the routine folder.
	Title     string    `json:"title"`      // The routine folder title.
	UpdatedAt time.Time `json:"updated_at"` // ISO 8601 timestamp of when the routine folder was last updated.
	CreatedAt time.Time `json:"created_at"` // ISO 8601 timestamp of when the routine folder was created.
}

type Event struct {
	EventType EventType `json:"type"`       // The Type of Event
	ID        uuid.UUID `json:"id"`         // When deleted, this references the workout that was removed
	DeletedAt time.Time `json:"deleted_at"` // when the type is deleted, when it was removed
	Workout   Workout   `json:"workout"`    // On an update, output the workout
}

type User struct {
	ID   uuid.UUID `json:"id"`   // The user ID.
	Name string    `json:"name"` // The user's name.
	URL  string    `json:"url"`  // The user's profile URL.
}

type BodyMeasurement struct {
	ID             int       `json:"id"`               // The body measurement ID.
	Date           string    `json:"date"`             // The date of the measurement in ISO 8601 format (YYYY-MM-DD).
	CreatedAt      time.Time `json:"created_at"`       // ISO 8601 timestamp of when the body measurement was created.
	WeightKG       float64   `json:"weight_kg"`        // The weight measurement in kilograms.
	WeightLB       float64   `json:"-"`                // The weight measurement in pounds (computed).
	LeanMassKG     float64   `json:"lean_mass_kg"`     // The lean mass measurement
	LeanMassLB     float64   `json:"-"`                // The lean mass measurement in pounds (computed).
	FatPercent     float64   `json:"fat_percent"`      // The fat percentage measurement.
	NeckCM         float64   `json:"neck_cm"`          // The neck measurement in centimeters.
	NeckInches     float64   `json:"-"`                // The neck measurement in inches (computed).
	ShoulderCM     float64   `json:"shoulder_cm"`      // The shoulder measurement in centimeters.
	ShoulderIN     float64   `json:"-"`                // The shoulder measurement in inches (computed).
	ChestCM        float64   `json:"chest_cm"`         // The chest measurement in centimeters.
	ChestIn        float64   `json:"-"`                // The chest measurement in inches (computed).
	LeftBicepCM    float64   `json:"left_bicep_cm"`    // The left bicep measurement in centimeters.
	LeftBicepIN    float64   `json:"-"`                // The left bicep measurement in inches (computed).
	RightBicepCM   float64   `json:"right_bicep_cm"`   // The right bicep measurement in centimeters.
	RightBicepIN   float64   `json:"-"`                // The right bicep measurement in inches (computed).
	LeftForearmCM  float64   `json:"left_forearm_cm"`  // The left forearm measurement in centimeters.
	LeftForearmIN  float64   `json:"-"`                // The left forearm measurement in inches (computed).
	RightForearmCM float64   `json:"right_forearm_cm"` // The right forearm measurement in centimeters.
	RightForearmIN float64   `json:"-"`                // The right forearm measurement in inches (computed).
	AbdomenCM      float64   `json:"abdomen"`          // The abdomen measurement in centimeters.
	AbdomenIN      float64   `json:"-"`                // The abdomen measurement
	WaistCM        float64   `json:"waist"`            // The waist measurement in centimeters.
	WaistIN        float64   `json:"-"`                // The waist measurement in inches (computed).
	HipsCM         float64   `json:"hips"`             // The hips measurement in centimeters.
	HipsIN         float64   `json:"-"`                // The hips measurement in inches (computed).
	LeftThighCM    float64   `json:"left_thigh"`       // The left thigh measurement in centimeters.
	LeftThighIN    float64   `json:"-"`                // The left thigh measurement in inches (computed).
	RightThighCM   float64   `json:"right_thigh"`      // The right thigh measurement in centimeters.
	RightThighIN   float64   `json:"-"`                // The right thigh measurement in inches (computed).
	LeftCalfCM     float64   `json:"left_calf"`        // The left calf measurement in centimeters.
	LeftCalfIN     float64   `json:"-"`                // The left calf measurement in inches (computed).
	RightCalfCM    float64   `json:"right_calf"`       // The right calf measurement in centimeters.
	RightCalfIN    float64   `json:"-"`                // The right calf measurement in inches (computed).
}

// UnmarshalJSON unmarshals the given struct, and also computes the
// following fields:
// WeightLB
// LeanMassLB
// NeckInches
// ShoulderIN
// ChestIn
// LeftBicepIN
// RightBicepIN
// LeftForearmIN
// RightForearmIN
// AbdomenIN
// WaistIN
// HipsIN
// LeftThighIN
// RightThighIN
// LeftCalfIN
// RightCalfIN
func (bm *BodyMeasurement) UnmarshalJSON(b []byte) error {
	type mask BodyMeasurement
	var base mask

	err := json.Unmarshal(b, &base)
	if err != nil {
		return err
	}

	bm.ID = base.ID
	bm.Date = base.Date
	bm.CreatedAt = base.CreatedAt
	bm.WeightKG = base.WeightKG
	bm.WeightLB = bm.WeightKG * 2.20462262185
	bm.LeanMassKG = base.LeanMassKG
	bm.LeanMassLB = bm.LeanMassKG * 2.20462262185
	bm.FatPercent = base.FatPercent
	bm.NeckCM = base.NeckCM
	bm.NeckInches = bm.NeckCM * 0.3937007874
	bm.ShoulderCM = base.ShoulderCM
	bm.ShoulderIN = bm.ShoulderCM * 0.3937007874
	bm.ChestCM = base.ChestCM
	bm.ChestIn = bm.ChestCM * 0.3937007874
	bm.LeftBicepCM = base.LeftBicepCM
	bm.LeftBicepIN = bm.LeftBicepCM * 0.3937007874
	bm.RightBicepCM = base.RightBicepCM
	bm.RightBicepIN = bm.RightBicepCM * 0.3937007874
	bm.LeftForearmCM = base.LeftForearmCM
	bm.LeftForearmIN = bm.LeftForearmCM * 0.3937007874
	bm.RightForearmCM = base.RightForearmCM
	bm.RightForearmIN = bm.RightForearmCM * 0.3937007874
	bm.AbdomenCM = base.AbdomenCM
	bm.AbdomenIN = bm.AbdomenCM * 0.3937007874
	bm.WaistCM = base.WaistCM
	bm.WaistIN = bm.WaistCM * 0.3937007874
	bm.HipsCM = base.HipsCM
	bm.HipsIN = bm.HipsCM * 0.3937007874
	bm.LeftThighCM = base.LeftThighCM
	bm.LeftThighIN = bm.LeftThighCM * 0.3937007874
	bm.RightThighCM = base.RightThighCM
	bm.RightThighIN = bm.RightThighCM * 0.3937007874
	bm.LeftCalfCM = base.LeftCalfCM
	bm.LeftCalfIN = bm.LeftCalfCM * 0.3937007874
	bm.RightCalfCM = base.RightCalfCM
	bm.RightCalfIN = bm.RightCalfCM * 0.3937007874

	return nil
}
