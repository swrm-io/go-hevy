<p align="center">
<a href="https://github.com/swrm-io/go-heyv/actions/workflows/github-code-scanning/codeql">
    <img alt="CodeQL Status" src="https://github.com/swrm-io/go-hevy/actions/workflows/github-code-scanning/codeql/badge.svg">
</a>
<a href="https://pkg.go.dev/github.com/swrm-io/go-hevy">
    <img alt="GoDoc" src="https://pkg.go.dev/badge/github.com/swrm-io/go-hevy.svg">
</a>

# go-hevy

> [!IMPORTANT]
> This is still a work in progress and may have bugs and is not feature complete.

Golang Client for working with the [Hevy](https://www.hevyapp.com/) API.

> [!NOTE]
> The API is only available to `Hevy Pro` users.

To generate your API key visit the [Developer Portal](https://hevy.com/settings?developer).

## Installation

```sh
go get github.com/swrm-io/go-hevy
```

## Usage

### Creating a Client

```go
import "github.com/swrm-io/go-hevy"

client := hevy.NewClient("your-api-key")
```

### User

```go
// Get authenticated user's profile
user, err := client.User()
if err != nil {
    // handle error
}
fmt.Println(user.Name)
```

### Workouts

```go
// Get all workouts (handles pagination automatically)
workouts, err := client.AllWorkouts()
if err != nil {
    // handle error
}

// Iterate over workouts one by one (memory efficient for large datasets)
for workout := range client.Workouts() {
    fmt.Println(workout.Title)
}

// Get a paginated list of workouts (max 10 per page)
workouts, nextPage, err := client.GetWorkouts(1, 10)
if err != nil {
    // handle error
}
// nextPage is 0 when there are no more pages
workouts, nextPage, err = client.GetWorkouts(nextPage, 10)

// Get a specific workout by ID
workoutID := uuid.MustParse("b459cba5-cd6d-463c-abd6-54f8eafcadcb")
workout, err := client.Workout(workoutID)
if err != nil {
    // handle error
}

// Get total workout count
count, err := client.WorkoutCount()
if err != nil {
    // handle error
}
```

### Workout Events

Workout events track changes (updates and deletes) since a given point in time, useful for syncing.

```go
since := time.Now().AddDate(0, -1, 0) // last month

// Get all workout events since a specific time
events, err := client.AllWorkoutEvents(since)
if err != nil {
    // handle error
}

// Iterate over workout events (memory efficient)
for event := range client.WorkoutEvents(since) {
    switch event.EventType {
    case hevy.EventTypeUpdated:
        // process updated workout
    case hevy.EventTypeDeleted:
        // process deleted workout
    }
}

// Get a paginated list of workout events
events, nextPage, err := client.GetWorkoutEvents(1, 10, since)
if err != nil {
    // handle error
}
```

### Routines

```go
// Get all routines (handles pagination automatically)
routines, err := client.AllRoutines()
if err != nil {
    // handle error
}

// Iterate over routines one by one
for routine := range client.Routines() {
    fmt.Println(routine.Title)
}

// Get a paginated list of routines (max 10 per page)
routines, nextPage, err := client.GetRoutines(1, 10)
if err != nil {
    // handle error
}

// Get a specific routine by ID
routineID := uuid.MustParse("routine-uuid")
routine, err := client.Routine(routineID)
if err != nil {
    // handle error
}
```

### Routine Folders

```go
// Get all routine folders (handles pagination automatically)
folders, err := client.AllRoutineFolders()
if err != nil {
    // handle error
}

// Iterate over routine folders one by one
for folder := range client.RoutineFolders() {
    fmt.Println(folder.Title)
}

// Get a paginated list of routine folders (max 10 per page)
folders, nextPage, err := client.GetRoutineFolders(1, 10)
if err != nil {
    // handle error
}

// Get a specific routine folder by ID
folder, err := client.RoutineFolder(42)
if err != nil {
    // handle error
}
```

### Body Measurements

```go
// Get all body measurements (handles pagination automatically)
measurements, err := client.AllBodyMeasurements()
if err != nil {
    // handle error
}

// Iterate over body measurements one by one
for m := range client.BodyMeasurements() {
    fmt.Printf("Weight: %.2f kg (%.2f lbs)\n", m.WeightKG, m.WeightLB)
}

// Get a paginated list of body measurements (max 10 per page)
measurements, nextPage, err := client.GetBodyMeasurements(1, 10)
if err != nil {
    // handle error
}

// Get a body measurement for a specific date
date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
measurement, err := client.BodyMeasurement(date)
if err != nil {
    // handle error
}
```
