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
> The API is only avaliable to `Hevy Pro` users.

To generate your API key visit the [Developer Portal](https://hevy.com/settings?developer).

## Usage

### Creating a Client

```go
import "github.com/swrm-io/go-hevy"

client := hevy.NewClient("your-api-key")
```

### Fetching Workouts

```go
// Get all workouts (handles pagination automatically)
workouts, err := client.AllWorkouts()
if err != nil {
    // handle error
}

// Iterate over workouts one by one (memory efficient for large datasets)
for workout := range client.Workouts() {
    // process workout
}

// Get a paginated list of workouts
workouts, nextPage, err := client.GetWorkouts(1, 10)
if err != nil {
    // handle error
}
currentPage = nextPage
workouts, nextPage, err := client.Getworkouts(currentPage, 10)

// Get a specific workout by ID
workoutID := uuid.MustParse("workout-uuid")
workout, err := client.Workout(workoutID)
if err != nil {
    // handle error
}

// Get workout count
count, err := client.WorkoutCount()
if err != nil {
    // handle error
}

// Get all workout events since a specific time (for syncing)
since := time.Now().AddDate(0, -1, 0) // Last month
events, err := client.AllWorkoutEvents(since)
if err != nil {
    // handle error
}

// Iterate over workout events (memory efficient)
since := time.Now().AddDate(0, -1, 0) // Last month
for event := range client.WorkoutEvents(since) {
    // process event
}

// Get a paginated list of workout events
since := time.Now().AddDate(0, -1, 0) // Last month
events, nextPage, err := client.GetWorkoutEvents(1, 10, since)
if err != nil {
    // handle error
}
currentPage = nextPage
events, nextPage, err := client.GetWorkoutEvents(currentPage, 10, since)
```

### Fetching Routines

```go
// Get all routines
routines, err := client.Routines()
if err != nil {
    // handle error
}
```