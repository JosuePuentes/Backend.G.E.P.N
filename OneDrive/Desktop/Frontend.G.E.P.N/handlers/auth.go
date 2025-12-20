package handlers

import (
	"encoding/json"
	"gepn/models"
	"net/http"
	"time"
)

// Almacenamiento temporal en memoria (en producción usar base de datos)
var usuarios = map[string]*models.Usuario{
	"POL001": {
		ID:         1,
		Credencial: "POL001",
		PIN:        "123456",
		Nombre:     "Juan Pérez",
		Rango:      "Oficial",
		Activo:     true,
	},
	"POL002": {
		ID:         2,
		Credencial: "POL002",
		PIN:        "654321",
		Nombre:     "María González",
		Rango:      "Sargento",
		Activo:     true,
	},
}

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

	// Validar credenciales
	usuario, exists := usuarios[req.Credencial]
	if !exists || usuario.PIN != req.PIN {
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
	return usuario, exists
}

