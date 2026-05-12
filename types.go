package hevy

import "time"

// SetType represents the type of a workout or routine set.
type SetType string

const (
	SetTypeWarmup  SetType = "warmup"
	SetTypeNormal  SetType = "normal"
	SetTypeFailure SetType = "failure"
	SetTypeDropset SetType = "dropset"
)

// RPE represents Rating of Perceived Exertion values.
type RPE float64

const (
	RPE6  RPE = 6
	RPE7  RPE = 7
	RPE75 RPE = 7.5
	RPE8  RPE = 8
	RPE85 RPE = 8.5
	RPE9  RPE = 9
	RPE95 RPE = 9.5
	RPE10 RPE = 10
)

// CustomExerciseType represents the tracking type of a custom exercise.
type CustomExerciseType string

const (
	ExerciseTypeWeightReps             CustomExerciseType = "weight_reps"
	ExerciseTypeRepsOnly               CustomExerciseType = "reps_only"
	ExerciseTypeBodyweightReps         CustomExerciseType = "bodyweight_reps"
	ExerciseTypeBodyweightAssistedReps CustomExerciseType = "bodyweight_assisted_reps"
	ExerciseTypeDuration               CustomExerciseType = "duration"
	ExerciseTypeWeightDuration         CustomExerciseType = "weight_duration"
	ExerciseTypeDistanceDuration       CustomExerciseType = "distance_duration"
	ExerciseTypeShortDistanceWeight    CustomExerciseType = "short_distance_weight"
)

// MuscleGroup represents a muscle group category.
type MuscleGroup string

const (
	MuscleGroupAbdominals MuscleGroup = "abdominals"
	MuscleGroupShoulders  MuscleGroup = "shoulders"
	MuscleGroupBiceps     MuscleGroup = "biceps"
	MuscleGroupTriceps    MuscleGroup = "triceps"
	MuscleGroupForearms   MuscleGroup = "forearms"
	MuscleGroupQuadriceps MuscleGroup = "quadriceps"
	MuscleGroupHamstrings MuscleGroup = "hamstrings"
	MuscleGroupCalves     MuscleGroup = "calves"
	MuscleGroupGlutes     MuscleGroup = "glutes"
	MuscleGroupAbductors  MuscleGroup = "abductors"
	MuscleGroupAdductors  MuscleGroup = "adductors"
	MuscleGroupLats       MuscleGroup = "lats"
	MuscleGroupUpperBack  MuscleGroup = "upper_back"
	MuscleGroupTraps      MuscleGroup = "traps"
	MuscleGroupLowerBack  MuscleGroup = "lower_back"
	MuscleGroupChest      MuscleGroup = "chest"
	MuscleGroupCardio     MuscleGroup = "cardio"
	MuscleGroupNeck       MuscleGroup = "neck"
	MuscleGroupFullBody   MuscleGroup = "full_body"
	MuscleGroupOther      MuscleGroup = "other"
)

// EquipmentCategory represents equipment used for an exercise.
type EquipmentCategory string

const (
	EquipmentNone           EquipmentCategory = "none"
	EquipmentBarbell        EquipmentCategory = "barbell"
	EquipmentDumbbell       EquipmentCategory = "dumbbell"
	EquipmentKettlebell     EquipmentCategory = "kettlebell"
	EquipmentMachine        EquipmentCategory = "machine"
	EquipmentPlate          EquipmentCategory = "plate"
	EquipmentResistanceBand EquipmentCategory = "resistance_band"
	EquipmentSuspension     EquipmentCategory = "suspension"
	EquipmentOther          EquipmentCategory = "other"
)

// WorkoutSet represents a single set within a workout exercise.
type WorkoutSet struct {
	Index           int      `json:"index"`
	Type            SetType  `json:"type"`
	WeightKg        *float64 `json:"weight_kg"`
	Reps            *int     `json:"reps"`
	DistanceMeters  *int     `json:"distance_meters"`
	DurationSeconds *int     `json:"duration_seconds"`
	RPE             *float64 `json:"rpe"`
	CustomMetric    *float64 `json:"custom_metric"`
}

