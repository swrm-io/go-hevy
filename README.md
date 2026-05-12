# go-hevy

A Go client library for the [Hevy](https://www.hevyapp.com) workout tracking API.

> **Note:** The Hevy API requires a Hevy Pro subscription. You can find your API key in the Hevy app under Settings → API.

## Installation

```bash
go get github.com/swrm-io/go-hevy
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/swrm-io/go-hevy"
)

func main() {
    client := hevy.New("your-api-key")
    ctx := context.Background()

    user, err := client.User.Info(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Logged in as:", user.Name)
}
```

## Examples

The examples below use a small generic helper to take the address of a literal value:

```go
func ptr[T any](v T) *T { return &v }
```

---

### Workouts

#### Get one page

```go
// pageSize max is 10
page, err := client.Workouts.List(ctx, 1, 10)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Page %d of %d — %d workouts\n", page.Page, page.PageCount, len(page.Workouts))
for _, w := range page.Workouts {
    fmt.Printf("  %s (%s)\n", w.Title, w.StartTime.Format("2006-01-02"))
}
```

#### Fetch all workouts across pages

```go
workouts, err := client.Workouts.ListAll(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Total fetched:", len(workouts))
```

#### Get total workout count

```go
count, err := client.Workouts.Count(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Total workouts:", count)
```

#### Get a single workout

```go
workout, err := client.Workouts.Get(ctx, "workout-uuid")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s: %d exercises\n", workout.Title, len(workout.Exercises))
```

#### Create a workout

```go
workout, err := client.Workouts.Create(ctx, hevy.WorkoutInput{
    Title:     "Morning Strength",
    StartTime: time.Now().Add(-1 * time.Hour),
    EndTime:   time.Now(),
    Exercises: []hevy.WorkoutExerciseInput{
        {
            ExerciseTemplateID: "exercise-template-uuid",
            Sets: []hevy.WorkoutSetInput{
                {Type: hevy.SetTypeWarmup, WeightKg: ptr(60.0), Reps: ptr(10)},
                {Type: hevy.SetTypeNormal, WeightKg: ptr(100.0), Reps: ptr(5)},
                {Type: hevy.SetTypeNormal, WeightKg: ptr(100.0), Reps: ptr(5)},
            },
        },
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Created workout:", workout.ID)
```

#### Update a workout

```go
notes := "felt strong today"
updated, err := client.Workouts.Update(ctx, "workout-uuid", hevy.WorkoutInput{
    Title:     "Morning Strength (updated)",
    StartTime: time.Now().Add(-1 * time.Hour),
    EndTime:   time.Now(),
    Exercises: []hevy.WorkoutExerciseInput{
        {
            ExerciseTemplateID: "exercise-template-uuid",
            Notes:              &notes,
            Sets: []hevy.WorkoutSetInput{
                {Type: hevy.SetTypeNormal, WeightKg: ptr(105.0), Reps: ptr(5)},
            },
        },
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Updated:", updated.Title)
```

#### Poll for workout changes since a date

```go
since := time.Now().Add(-24 * time.Hour)

// One page
page, err := client.Workouts.Events(ctx, 1, 10, &hevy.WorkoutEventsOptions{Since: &since})
if err != nil {
    log.Fatal(err)
}
for _, e := range page.Events {
    switch e.Type {
    case hevy.WorkoutEventUpdated:
        fmt.Println("Updated workout:", e.Workout.ID)
    case hevy.WorkoutEventDeleted:
        fmt.Println("Deleted workout:", e.ID)
    }
}

// Or fetch all pages at once
allEvents, err := client.Workouts.EventsAll(ctx, &hevy.WorkoutEventsOptions{Since: &since})
if err != nil {
    log.Fatal(err)
}
```

---

### Routines

```go
// One page (pageSize max is 10)
page, err := client.Routines.List(ctx, 1, 10)

// All pages
routines, err := client.Routines.ListAll(ctx)

// Single routine
routine, err := client.Routines.Get(ctx, "routine-uuid")

// Create
routine, err = client.Routines.Create(ctx, hevy.RoutineInput{
    Title: "Push Day",
    Notes: "Chest, shoulders, triceps",
    Exercises: []hevy.RoutineExerciseInput{
        {
            ExerciseTemplateID: "exercise-template-uuid",
            RestSeconds:        ptr(90),
            Sets: []hevy.RoutineSetInput{
                {Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
                {Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
                {Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
            },
        },
    },
})
```

---

### Routine Folders

```go
// One page (pageSize max is 10)
page, err := client.RoutineFolders.List(ctx, 1, 10)

// All pages
folders, err := client.RoutineFolders.ListAll(ctx)

// Single folder
folder, err := client.RoutineFolders.Get(ctx, 42)

// Create
folder, err = client.RoutineFolders.Create(ctx, "Strength")
```

---

### Exercise Templates

```go
// One page (pageSize max is 100)
page, err := client.ExerciseTemplates.List(ctx, 1, 100)

// All pages
templates, err := client.ExerciseTemplates.ListAll(ctx)

// Single template
tmpl, err := client.ExerciseTemplates.Get(ctx, "exercise-template-uuid")

// Create a custom exercise
id, err := client.ExerciseTemplates.Create(ctx, hevy.CreateExerciseInput{
    Title:             "Viking Press",
    ExerciseType:      hevy.ExerciseTypeWeightReps,
    EquipmentCategory: hevy.EquipmentMachine,
    MuscleGroup:       hevy.MuscleGroupShoulders,
    OtherMuscles:      []hevy.MuscleGroup{hevy.MuscleGroupTriceps},
})
```

---

### Exercise History

```go
// All history for an exercise
entries, err := client.ExerciseHistory.Get(ctx, "exercise-template-uuid", nil)
if err != nil {
    log.Fatal(err)
}

// Filter by date range
start := time.Now().AddDate(0, -3, 0)
entries, err = client.ExerciseHistory.Get(ctx, "exercise-template-uuid", &hevy.GetHistoryOptions{
    StartDate: &start,
})
if err != nil {
    log.Fatal(err)
}
for _, e := range entries {
    fmt.Printf("%s — %.1f kg × %d reps\n",
        e.WorkoutStartTime.Format("2006-01-02"),
        *e.WeightKg,
        *e.Reps,
    )
}
```

---

### Body Measurements

```go
// One page (pageSize max is 10)
page, err := client.BodyMeasurements.List(ctx, 1, 10)

// All pages
measurements, err := client.BodyMeasurements.ListAll(ctx)

// Single date (YYYY-MM-DD)
m, err := client.BodyMeasurements.Get(ctx, "2024-06-01")

// Create
err = client.BodyMeasurements.Create(ctx, hevy.BodyMeasurement{
    Date:       "2024-06-01",
    WeightKg:   ptr(82.5),
    FatPercent: ptr(15.2),
})

// Update (full replacement — nil fields are cleared to null on the server)
err = client.BodyMeasurements.Update(ctx, "2024-06-01", hevy.BodyMeasurementUpdate{
    WeightKg:   ptr(83.0),
    FatPercent: ptr(14.9),
})

// Convert to imperial units (pounds and inches)
imp := m.Imperial()
fmt.Printf("%.1f lbs, %.1f%% body fat\n", *imp.WeightLbs, *imp.FatPercent)
fmt.Printf("waist: %.1f in\n", *imp.WaistIn)
```

---

### Error Handling

The client exposes sentinel errors that can be checked with `errors.Is`:

```go
_, err := client.Workouts.Get(ctx, "nonexistent-id")
if errors.Is(err, hevy.ErrNotFound) {
    fmt.Println("workout not found")
}
```

| Sentinel | Triggered by |
|---|---|
| `hevy.ErrNotFound` | 404 response |
| `hevy.ErrUnauthorized` | 401 or 403 response |
| `hevy.ErrConflict` | 409 response (e.g. duplicate body measurement date) |
| `hevy.ErrBadRequest` | 400 response |
| `hevy.ErrNoMorePages` | requested page exceeds total page count |
| `hevy.ErrInvalidPageSize` | pageSize exceeds the endpoint maximum |
| `hevy.ErrRoutineLimitExceeded` | 403 on `POST /v1/routines` |
| `hevy.ErrExerciseLimitExceeded` | 403 on `POST /v1/exercise_templates` |

For the raw HTTP status code and response body, use `errors.As`:

```go
var apiErr *hevy.APIError
if errors.As(err, &apiErr) {
    fmt.Println("status:", apiErr.StatusCode)
    fmt.Println("body:", apiErr.Body)
}
```

---

### Client Options

```go
// Custom HTTP client
client := hevy.New("your-api-key",
    hevy.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
)

// Override base URL (useful for testing against a mock server)
client = hevy.New("your-api-key",
    hevy.WithBaseURL("http://localhost:8080"),
)
```

---

## API Reference

Full API documentation: https://api.hevyapp.com/docs


