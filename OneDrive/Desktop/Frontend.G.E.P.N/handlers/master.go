package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Módulos disponibles en el sistema
var ModulosDisponibles = []string{
	"rrhh",         // RRHH - Recursos Humanos
	"policial",     // Módulo Policial
	"denuncias",    // Denuncias
	"detenidos",    // Detenidos
	"minutas",      // Minutas Digitales
	"buscados",     // Más Buscados
	"verificacion", // Verificación de Cédulas
	"panico",       // Botón de Pánico
}

// JWT Secret Key (en producción usar variable de entorno)
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "gepn-secret-key-change-in-production" // Cambiar en producción
	}
	return secret
}

// Claims para JWT
type MasterClaims struct {
	UsuarioID primitive.ObjectID `json:"usuario_id"`
	Usuario   string            `json:"usuario"`
	Permisos  []string          `json:"permisos"`
	jwt.RegisteredClaims
}

// Almacenamiento temporal de tokens de master en memoria (para compatibilidad)
var masterTokens = make(map[string]*models.UsuarioMaster)

// LoginMasterHandler maneja el login de usuarios master con JWT
func LoginMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Error al decodificar login request: %v", err)
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Intento de login - Usuario: %s", req.Usuario)

	// Buscar master por usuario
	master, err := database.ObtenerUsuarioMasterPorUsuario(req.Usuario)
	if err != nil {
		log.Printf("❌ Error al buscar usuario master '%s': %v", req.Usuario, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	log.Printf("✅ Usuario encontrado: %s, Activo: %v", master.Usuario, master.Activo)

	// Verificar contraseña con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(master.Contraseña), []byte(req.Contraseña))
	if err != nil {
		log.Printf("❌ Error al verificar contraseña para usuario '%s': %v", req.Usuario, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	log.Printf("✅ Contraseña correcta para usuario: %s", req.Usuario)

	// Verificar que esté activo
	if !master.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario inactivo",
		})
		return
	}

	// Actualizar último acceso
	database.ActualizarUltimoAccesoMaster(master.ID)

	// Generar JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &MasterClaims{
		UsuarioID: master.ID,
		Usuario:   master.Usuario,
		Permisos:  master.Permisos,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   master.Usuario,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error al generar token", http.StatusInternalServerError)
		return
	}

	// Guardar en memoria para compatibilidad
	masterTokens[tokenString] = master

	// No retornar la contraseña
	masterResp := *master
	masterResp.Contraseña = ""

	response := models.LoginMasterResponse{
		Token:   tokenString,
		Master:  masterResp,
		Mensaje: "Login exitoso",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CrearUsuarioMasterHandler crea un nuevo usuario master (requiere autenticación)
func CrearUsuarioMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el usuario esté autenticado
	master, ok := GetMasterFromRequest(r)
	if !ok || master == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	var req models.CrearUsuarioMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Usuario == "" {
		http.Error(w, "El usuario es obligatorio", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "El email es obligatorio", http.StatusBadRequest)
		return
	}

	if req.Contraseña == "" || len(req.Contraseña) < 6 {
		http.Error(w, "La contraseña debe tener al menos 6 caracteres", http.StatusBadRequest)
		return
	}

	// Validar permisos
	if len(req.Permisos) == 0 {
		http.Error(w, "Debe asignar al menos un permiso", http.StatusBadRequest)
		return
	}

	// Verificar que los permisos sean válidos
	for _, permiso := range req.Permisos {
		valido := false
		for _, modulo := range ModulosDisponibles {
			if permiso == modulo {
				valido = true
				break
			}
		}
		if !valido {
			http.Error(w, "Permiso inválido: "+permiso, http.StatusBadRequest)
			return
		}
	}

	// Verificar que el usuario no exista
	_, err := database.ObtenerUsuarioMasterPorUsuario(req.Usuario)
	if err == nil {
		http.Error(w, "El usuario ya está registrado", http.StatusConflict)
		return
	}

	// Verificar que el email no exista
	_, err = database.ObtenerUsuarioMasterPorEmail(req.Email)
	if err == nil {
		http.Error(w, "El email ya está registrado", http.StatusConflict)
		return
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contraseña), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al procesar la contraseña", http.StatusInternalServerError)
		return
	}

	// Crear usuario master
	newMaster := models.UsuarioMaster{
		Usuario:       req.Usuario,
		Nombre:        req.Nombre,
		Email:         req.Email,
		Contraseña:    string(hashedPassword),
		Permisos:      req.Permisos,
		Activo:        true,
		CreadoPor:     master.Usuario,
		FechaCreacion: time.Now(),
	}

	if err := database.CrearUsuarioMaster(&newMaster); err != nil {
		http.Error(w, "Error al crear el usuario master: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// No retornar la contraseña
	newMaster.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Usuario master creado correctamente",
		"master":  newMaster,
	})
}

// ListarUsuariosMasterHandler lista todos los usuarios master
func ListarUsuariosMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el usuario esté autenticado
	_, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	masters, err := database.ListarUsuariosMaster()
	if err != nil {
		http.Error(w, "Error al listar usuarios master", http.StatusInternalServerError)
		return
	}

	// Ocultar contraseñas
	for i := range masters {
		masters[i].Contraseña = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masters)
}

