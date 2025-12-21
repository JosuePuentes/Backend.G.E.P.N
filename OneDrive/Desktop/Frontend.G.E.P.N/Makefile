# Makefile para GEPN Backend

.PHONY: run build clean test docker-build docker-run help

# Variables
BINARY_NAME=app
MAIN_FILE=main.go

# Ejecutar la aplicación
run:
	@echo "🚀 Ejecutando servidor GEPN..."
	@go run $(MAIN_FILE)

# Compilar la aplicación
build:
	@echo "🔨 Compilando aplicación..."
	@go build -tags netgo -ldflags '-s -w' -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Compilación exitosa: $(BINARY_NAME)"

# Limpiar archivos compilados
clean:
	@echo "🧹 Limpiando archivos..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME).exe
	@echo "✅ Limpieza completada"

# Ejecutar tests
test:
	@echo "🧪 Ejecutando tests..."
	@go test -v ./...

# Formatear código
fmt:
	@echo "📝 Formateando código..."
	@go fmt ./...
	@echo "✅ Formateo completado"

# Verificar código
vet:
	@echo "🔍 Verificando código..."
	@go vet ./...
	@echo "✅ Verificación completada"

# Actualizar dependencias
deps:
	@echo "📦 Actualizando dependencias..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencias actualizadas"

# Construir imagen Docker
docker-build:
	@echo "🐳 Construyendo imagen Docker..."
	@docker build -t gepn-backend:latest .
	@echo "✅ Imagen construida: gepn-backend:latest"

# Ejecutar contenedor Docker
docker-run:
	@echo "🐳 Ejecutando contenedor Docker..."
	@docker run -p 8080:8080 --env PORT=8080 gepn-backend:latest

# Ayuda
help:
	@echo "📋 Comandos disponibles:"
	@echo "  make run          - Ejecutar la aplicación"
	@echo "  make build        - Compilar la aplicación"
	@echo "  make clean        - Limpiar archivos compilados"
	@echo "  make test         - Ejecutar tests"
	@echo "  make fmt          - Formatear código"
	@echo "  make vet          - Verificar código"
	@echo "  make deps         - Actualizar dependencias"
	@echo "  make docker-build - Construir imagen Docker"
	@echo "  make docker-run   - Ejecutar contenedor Docker"
	@echo "  make help         - Mostrar esta ayuda"