// WorkoutExercise represents an exercise within a workout.
type WorkoutExercise struct {
	Index              int          `json:"index"`
	Title              string       `json:"title"`
	Notes              string       `json:"notes"`
	ExerciseTemplateID string       `json:"exercise_template_id"`
	SupersetID         *int         `json:"superset_id"`
	Sets               []WorkoutSet `json:"sets"`
}

// Workout represents a completed workout session.
type Workout struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	RoutineID   string            `json:"routine_id"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedAt   time.Time         `json:"created_at"`
	Exercises   []WorkoutExercise `json:"exercises"`
}

// RepRange represents a rep range for routine sets.
type RepRange struct {
	Start *float64 `json:"start"`
	End   *float64 `json:"end"`
}

// RoutineSet represents a single set within a routine exercise.
type RoutineSet struct {
	Index           int       `json:"index"`
	Type            SetType   `json:"type"`
	WeightKg        *float64  `json:"weight_kg"`
	Reps            *int      `json:"reps"`
	RepRange        *RepRange `json:"rep_range"`
	DistanceMeters  *int      `json:"distance_meters"`
	DurationSeconds *int      `json:"duration_seconds"`
	RPE             *float64  `json:"rpe"`
	CustomMetric    *float64  `json:"custom_metric"`
}

// RoutineExercise represents an exercise within a routine.
type RoutineExercise struct {
	Index              int          `json:"index"`
	Title              string       `json:"title"`
	RestSeconds        int          `json:"rest_seconds"`
	Notes              string       `json:"notes"`
	ExerciseTemplateID string       `json:"exercise_template_id"`
	SupersetID         *int         `json:"superset_id"`
	Sets               []RoutineSet `json:"sets"`
}

// Routine represents a workout routine (template).
type Routine struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	FolderID  *float64          `json:"folder_id"`
	UpdatedAt time.Time         `json:"updated_at"`
	CreatedAt time.Time         `json:"created_at"`
	Exercises []RoutineExercise `json:"exercises"`
}

// RoutineFolder represents a folder that groups routines.
type RoutineFolder struct {
	ID        int       `json:"id"`
	Index     int       `json:"index"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ExerciseTemplate represents a reusable exercise definition.
type ExerciseTemplate struct {
	ID                    string   `json:"id"`
	Title                 string   `json:"title"`
	Type                  string   `json:"type"`
	PrimaryMuscleGroup    string   `json:"primary_muscle_group"`
	SecondaryMuscleGroups []string `json:"secondary_muscle_groups"`
	Equipment             string   `json:"equipment"`
	IsCustom              bool     `json:"is_custom"`
}

// ExerciseHistoryEntry represents a single set entry from exercise history.
type ExerciseHistoryEntry struct {
	WorkoutID          string    `json:"workout_id"`
	WorkoutTitle       string    `json:"workout_title"`
	WorkoutStartTime   time.Time `json:"workout_start_time"`
	WorkoutEndTime     time.Time `json:"workout_end_time"`
	ExerciseTemplateID string    `json:"exercise_template_id"`
	WeightKg           *float64  `json:"weight_kg"`
	Reps               *int      `json:"reps"`
	DistanceMeters     *int      `json:"distance_meters"`
	DurationSeconds    *int      `json:"duration_seconds"`
	RPE                *float64  `json:"rpe"`
	CustomMetric       *float64  `json:"custom_metric"`
	SetType            SetType   `json:"set_type"`
}

