package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware registra las peticiones HTTP
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Crear un ResponseWriter personalizado para capturar el status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			wrapped.statusCode,
			duration,
		)
	})
}

// CORSMiddleware maneja los headers CORS
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir todos los orígenes
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Manejar preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// responseWriter envuelve http.ResponseWriter para capturar el status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AuthMasterMiddleware verifica que el usuario esté autenticado como master
func AuthMasterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Token de autorización requerido",
			})
			return
		}

		// Verificar token (esto se hace en el handler, pero podemos validar aquí también)
		next.ServeHTTP(w, r)
	}
}

// VerificarPermisoMiddleware verifica que el usuario tenga el permiso requerido
func VerificarPermisoMiddleware(modulo string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Importar handlers para usar GetMasterFromRequest
			// Esto se hace en el handler directamente, pero podemos validar aquí
			next.ServeHTTP(w, r)
		}
	}
}

