package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"
)

// Almacenamiento temporal de tokens en memoria (los tokens pueden estar en memoria)
var tokens = make(map[string]*models.Usuario)

// LoginPolicialHandler maneja el login de policiales
func LoginPolicialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validar credenciales desde MongoDB
	usuario, err := database.ObtenerUsuarioPorCredencial(req.Credencial)
	if err != nil || usuario.PIN != req.PIN {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Credenciales inválidas",
		})
		return
	}

	if !usuario.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario inactivo",
		})
		return
	}

	// Generar token simple (en producción usar JWT)
	token := generateToken()
	tokens[token] = usuario

	// Crear respuesta sin el PIN
	usuarioResp := *usuario
	usuarioResp.PIN = ""

	response := models.LoginResponse{
		Token:   token,
		Usuario: usuarioResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateToken genera un token simple (en producción usar JWT)
func generateToken() string {
	return time.Now().Format("20060102150405") + "-token"
}

// GetUsuarioFromToken obtiene el usuario desde el token
func GetUsuarioFromToken(token string) (*models.Usuario, bool) {
	usuario, exists := tokens[token]
	if !exists {
		return nil, false
	}
	// Actualizar usuario desde la base de datos para asegurar datos actualizados
	usuarioActualizado, err := database.ObtenerUsuarioPorID(usuario.ID)
	if err != nil {
		return usuario, true // Retornar el usuario del token si falla la actualización
	}
	return usuarioActualizado, true
}

