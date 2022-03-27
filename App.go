package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/mehmetkule/go-restapi/internal/dto"
	"github.com/mehmetkule/go-restapi/internal/store"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"log"

	"net/http"
	"os"
	"time"
)

// ServerConfig is config struct for server
type ServerConfig struct {
	Addr    string `yaml:"listen_addr" envconfig:"LISTEN_ADDR"  default:":8080"`
	Timeout struct {
		// graceful shutdown
		Server time.Duration `yaml:"server" default:90`
		// write operation
		Write time.Duration `yaml:"write" default:90`
		// read operation
		Read time.Duration `yaml:"read" default:90`
		// time until idle session is closed
		Idle time.Duration `yaml:"idle" default:90`
	} `yaml:"timeout"`
}

type EmailConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Target   string `yaml:"target"`
}

type Config struct {
	AppName      string
	PasswordKey  string         `yaml:"password_key" envconfig:"CERCI_PASSWORD_KEY"`
	ServerConfig ServerConfig   `yaml:"server"`
	DBConfig     DbConfig `yaml:"database"`
	EmailConfig  EmailConfig    `yaml:"smtp"`
}

type App struct {
	conf         *Config
	db           *sql.DB
	Router       *mux.Router
	ShutdownHook func()
	userRepo     *store.UserRepo
	filesRepo    *store.FilesRepo
}

// NewConfig creates a new config from yaml file
// It first reads from config.yaml file. It overrides values from environment values
func NewConfig(configPath string) (*Config, error) {
	logger.Logger().Info("reading from config path", zap.String("configPath", configPath))

	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	logger.Logger().Info("Reading environment variables")
	err = envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// NewApp creates a new instance of App
func NewApp(conf *Config) *App {
	return &App{
		conf:   conf,
		Router: mux.NewRouter(),
	}
}


func (app *App) Initialize() error {
	database, err := app.conf.DBConfig.GetDatabase()
	if err != nil {
		return err
	}
	if err != nil {
		logger.Logger().Error("failed to get postgres connection", zap.Error(err))
		return err
	}
	app.filesRepo = &store.FilesRepo{DB: database}
	app.userRepo = &store.UserRepo{DB: database}
	app.AddRoutes()

	app.ShutdownHook = func() {
		logger.Logger().Info("Closing database connections....")
		if app != nil && app.db != nil {
			app.db.Close()
		}
	}
	return nil
}

// Run runs the server
func (app *App) Run() error {
	var runChan = make(chan os.Signal, 1)
	handler := Logging(app.Router)
	server := &http.Server{
		Addr:         app.conf.ServerConfig.Addr,
		Handler:      handler,
		ReadTimeout:  app.conf.ServerConfig.Timeout.Read * time.Second,
		WriteTimeout: app.conf.ServerConfig.Timeout.Write * time.Second,
		IdleTimeout:  app.conf.ServerConfig.Timeout.Idle * time.Second,
	}
	logger.Logger().Info("Staring server :", zap.String("address", server.Addr))
	// Run the server on a new goroutine: we will wait on signals
	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				logger.Logger().Warn("Server closing under normal conditions")
			} else {
				logger.Logger().Fatal("Server failed to start due to error", zap.Error(err))
			}
		}
	}()

	 // Handle ctrl+c/ctrl+x interrupts
	//signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)
	 // Block on this channel and waith for signals
	interrupt := <-runChan

	// Call shutdown hook so that application can finalize itself
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	app.ShutdownHook()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShutDown); err != nil {
		logger.Logger().Fatal("Server was unable to gracefully shutdown due to err", zap.Error(err))
	}
	return nil
}

// AddRoute adds a route to applicatoin
func (app *App) AddRoute(method string, route string, apiHandler func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(route, apiHandler).Methods(method)
}

// AddRouteWithMiddleware wraps the route in middleware
func (app *App) AddRouteWithMiddleware(method string, route string, apiHandler func(w http.ResponseWriter, r *http.Request), middleware func(next http.Handler) http.Handler) {
	handler := middleware(http.HandlerFunc(apiHandler))
	app.Router.Handle(route, handler).Methods(method)
}

// RenderErrorResponse render error response
func (app *App) RenderErrorResponse(writer http.ResponseWriter, httpStatus int, err error, message string) {
	logger.Logger().Error("Error render response", zap.Error(err), zap.String("", message))
	writer.WriteHeader(httpStatus)
	app.RenderJSON(writer, httpStatus, dto.ErrorResponse{Status: httpStatus, Error: err, Message: message})
}
func (app *App) RenderJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Write(jsonData)
	}
}
