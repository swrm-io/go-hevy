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
// Get all workouts
workouts, err := client.Workouts()
if err != nil {
    // handle error
}

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

// Get workout events since a specific time (for syncing)
since := time.Now().AddDate(0, -1, 0) // Last month
events, err := client.WorkoutEvents(since)
if err != nil {
    // handle error
}
```

### Fetching Routines

```go
// Get all routines
routines, err := client.Routines()
if err != nil {
    // handle error
}
```