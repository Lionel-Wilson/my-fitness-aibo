// Command api is the my-fitness-aibo HTTP backend.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiauth "github.com/lionel/my-fitness-aibo/backend/internal/api/auth"
	apicycle "github.com/lionel/my-fitness-aibo/backend/internal/api/cycle"
	apiexercise "github.com/lionel/my-fitness-aibo/backend/internal/api/exercise"
	apiexerciselog "github.com/lionel/my-fitness-aibo/backend/internal/api/exerciselog"
	apiplan "github.com/lionel/my-fitness-aibo/backend/internal/api/plan"
	apiprogress "github.com/lionel/my-fitness-aibo/backend/internal/api/progress"
	apiworkout "github.com/lionel/my-fitness-aibo/backend/internal/api/workout"
	internalauth "github.com/lionel/my-fitness-aibo/backend/internal/auth"
	"github.com/lionel/my-fitness-aibo/backend/internal/config"
	"github.com/lionel/my-fitness-aibo/backend/internal/cycle"
	cyclestorage "github.com/lionel/my-fitness-aibo/backend/internal/cycle/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/db"
	"github.com/lionel/my-fitness-aibo/backend/internal/exercise"
	exercisestorage "github.com/lionel/my-fitness-aibo/backend/internal/exercise/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/exerciselog"
	exerciselogstorage "github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/http/router"
	"github.com/lionel/my-fitness-aibo/backend/internal/plan"
	planstorage "github.com/lionel/my-fitness-aibo/backend/internal/plan/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/progress"
	progressstorage "github.com/lionel/my-fitness-aibo/backend/internal/progress/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/user"
	userstorage "github.com/lionel/my-fitness-aibo/backend/internal/user/storage"
	"github.com/lionel/my-fitness-aibo/backend/internal/workout"
	workoutstorage "github.com/lionel/my-fitness-aibo/backend/internal/workout/storage"
	commonauth "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	logger := log.Default()

	log.Println("running migrations...")
	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	tokenManager := commonauth.NewTokenManager(cfg.JWTSecret, time.Duration(cfg.JWTTTLHours)*time.Hour)

	userRepo := userstorage.NewUserRepository(pool)
	planRepo := planstorage.NewPlanRepository(pool)
	workoutRepo := workoutstorage.NewWorkoutRepository(pool)
	exerciseRepo := exercisestorage.NewExerciseRepository(pool)
	cycleRepo := cyclestorage.NewCycleRepository(pool)
	exerciseLogRepo := exerciselogstorage.NewExerciseLogRepository(pool)
	progressRepo := progressstorage.NewProgressRepository(pool)

	userService := user.NewService(logger, userRepo)
	authService := internalauth.NewService(logger, userService, tokenManager)
	planService := plan.NewService(logger, planRepo)
	workoutService := workout.NewService(logger, workoutRepo)
	exerciseService := exercise.NewService(logger, exerciseRepo)
	cycleService := cycle.NewService(logger, cycleRepo)
	exerciseLogService := exerciselog.NewService(logger, exerciseLogRepo)
	progressService := progress.NewService(logger, progressRepo)

	handlers := router.Handlers{
		Auth:        apiauth.NewHandler(logger, authService),
		Plan:        apiplan.NewHandler(logger, planService),
		Workout:     apiworkout.NewHandler(logger, workoutService),
		Exercise:    apiexercise.NewHandler(logger, exerciseService),
		Cycle:       apicycle.NewHandler(logger, cycleService),
		ExerciseLog: apiexerciselog.NewHandler(logger, exerciseService, exerciseLogService),
		Progress:    apiprogress.NewHandler(logger, exerciseService, progressService),
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(cfg.CORSOrigins, tokenManager, handlers),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