// BodyMeasurement represents a body measurement snapshot for a given date.
type BodyMeasurement struct {
	Date           string   `json:"date"`
	WeightKg       *float64 `json:"weight_kg"`
	LeanMassKg     *float64 `json:"lean_mass_kg"`
	FatPercent     *float64 `json:"fat_percent"`
	NeckCm         *float64 `json:"neck_cm"`
	ShoulderCm     *float64 `json:"shoulder_cm"`
	ChestCm        *float64 `json:"chest_cm"`
	LeftBicepCm    *float64 `json:"left_bicep_cm"`
	RightBicepCm   *float64 `json:"right_bicep_cm"`
	LeftForearmCm  *float64 `json:"left_forearm_cm"`
	RightForearmCm *float64 `json:"right_forearm_cm"`
	AbdomenCm      *float64 `json:"abdomen"`
	WaistCm        *float64 `json:"waist"`
	HipsCm         *float64 `json:"hips"`
	LeftThighCm    *float64 `json:"left_thigh"`
	RightThighCm   *float64 `json:"right_thigh"`
	LeftCalfCm     *float64 `json:"left_calf"`
	RightCalfCm    *float64 `json:"right_calf"`
}

// UserInfo contains basic information about the authenticated user.
type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// WorkoutEventType indicates whether a workout event is an update or deletion.
type WorkoutEventType string

const (
	WorkoutEventUpdated WorkoutEventType = "updated"
	WorkoutEventDeleted WorkoutEventType = "deleted"
)

// WorkoutEvent represents a single event from the workout events feed.
// Check Type to determine whether Workout or DeletedAt/DeletedID is populated.
type WorkoutEvent struct {
	Type WorkoutEventType `json:"type"`
	// Populated when Type == WorkoutEventUpdated.
	Workout *Workout `json:"workout,omitempty"`
	// Populated when Type == WorkoutEventDeleted.
	ID        string     `json:"id,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// --- Paginated response types ---

// PaginatedWorkouts is the response from GET /v1/workouts.
type PaginatedWorkouts struct {
	Page      int       `json:"page"`
	PageCount int       `json:"page_count"`
	Workouts  []Workout `json:"workouts"`
}

// WorkoutCount is the response from GET /v1/workouts/count.
type WorkoutCount struct {
	WorkoutCount int `json:"workout_count"`
}

// PaginatedWorkoutEvents is the response from GET /v1/workouts/events.
type PaginatedWorkoutEvents struct {
	Page      int            `json:"page"`
	PageCount int            `json:"page_count"`
	Events    []WorkoutEvent `json:"events"`
}

// PaginatedRoutines is the response from GET /v1/routines.
type PaginatedRoutines struct {
	Page      int       `json:"page"`
	PageCount int       `json:"page_count"`
	Routines  []Routine `json:"routines"`
}

// PaginatedRoutineFolders is the response from GET /v1/routine_folders.
type PaginatedRoutineFolders struct {
	Page           int             `json:"page"`
	PageCount      int             `json:"page_count"`
	RoutineFolders []RoutineFolder `json:"routine_folders"`
}

// PaginatedExerciseTemplates is the response from GET /v1/exercise_templates.
type PaginatedExerciseTemplates struct {
	Page              int                `json:"page"`
	PageCount         int                `json:"page_count"`
	ExerciseTemplates []ExerciseTemplate `json:"exercise_templates"`
}

// ExerciseHistory is the response from GET /v1/exercise_history/{exerciseTemplateId}.
type ExerciseHistory struct {
	ExerciseHistory []ExerciseHistoryEntry `json:"exercise_history"`
}

// PaginatedBodyMeasurements is the response from GET /v1/body_measurements.
type PaginatedBodyMeasurements struct {
	Page         int               `json:"page"`
	PageCount    int               `json:"page_count"`
	Measurements []BodyMeasurement `json:"body_measurements"`
}

// --- Request body types ---

// WorkoutSetInput is a set as sent in create/update workout requests.
type WorkoutSetInput struct {
	Type            SetType  `json:"type"`
	WeightKg        *float64 `json:"weight_kg,omitempty"`
	Reps            *int     `json:"reps,omitempty"`
	DistanceMeters  *int     `json:"distance_meters,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	CustomMetric    *float64 `json:"custom_metric,omitempty"`
	RPE             *float64 `json:"rpe,omitempty"`
}