// ActualizarPermisosHandler actualiza los permisos de un usuario master
func ActualizarPermisosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el usuario esté autenticado
	_, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	// Obtener ID de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "ID de usuario requerido", http.StatusBadRequest)
		return
	}

	usuarioID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req models.ActualizarPermisosRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validar permisos
	for _, permiso := range req.Permisos {
		valido := false
		for _, modulo := range ModulosDisponibles {
			if permiso == modulo {
				valido = true
				break
			}
		}
		if !valido {
			http.Error(w, "Permiso inválido: "+permiso, http.StatusBadRequest)
			return
		}
	}

	if err := database.ActualizarPermisosMaster(usuarioID, req.Permisos); err != nil {
		http.Error(w, "Error al actualizar permisos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Permisos actualizados correctamente",
	})
}

// ActivarUsuarioMasterHandler activa/desactiva un usuario master
func ActivarUsuarioMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el usuario esté autenticado
	_, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	// Obtener ID de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "ID de usuario requerido", http.StatusBadRequest)
		return
	}

	usuarioID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		Activo bool `json:"activo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	if err := database.ActualizarEstadoMaster(usuarioID, req.Activo); err != nil {
		http.Error(w, "Error al actualizar estado", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Estado actualizado correctamente",
	})
}

// ListarModulosHandler retorna la lista de módulos disponibles
func ListarModulosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"modulos": ModulosDisponibles,
	})
}

// VerificarMasterHandler verifica si el token es válido y retorna información del master
func VerificarMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	master, ok := GetMasterFromRequest(r)
	if !ok {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// No retornar la contraseña
	masterResp := *master
	masterResp.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masterResp)
}

// GetMasterFromRequest obtiene el master desde el token JWT en la request
func GetMasterFromRequest(r *http.Request) (*models.UsuarioMaster, bool) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, false
	}

	// Verificar JWT
	claims := &MasterClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !tkn.Valid {
		// Intentar con el sistema antiguo de tokens en memoria
		master, exists := masterTokens[token]
		if exists {
			return master, true
		}
		return nil, false
	}

	// Obtener master desde la base de datos
	master, err := database.ObtenerUsuarioMasterPorID(claims.UsuarioID)
	if err != nil {
		return nil, false
	}

	return master, true
}

// GetMasterFromToken obtiene el master desde el token (compatibilidad)
func GetMasterFromToken(token string) (*models.UsuarioMaster, bool) {
	// Intentar con JWT primero
	claims := &MasterClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err == nil && tkn.Valid {
		master, err := database.ObtenerUsuarioMasterPorID(claims.UsuarioID)
		if err == nil {
			return master, true
		}
	}

	// Fallback a tokens en memoria
	master, exists := masterTokens[token]
	if !exists {
		return nil, false
	}
	masterActualizado, err := database.ObtenerUsuarioMasterPorID(master.ID)
	if err != nil {
		return master, true
	}
	return masterActualizado, true
}

// InicializarUsuarioAdmin crea el usuario admin inicial si no existe
func InicializarUsuarioAdmin() error {
	// Verificar si ya existe el usuario admin (sin filtrar por activo)
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("usuarios_master")
	var master models.UsuarioMaster
	err := collection.FindOne(ctx, bson.M{"usuario": "admin"}).Decode(&master)
	if err == nil {
		log.Println("ℹ️  Usuario admin ya existe")
		// Asegurar que esté activo
		if !master.Activo {
			log.Println("⚠️  Usuario admin está inactivo, activándolo...")
			master.Activo = true
			if err := database.ActualizarEstadoMaster(master.ID, true); err != nil {
				log.Printf("⚠️  Error al activar usuario admin: %v", err)
			}
		}
		return nil
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Error al hashear contraseña: %v", err)
		return err
	}

	// Crear usuario admin con todos los permisos
	admin := models.UsuarioMaster{
		Usuario:       "admin",
		Nombre:        "Administrador",
		Email:         "admin@gepn.gob.ve",
		Contraseña:    string(hashedPassword),
		Permisos:      ModulosDisponibles, // Todos los permisos
		Activo:        true,
		CreadoPor:     "sistema",
		FechaCreacion: time.Now(),
	}

	if err := database.CrearUsuarioMaster(&admin); err != nil {
		log.Printf("❌ Error al crear usuario admin: %v", err)
		return err
	}

	log.Println("✅ Usuario admin creado automáticamente")
	log.Println("   Usuario: admin")
	log.Println("   Contraseña: Admin123! (CAMBIAR EN PRODUCCIÓN)")
	return nil
}

// InicializarAdminHandler endpoint temporal para inicializar el admin manualmente
func InicializarAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	err := InicializarUsuarioAdmin()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error al inicializar admin: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Usuario admin inicializado correctamente",
		"usuario": "admin",
		"contraseña": "Admin123!",
	})
}

// ResetearPasswordAdminHandler resetea la contraseña del admin a Admin123!
func ResetearPasswordAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Buscar usuario admin
	master, err := database.ObtenerUsuarioMasterPorUsuario("admin")
	if err != nil {
		// Si no existe, crearlo
		err = InicializarUsuarioAdmin()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Error al crear/resetear admin: " + err.Error(),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"mensaje": "Usuario admin creado",
			"usuario": "admin",
			"contraseña": "Admin123!",
		})
		return
	}

	// Hashear nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error al hashear contraseña: " + err.Error(),
		})
		return
	}

	// Actualizar contraseña
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("usuarios_master")
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": master.ID},
		bson.M{"$set": bson.M{
			"contraseña": string(hashedPassword),
			"activo":     true,
		}},
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error al actualizar contraseña: " + err.Error(),
		})
		return
	}

	log.Printf("✅ Contraseña del admin reseteada correctamente")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Contraseña del admin reseteada correctamente",
		"usuario": "admin",
		"contraseña": "Admin123!",
	})
}
