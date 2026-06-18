package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	apiauth "github.com/lionel/my-fitness-aibo/backend/internal/api/auth"
	apicycle "github.com/lionel/my-fitness-aibo/backend/internal/api/cycle"
	apiexercise "github.com/lionel/my-fitness-aibo/backend/internal/api/exercise"
	apiexerciselog "github.com/lionel/my-fitness-aibo/backend/internal/api/exerciselog"
	apiplan "github.com/lionel/my-fitness-aibo/backend/internal/api/plan"
	apiprogress "github.com/lionel/my-fitness-aibo/backend/internal/api/progress"
	apiworkout "github.com/lionel/my-fitness-aibo/backend/internal/api/workout"
	"github.com/lionel/my-fitness-aibo/backend/internal/middleware"
	commonauth "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/auth"
	"github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/render"
)

type Handlers struct {
	Auth         apiauth.Handler
	Plan         apiplan.Handler
	Workout      apiworkout.Handler
	Exercise     apiexercise.Handler
	Cycle        apicycle.Handler
	ExerciseLog  apiexerciselog.Handler
	Progress     apiprogress.Handler
}

func New(corsOrigins []string, tokenManager *commonauth.TokenManager, handlers Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		render.Json(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/signup", handlers.Auth.Signup())
		r.Post("/auth/login", handlers.Auth.Login())
		r.Post("/auth/logout", handlers.Auth.Logout())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator(tokenManager))

			r.Get("/auth/me", handlers.Auth.Me())

			r.Route("/plans", func(r chi.Router) {
				r.Get("/", handlers.Plan.ListPlans())
				r.Post("/", handlers.Plan.CreatePlan())
				r.Route("/{planID}", func(r chi.Router) {
					r.Get("/", handlers.Plan.GetPlan())
					r.Patch("/", handlers.Plan.UpdatePlan())
					r.Delete("/", handlers.Plan.DeletePlan())
					r.Get("/workouts", handlers.Workout.ListWorkouts())
					r.Post("/workouts", handlers.Workout.CreateWorkout())
					r.Get("/cycles", handlers.Cycle.ListCycles())
					r.Post("/cycles", handlers.Cycle.CreateCycle())
				})
			})

			r.Route("/workouts/{workoutID}", func(r chi.Router) {
				r.Get("/", handlers.Workout.GetWorkout())
				r.Patch("/", handlers.Workout.UpdateWorkout())
				r.Delete("/", handlers.Workout.DeleteWorkout())
				r.Get("/exercises", handlers.Exercise.ListExercises())
				r.Post("/exercises", handlers.Exercise.CreateExercise())
			})

			r.Route("/exercises/{exerciseID}", func(r chi.Router) {
				r.Get("/", handlers.Exercise.GetExercise())
				r.Patch("/", handlers.Exercise.UpdateExercise())
				r.Delete("/", handlers.Exercise.DeleteExercise())
				r.Get("/logs", handlers.ExerciseLog.ListExerciseLogs())
				r.Put("/logs/{cycleID}", handlers.ExerciseLog.UpsertExerciseLog())
				r.Get("/progress", handlers.Progress.ExerciseProgress())
			})

			r.Route("/cycles/{cycleID}", func(r chi.Router) {
				r.Patch("/", handlers.Cycle.UpdateCycle())
				r.Delete("/", handlers.Cycle.DeleteCycle())
			})
		})
	})

	return r
}
