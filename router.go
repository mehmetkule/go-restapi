package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// AddRoutes for api creates routes
func (app *App) AddRoutes() {

	//Health Check Status
	app.AddRoute("GET", "/health", app.HealthCheck)

	// User
	app.AddRoute("POST", "/login", app.Login)
	app.AddRoute("POST", "/register", app.Register)
	app.AddRouteWithMiddleware("GET", "/rap/users", app.GetUsers,app.JWTHandler)
	app.AddRouteWithMiddleware("GET", "/rap/{id}", app.FindUserByID, app.JWTHandler)
	app.AddRouteWithMiddleware("GET", "/rap/email/{email}", app.FindUserByEmail, app.JWTHandler)
	app.AddRouteWithMiddleware("DELETE", "/rap/user/{id}", app.DeleteUser, app.JWTHandler)
	//Relation API

	//Health Check Status
	app.AddRoute("GET", "/health", app.HealthCheck)

	//File Upload API
	app.AddRoute("POST", "/rap/file/{parent_id}", app.AddFile)
	app.AddRouteWithMiddleware("GET", "/rap/file/{id}", app.FindFile, app.JWTHandler)
	app.AddRouteWithMiddleware("GET", "/rap/files/{parent_id}", app.FindFiles, app.JWTHandler)
	app.AddRouteWithMiddleware("DELETE", "/rap/file/{id}", app.DeleteFile, app.JWTHandler)
	app.AddRouteWithMiddleware("DELETE", "/rap/files/{parent_id}", app.DeleteFiles, app.JWTHandler)
}

//HealthCheck checks application status
func (app *App) HealthCheck(writer http.ResponseWriter, request *http.Request) {
	health := map[string]string{}
	err := app.db.Ping()
	if err != nil {
		health["database"] = "down"
		health["error"] = err.Error()
		health["status"] = "down"
	} else {
		health["database"] = "up"
		health["status"] = "up"
	}
	app.RenderJSON(writer, http.StatusOK, health)
}

// JWTHandler jwt handler
func (app *App) JWTHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		const bearerSchema = "Bearer "
		authorization := request.Header.Get("Authorization")
		if authorization == "" {
			app.RenderErrorResponse(response, http.StatusForbidden, nil, "No token given")
			return
		}
		if strings.Contains(authorization, bearerSchema) {
			var token string
			if token = authorization[len(bearerSchema):]; len(token) == 0 {
				app.RenderErrorResponse(response, http.StatusForbidden, nil, "Invalid authorization token")
				return
			}
			result, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.conf.PasswordKey), nil
			})
			if err == nil && result.Valid {
				next.ServeHTTP(response, request)
				return
			} else {
				app.RenderErrorResponse(response, http.StatusForbidden, nil, "Invalid login token or expired")
				return
			}
		}
		app.RenderErrorResponse(response, http.StatusForbidden, nil, "Bad token")
		return
	})
}

// Logging traces all API endpoints
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		logger.Logger().Debug("API transcation",zap.String("method",req.Method),zap.Time("start",start))
	})
}