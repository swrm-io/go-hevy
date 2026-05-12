/*
Package hevy provides a Go client for the Hevy workout tracking API.

A Hevy Pro subscription is required. Find your API key in the Hevy app under
Settings → API.

# Creating a Client

	client := hevy.New("your-api-key")

Options can be passed to override defaults:

	client := hevy.New("your-api-key",
		hevy.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
		hevy.WithBaseURL("https://api.hevyapp.com"),
	)

# Accessing resources

The client exposes a service field per resource group:

	client.User
	client.Workouts
	client.Routines
	client.RoutineFolders
	client.ExerciseTemplates
	client.ExerciseHistory
	client.BodyMeasurements

# Pagination

List methods accept explicit page and pageSize parameters. Page is 1-based.
Maximum page sizes vary by resource: 10 for most, 100 for exercise templates.

	page, err := client.Workouts.List(ctx, 1, 10)
	fmt.Printf("page %d of %d, %d workouts\n", page.Page, page.PageCount, len(page.Workouts))

To fetch all pages automatically use the ListAll variant:

	workouts, err := client.Workouts.ListAll(ctx)

When you request a page beyond the last, List returns ErrNoMorePages. Passing
a pageSize that exceeds the endpoint maximum returns ErrInvalidPageSize before
making any network request.

# Error handling

Sentinel errors can be checked with errors.Is:

	_, err := client.Workouts.Get(ctx, id)
	if errors.Is(err, hevy.ErrNotFound) {
		// handle missing workout
	}

Available sentinels:

	hevy.ErrNotFound              // 404
	hevy.ErrUnauthorized          // 401 or 403
	hevy.ErrConflict              // 409 (e.g. duplicate body measurement date)
	hevy.ErrBadRequest            // 400
	hevy.ErrNoMorePages           // page exceeds total page count
	hevy.ErrInvalidPageSize       // pageSize exceeds endpoint maximum
	hevy.ErrRoutineLimitExceeded  // account routine limit reached
	hevy.ErrExerciseLimitExceeded // account custom exercise limit reached

For the raw HTTP status code and response body use errors.As:

	var apiErr *hevy.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.Body)
	}

# Examples

The following examples assume:

	ctx := context.Background()
	client := hevy.New("your-api-key")

## User info

	user, err := client.User.Info(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.Name, user.URL)

## Listing workouts

	page, err := client.Workouts.List(ctx, 1, 10)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range page.Workouts {
		fmt.Printf("%s  %s\n", w.StartTime.Format("2006-01-02"), w.Title)
	}

## Fetching all workouts across pages

	workouts, err := client.Workouts.ListAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("total:", len(workouts))

## Getting a single workout

	w, err := client.Workouts.Get(ctx, "workout-uuid")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d exercises\n", w.Title, len(w.Exercises))

## Creating a workout

	weight := 100.0
	reps := 5

	w, err := client.Workouts.Create(ctx, hevy.WorkoutInput{
		Title:     "Morning Strength",
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Exercises: []hevy.WorkoutExerciseInput{
			{
				ExerciseTemplateID: "exercise-template-id",
				Sets: []hevy.WorkoutSetInput{
					{Type: hevy.SetTypeWarmup, WeightKg: &weight, Reps: &reps},
					{Type: hevy.SetTypeNormal, WeightKg: &weight, Reps: &reps},
					{Type: hevy.SetTypeNormal, WeightKg: &weight, Reps: &reps},
				},
			},
		},
	})

## Polling for workout changes

	since := time.Now().Add(-24 * time.Hour)
	events, err := client.Workouts.EventsAll(ctx, &hevy.WorkoutEventsOptions{Since: &since})
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range events {
		switch e.Type {
		case hevy.WorkoutEventUpdated:
			fmt.Println("updated:", e.Workout.ID)
		case hevy.WorkoutEventDeleted:
			fmt.Println("deleted:", e.ID)
		}
	}

## Working with routines

	// List
	page, err := client.Routines.List(ctx, 1, 10)

	// Get
	r, err := client.Routines.Get(ctx, "routine-uuid")

	// Create
	r, err := client.Routines.Create(ctx, hevy.RoutineInput{
		Title: "Push Day",
		Notes: "Chest, shoulders, triceps",
		Exercises: []hevy.RoutineExerciseInput{
			{
				ExerciseTemplateID: "exercise-template-id",
				RestSeconds:        ptr(90),
				Sets: []hevy.RoutineSetInput{
					{Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
					{Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
					{Type: hevy.SetTypeNormal, Reps: ptr(8), WeightKg: ptr(80.0)},
				},
			},
		},
	})

## Exercise history

	start := time.Now().AddDate(0, -3, 0)
	entries, err := client.ExerciseHistory.Get(ctx, "exercise-template-id", &hevy.GetHistoryOptions{
		StartDate: &start,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		fmt.Printf("%s  %.1f kg × %d reps\n",
			e.WorkoutStartTime.Format("2006-01-02"), *e.WeightKg, *e.Reps)
	}

## Body measurements

	// Record
	err := client.BodyMeasurements.Create(ctx, hevy.BodyMeasurement{
		Date:     "2024-06-01",
		WeightKg: ptr(82.5),
	})

	// Get by date
	m, err := client.BodyMeasurements.Get(ctx, "2024-06-01")

	// Update (full replacement — nil fields are cleared to null on the server)
	err = client.BodyMeasurements.Update(ctx, "2024-06-01", hevy.BodyMeasurementUpdate{
		WeightKg: ptr(83.0),
	})
*/
package hevy
