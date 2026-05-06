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
