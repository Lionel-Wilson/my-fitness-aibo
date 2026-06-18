package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/lionel/my-fitness-aibo/backend/internal/middleware"
)

// Router builds the full HTTP router: health check, auth routes, and the
// authenticated API.
func (h *Handler) Router(corsOrigins []string) http.Handler {
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
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api", func(r chi.Router) {
		// Public auth routes.
		r.Post("/auth/signup", h.Signup)
		r.Post("/auth/login", h.Login)
		r.Post("/auth/logout", h.Logout)

		// Authenticated routes.
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticator(h.Tokens))

			r.Get("/auth/me", h.Me)

			r.Route("/plans", func(r chi.Router) {
				r.Get("/", h.ListPlans)
				r.Post("/", h.CreatePlan)
				r.Route("/{planID}", func(r chi.Router) {
					r.Get("/", h.GetPlan)
					r.Patch("/", h.UpdatePlan)
					r.Delete("/", h.DeletePlan)
					r.Get("/workouts", h.ListWorkouts)
					r.Post("/workouts", h.CreateWorkout)
					r.Get("/cycles", h.ListCycles)
					r.Post("/cycles", h.CreateCycle)
				})
			})

			r.Route("/workouts/{workoutID}", func(r chi.Router) {
				r.Get("/", h.GetWorkout)
				r.Patch("/", h.UpdateWorkout)
				r.Delete("/", h.DeleteWorkout)
				r.Get("/exercises", h.ListExercises)
				r.Post("/exercises", h.CreateExercise)
			})

			r.Route("/exercises/{exerciseID}", func(r chi.Router) {
				r.Get("/", h.GetExercise)
				r.Patch("/", h.UpdateExercise)
				r.Delete("/", h.DeleteExercise)
				r.Get("/logs", h.ListExerciseLogs)
				r.Put("/logs/{cycleID}", h.UpsertExerciseLog)
				r.Get("/progress", h.ExerciseProgress)
			})

			r.Route("/cycles/{cycleID}", func(r chi.Router) {
				r.Patch("/", h.UpdateCycle)
				r.Delete("/", h.DeleteCycle)
			})
		})
	})

	return r
}