// WorkoutExerciseInput is an exercise as sent in create/update workout requests.
type WorkoutExerciseInput struct {
	ExerciseTemplateID string            `json:"exercise_template_id"`
	SupersetID         *int              `json:"superset_id,omitempty"`
	Notes              *string           `json:"notes,omitempty"`
	Sets               []WorkoutSetInput `json:"sets"`
}

// WorkoutInput is the workout payload for create/update requests.
type WorkoutInput struct {
	Title       string                 `json:"title"`
	Description *string                `json:"description,omitempty"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
	IsPrivate   *bool                  `json:"is_private,omitempty"`
	Exercises   []WorkoutExerciseInput `json:"exercises"`
}

// RoutineSetInput is a set as sent in create/update routine requests.
type RoutineSetInput struct {
	Type            SetType   `json:"type"`
	WeightKg        *float64  `json:"weight_kg,omitempty"`
	Reps            *int      `json:"reps,omitempty"`
	DistanceMeters  *int      `json:"distance_meters,omitempty"`
	DurationSeconds *int      `json:"duration_seconds,omitempty"`
	CustomMetric    *float64  `json:"custom_metric,omitempty"`
	RepRange        *RepRange `json:"rep_range,omitempty"`
}

// RoutineExerciseInput is an exercise as sent in create/update routine requests.
type RoutineExerciseInput struct {
	ExerciseTemplateID string            `json:"exercise_template_id"`
	SupersetID         *int              `json:"superset_id,omitempty"`
	RestSeconds        *int              `json:"rest_seconds,omitempty"`
	Notes              *string           `json:"notes,omitempty"`
	Sets               []RoutineSetInput `json:"sets"`
}

// RoutineInput is the routine payload for create/update requests.
type RoutineInput struct {
	Title     string                 `json:"title"`
	FolderID  *float64               `json:"folder_id,omitempty"`
	Notes     string                 `json:"notes"`
	Exercises []RoutineExerciseInput `json:"exercises"`
}

// RoutineUpdateInput is the payload for PUT /v1/routines/{routineId}.
// Notes is optional (unlike create).
type RoutineUpdateInput struct {
	Title     string                 `json:"title"`
	FolderID  *float64               `json:"folder_id,omitempty"`
	Notes     *string                `json:"notes,omitempty"`
	Exercises []RoutineExerciseInput `json:"exercises"`
}

// CreateExerciseInput is the payload for POST /v1/exercise_templates.
type CreateExerciseInput struct {
	Title             string             `json:"title"`
	ExerciseType      CustomExerciseType `json:"exercise_type"`
	EquipmentCategory EquipmentCategory  `json:"equipment_category"`
	MuscleGroup       MuscleGroup        `json:"muscle_group"`
	OtherMuscles      []MuscleGroup      `json:"other_muscles,omitempty"`
}

// BodyMeasurementUpdate is the payload for PUT /v1/body_measurements/{date}.
// All fields are optional; omitted fields will be set to null by the API.
type BodyMeasurementUpdate struct {
	WeightKg       *float64 `json:"weight_kg"`
	LeanMassKg     *float64 `json:"lean_mass_kg"`
	FatPercent     *float64 `json:"fat_percent"`
	NeckCm         *float64 `json:"neck_cm"`
	ShoulderCm     *float64 `json:"shoulder_cm"`
	ChestCm        *float64 `json:"chest_cm"`
	LeftBicepCm    *float64 `json:"left_bicep_cm"`
	RightBicepCm   *float64 `json:"right_bicep_cm"`
	LeftForearmCm  *float64 `json:"left_forearm_cm"`
	RightForearmCm *float64 `json:"right_forearm_cm"`
	AbdomenCm      *float64 `json:"abdomen"`
	WaistCm        *float64 `json:"waist"`
	HipsCm         *float64 `json:"hips"`
	LeftThighCm    *float64 `json:"left_thigh"`
	RightThighCm   *float64 `json:"right_thigh"`
	LeftCalfCm     *float64 `json:"left_calf"`
	RightCalfCm    *float64 `json:"right_calf"`
}
