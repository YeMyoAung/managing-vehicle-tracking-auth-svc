package app

import (
    "context"
    "errors"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/go-playground/validator/v10"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/config"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/handler"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/repositories"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/services"
    "github.com/yemyoaung/managing-vehicle-tracking-common"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var (
    ErrConfigMissing = errors.New("config is missing")
)

type App struct {
    validator *validator.Validate
    cfg       *config.EnvConfig
    db        *mongo.Client
    shutdown  chan error
    exit      chan os.Signal
}

// NewApp creates a new App instance
func NewApp() *App {
    exit := make(chan os.Signal, 1)
    shutdown := make(chan error, 1)

    signal.Notify(exit, os.Interrupt, syscall.SIGTERM) // listen for termination signals

    go func() {
        defer close(exit)
        <-exit
        shutdown <- nil // shutdown 
    }()

    return &App{shutdown: shutdown}
}

// SetValidator sets the validator for the app
func (a *App) SetValidator(validator *validator.Validate) *App {
    a.validator = validator
    return a
}

// SetConfig sets the configuration for the app
func (a *App) SetConfig(cfg *config.EnvConfig) *App {
    a.cfg = cfg
    return a
}

// Run starts the app, connects to MongoDB, and starts the HTTP server
func (a *App) Run(ctx context.Context) {
    var err error
    if a.cfg == nil {
        a.shutdown <- ErrConfigMissing
        return
    }

    // Connect to MongoDB
    a.db, err = mongo.Connect(ctx, options.Client().ApplyURI(a.cfg.DatabaseURL))
    if err != nil {
        a.shutdown <- err
        return
    }

    // Initialize repository and service
    adminRepo, err := repositories.NewMongoAuthRepository(ctx, a.db.Database("users"))
    if err != nil {
        a.shutdown <- err
        return
    }
    authService := services.NewMongoAuthService(adminRepo, common.NewJwtMaker(), a.cfg)
    authHandler := handler.NewV1AuthHandler(authService, a.validator)

    // Set up the HTTP server
    server := http.NewServeMux()

    // Set up the API routes
    v1Router := http.NewServeMux()
    v1Router.HandleFunc("/api/v1/login", authHandler.Login)
    v1Router.HandleFunc("/api/v1/me", authHandler.Me)

    // Apply middlewares and handle requests
    // The v1Router (which holds our API routes) will have three middlewares applied:
    // - CorsMiddleware: Adds CORS headers to the response
    // - LoggingMiddleware: Logs each incoming request for debugging and monitoring
    // - VerifySignatureMiddleware: Verifies the request's signature (ensuring it's from a trusted source)
    server.Handle(
        "/", common.CorsMiddleware(nil)(
            common.LoggingMiddleware(log.Default())(
                common.VerifySignatureMiddleware(a.cfg.SignatureKey)(v1Router),
            ),
        ),
    )

    log.Println("Auth service started on Port: ", a.cfg.Port)

    // Start the HTTP server
    go func() {
        err = http.ListenAndServe(a.cfg.Host+":"+a.cfg.Port, server)
        if !errors.Is(err, http.ErrServerClosed) {
            a.shutdown <- err
        }
    }()
}

// Shutdown gracefully shuts down the app
func (a *App) Shutdown(ctx context.Context) error {
    defer close(a.shutdown)

    // Disconnect from MongoDB
    defer func(ctx context.Context, client *mongo.Client) {
        if client == nil {
            return
        }
        err := client.Disconnect(ctx)
        if err != nil {
            log.Println("Failed to disconnect from database", err)
        }
    }(ctx, a.db)

    return <-a.shutdown
}
