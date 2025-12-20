package main

import (
	"gepn/middleware"
	"gepn/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	// Obtener el puerto de la variable de entorno, usar 8080 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configurar las rutas
	mux := routes.SetupRoutes()

	// Aplicar middlewares
	handler := middleware.CORSMiddleware(middleware.LoggingMiddleware(mux))

	// Iniciar el servidor
	addr := ":" + port
	log.Printf("🚀 Servidor GEPN iniciado en el puerto %s", port)
	log.Printf("📍 Health check disponible en: http://localhost:%s/health", port)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

