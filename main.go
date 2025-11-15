package main

import (
	"context"
	"errors"
	"fmt"
	"microservicetest/app/vehicle"
	"microservicetest/infra/azure"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"microservicetest/app/healthcheck"
	"microservicetest/infra/couchbase"
	"microservicetest/pkg/config"
	apperrors "microservicetest/pkg/errors"
	_ "microservicetest/pkg/log"
)

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)
		return c.Next()
	}
}

func RequestDurationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		requestID := c.Locals("requestID").(string)
		zap.L().Info("Request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", c.Response().StatusCode()),
			zap.Float64("duration_seconds", duration),
			zap.Int("response_size", len(c.Response().Body())),
		)

		return err
	}
}

type Request any
type Response any

// Define an interface for handlers
type HandlerInterface[R Request, Res Response] interface {
	Handle(ctx context.Context, req *R) (*Res, error)
}

type HandlerCtxInterface[R Request, Res Response] interface {
	Handle(ctx *fiber.Ctx, req *R) (*Res, error)
}

// Update handle function to accept HandlerInterface instead of Handler function
func handle[R Request, Res Response](handler HandlerInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ParamsParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		/*
			ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
			defer cancel()
		*/

		ctx := c.UserContext()

		res, err := handler.Handle(ctx, &req)
		if err != nil {
			return apperrors.HandleError(c, err)
		}

		return c.JSON(res)
	}
}

func handleFiberCtx[R Request, Res Response](handler HandlerCtxInterface[R, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := c.BodyParser(&req); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ParamsParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.ReqHeaderParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		/*
			ctx, cancel := context.WithTimeout(c.UserContext(), 3*time.Second)
			defer cancel()
		*/

		res, err := handler.Handle(c, &req)
		if err != nil {
			return apperrors.HandleError(c, err)
		}

		return c.JSON(res)
	}
}

// Handler interface for raw fiber context (no response struct)
type HandlerRawInterface[R Request] interface {
	Handle(ctx *fiber.Ctx, req *R) error
}

func handleRaw[R Request](handler HandlerRawInterface[R]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req R

		if err := c.ParamsParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := c.QueryParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return handler.Handle(c, &req)
	}
}

func main() {
	appConfig := config.Read()
	defer zap.L().Sync()
	zap.L().Info("app starting...")
	zap.L().Info("app config", zap.Any("appConfig", appConfig))

	storageService, err := azure.NewStorage(appConfig.AzureConnectionString, "documents")
	if err != nil {
		zap.L().Error("Failed to initialize Azure Blob service", zap.Error(err))
	}

	couchbaseRepository := couchbase.NewVehicleRepository(appConfig.CouchbaseUrl, appConfig.CouchbaseUsername, appConfig.CouchbasePassword)

	healthcheckHandler := healthcheck.NewHealthCheckHandler()

	// Vehicle handlers
	createVehicleHandler := vehicle.NewCreateVehicleHandler(couchbaseRepository)
	getVehicleHandler := vehicle.NewGetVehicleHandler(couchbaseRepository)
	updateVehicleHandler := vehicle.NewUpdateVehicleHandler(couchbaseRepository)
	addDocumentHandler := vehicle.NewAddDocumentHandler(couchbaseRepository, storageService)
	getDocumentHandler := vehicle.NewGetDocumentsHandler(couchbaseRepository)
	deleteDocumentHandler := vehicle.NewDeleteDocumentHandler(couchbaseRepository, storageService)
	downloadDocumentHandler := vehicle.NewDownloadDocumentHandler(couchbaseRepository, storageService)

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	app.Use(RequestIDMiddleware())
	app.Use(RequestDurationMiddleware())

	// Health check endpoint
	app.Get("/healthcheck", handle[healthcheck.HealthCheckRequest, healthcheck.HealthCheckResponse](healthcheckHandler))

	// Vehicle endpoints
	app.Post("/vehicles", handle[vehicle.CreateVehicleRequest, vehicle.CreateVehicleResponse](createVehicleHandler))
	app.Get("/vehicles/:id", handle[vehicle.GetVehicleRequest, vehicle.GetVehicleResponse](getVehicleHandler))
	app.Put("/vehicles/:id", handle[vehicle.UpdateVehicleRequest, vehicle.UpdateVehicleResponse](updateVehicleHandler))
	app.Post("/vehicles/:id/documents", handleFiberCtx[vehicle.AddDocumentRequest, vehicle.AddDocumentResponse](addDocumentHandler))
	app.Get("/vehicles/:id/documents", handleFiberCtx[vehicle.GetDocumentsRequest, vehicle.GetDocumentsResponse](getDocumentHandler))
	app.Get("/vehicles/:id/documents/:doc_id/download", handleRaw[vehicle.DownloadDocumentRequest](downloadDocumentHandler))
	app.Delete("/vehicles/:id/documents/:doc_id", handleFiberCtx[vehicle.DeleteDocumentRequest, vehicle.DeleteDocumentResponse](deleteDocumentHandler))

	// Start server in a goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", appConfig.Port)); err != nil {
			zap.L().Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	zap.L().Info("Server started on port", zap.String("port", appConfig.Port))

	gracefulShutdown(app)
}

func gracefulShutdown(app *fiber.App) {
	// Create channel for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	zap.L().Info("Shutting down server...")

	// Shutdown with 5 second timeout
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Error during server shutdown", zap.Error(err))
	}

	zap.L().Info("Server gracefully stopped")
}
