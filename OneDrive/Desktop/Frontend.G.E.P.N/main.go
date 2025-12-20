package main

import (
	"gepn/database"
	"gepn/middleware"
	"gepn/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Conectar a MongoDB
	log.Println("🔌 Conectando a MongoDB...")
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Error al conectar a MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Inicializar datos por defecto
	log.Println("📦 Inicializando datos por defecto...")
	if err := database.InicializarDatos(); err != nil {
		log.Printf("⚠️  Error al inicializar datos: %v", err)
	}

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

	// Manejar señales para cierre graceful
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(addr, handler); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-sigChan
	log.Println("🛑 Cerrando servidor...")
}

