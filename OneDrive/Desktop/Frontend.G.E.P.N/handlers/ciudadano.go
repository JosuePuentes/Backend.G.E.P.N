package handlers

import (
	"encoding/json"
	"errors"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Almacenamiento temporal de tokens de ciudadanos
var ciudadanoTokens = make(map[string]*models.Ciudadano)

// RegistroCiudadanoHandler maneja el registro de ciudadanos
func RegistroCiudadanoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegistroCiudadanoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Nombre == "" || req.Cedula == "" || req.Telefono == "" || req.Contraseña == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Todos los campos son requeridos",
		})
		return
	}

	if len(req.Contraseña) < 6 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "La contraseña debe tener al menos 6 caracteres",
		})
		return
	}

	// Verificar si la cédula ya existe
	ciudadanoExistente, err := database.ObtenerCiudadanoPorCedula(req.Cedula)
	if err == nil && ciudadanoExistente != nil {
		// Ciudadano ya existe
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "La cédula ya está registrada",
		})
		return
	}
	// Si el error es "no documents", significa que no existe y podemos continuar
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		// Error real de base de datos
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al verificar cédula: " + err.Error(),
		})
		return
	}
	// Si llegamos aquí, la cédula no existe y podemos continuar

	// Crear ciudadano
	ciudadano := &models.Ciudadano{
		Nombre:        req.Nombre,
		Cedula:        req.Cedula,
		Telefono:      req.Telefono,
		Contraseña:    req.Contraseña, // En producción, hashear con bcrypt
		FechaRegistro: time.Now(),
		Activo:        true,
	}

	if err := database.CrearCiudadano(ciudadano); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al registrar usuario: " + err.Error(),
		})
		return
	}

	// Generar token simple (en producción usar JWT)
	token := generateCiudadanoToken()
	ciudadanoTokens[token] = ciudadano

	// Respuesta sin contraseña
	ciudadanoResp := *ciudadano
	ciudadanoResp.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Usuario registrado correctamente",
		"data": map[string]interface{}{
			"id":       ciudadanoResp.ID.Hex(),
			"nombre":   ciudadanoResp.Nombre,
			"cedula":   ciudadanoResp.Cedula,
			"telefono": ciudadanoResp.Telefono,
		},
		"token": token,
	})
}

// LoginCiudadanoHandler maneja el login de ciudadanos
func LoginCiudadanoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginCiudadanoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Cedula == "" || req.Contraseña == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Cédula y contraseña son requeridos",
		})
		return
	}

	// Buscar ciudadano
	ciudadano, err := database.ObtenerCiudadanoPorCedula(req.Cedula)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Cédula o contraseña incorrectos",
		})
		return
	}

	// Verificar contraseña (en producción usar bcrypt)
	if ciudadano.Contraseña != req.Contraseña {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Cédula o contraseña incorrectos",
		})
		return
	}

	if !ciudadano.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Usuario inactivo",
		})
		return
	}

	// Generar token
	token := generateCiudadanoToken()
	ciudadanoTokens[token] = ciudadano

	// Respuesta sin contraseña
	ciudadanoResp := *ciudadano
	ciudadanoResp.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login exitoso",
		"data": map[string]interface{}{
			"id":       ciudadanoResp.ID.Hex(),
			"nombre":   ciudadanoResp.Nombre,
			"cedula":   ciudadanoResp.Cedula,
			"telefono": ciudadanoResp.Telefono,
		},
		"token": token,
	})
}

// generateCiudadanoToken genera un token simple para ciudadanos
func generateCiudadanoToken() string {
	return time.Now().Format("20060102150405") + "-ciudadano-token"
}

// GetCiudadanoFromToken obtiene el ciudadano desde el token
func GetCiudadanoFromToken(token string) (*models.Ciudadano, bool) {
	ciudadano, exists := ciudadanoTokens[token]
	if !exists {
		return nil, false
	}
	// Actualizar desde la base de datos para asegurar datos actualizados
	ciudadanoActualizado, err := database.ObtenerCiudadanoPorID(ciudadano.ID)
	if err != nil {
		return ciudadano, true // Retornar el ciudadano del token si falla la actualización
	}
	return ciudadanoActualizado, true
}

