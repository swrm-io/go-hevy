package hevy

import (
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
}

func (w Workout) VolumeKG() float64 {
	volume := float64(0)
	for _, e := range w.Exercises {
		volume = volume + e.VolumeKG()
	}
	return volume
}

func (w Workout) VolumeLB() float64 {
	volume := float64(0)
	for _, e := range w.Exercises {
		volume = volume + e.VolumeLB()
	}
	return volume
}

type Exercise struct {
	Index               int    `json:"index"`                // Index indicating the order of the exercise in the workout / routine.
	Title               string `json:"title"`                // Title of the exercise
	Notes               string `json:"notes"`                // Notes on the exercise
	ExcersiseTemplateID string `json:"exercise_template_id"` // The id of the exercise template. This can be used to fetch the exercise template.
	SupersetID          int    `json:"supersets_id"`         // The id of the superset that the exercise belongs to. A value of null indicates the exercise is not part of a superset.
	Sets                []Set  `json:"sets"`                 // List of sets for the exercise.
}

func (e Exercise) VolumeKG() float64 {
	volume := float64(0)
	for _, s := range e.Sets {
		volume = volume + s.VolumeKG()
	}
	return volume
}

func (e Exercise) VolumeLB() float64 {
	volume := float64(0)
	for _, s := range e.Sets {
		volume = volume + s.VolumeLB()
	}
	return volume
}

// Set of the specifc workout
type Set struct {
	Index           int     `json:"index"`            // Index indicating the order of the set in the workout.
	SetType         SetType `json:"set_type"`         // The type of set.
	WeightKG        float64 `json:"weight_kg"`        // Weight lifted in kilograms.
	Reps            int     `json:"reps"`             // Number of reps logged for the set
	DistanceMeters  float64 `json:"distance_meters"`  // Number of meters logged for the set
	DurationSeconds int     `json:"duration_seconds"` // Number of seconds logged for the set
	RPE             float64 `json:"rpe"`              // RPE (Relative perceived exertion) value logged for the set
}

// Convert KG to LBs
func (s Set) KGtoLB() float64 {
	return math.Round(s.WeightKG * 2.20462262185)
}

// Return Volume in KG
func (s Set) VolumeKG() float64 {
	return float64(s.Reps) * s.WeightKG
}

// Return Volume in LB
func (s Set) VolumeLB() float64 {
	return float64(s.Reps) * s.KGtoLB()
}

type Routine struct {
	ID        uuid.UUID  `json:"id"`         // The routine ID.
	Title     string     `json:"title"`      // The routine title.
	CreatedAt time.Time  `json:"created_at"` // ISO 8601 timestamp of when the routine was created.
	UpdatedAt time.Time  `json:"updated_at"` // ISO 8601 timestamp of when the routine was last updated.
	Exercises []Exercise `json:"exercises"`  // Exercise that belong to the workout.
}

type Event struct {
	EventType EventType `json:"type"`       // The Type of Event
	ID        uuid.UUID `json:"id"`         // When deleted, this references the workout that was removed
	DeletedAt time.Time `json:"deleted_at"` // when the type is deleted, when it was removed
	Workout   Workout   `json:"workout"`    // On an update, output the workout
}
