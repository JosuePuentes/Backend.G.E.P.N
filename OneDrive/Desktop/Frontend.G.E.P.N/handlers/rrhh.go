package handlers

import (
	"encoding/base64"
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Rangos válidos del sistema
var rangosValidos = []string{
	"Oficial",
	"Primer Oficial",
	"Oficial Jefe",
	"Inspector",
	"Primer Inspector",
	"Inspector Jefe",
	"Comisario",
	"Primer Comisario",
	"Comisario Jefe",
	"Comisario General",
	"Comisario Mayor",
	"Comisario Superior",
	"Subcomisario",
	"Comisario General de Brigada",
	"Comisario General de División",
	"Comisario General Inspector",
	"Comisario General en Jefe",
}

// validarRango verifica si el rango es válido
func validarRango(rango string) bool {
	for _, r := range rangosValidos {
		if r == rango {
			return true
		}
	}
	return false
}

// calcularAntiguedad calcula la antigüedad desde la fecha de graduación
func calcularAntiguedad(fechaGraduacion string) float64 {
	if fechaGraduacion == "" {
		return 0
	}
	fecha, err := time.Parse("2006-01-02", fechaGraduacion)
	if err != nil {
		return 0
	}
	ahora := time.Now()
	antiguedad := ahora.Sub(fecha).Hours() / 24 / 365.25
	return antiguedad
}

// generarQR genera el código QR del oficial (sin información sensible)
func generarQR(oficial *models.Oficial) (string, error) {
	nombreCompleto := oficial.PrimerNombre + " " + oficial.SegundoNombre + " " + oficial.PrimerApellido + " " + oficial.SegundoApellido
	
	// IMPORTANTE: NO incluir información sensible
	datosQR := map[string]interface{}{
		"id":               oficial.ID.Hex(),
		"credencial":       oficial.Credencial,
		"nombre_completo":  nombreCompleto,
		"rango":            oficial.Rango,
		"foto_cara":        oficial.FotoCara,
		"foto_carnet":      oficial.FotoCarnet,
		"destacado":        oficial.Destacado,
		"antiguedad":       oficial.Antiguedad,
		"fecha_graduacion": oficial.FechaGraduacion,
		"fecha_registro":   oficial.FechaRegistro.Format(time.RFC3339),
		// ❌ NO incluir: parientes, licencia_conducir, carnet_medico, contraseña
	}

	jsonData, err := json.Marshal(datosQR)
	if err != nil {
		return "", err
	}

	png, err := qrcode.Encode(string(jsonData), qrcode.Medium, 256)
	if err != nil {
		return "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + qrBase64, nil
}

// RegistrarOficialHandler registra un nuevo oficial
// Requiere autenticación como master
func RegistrarOficialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar que el usuario sea master
	token := r.Header.Get("Authorization")
	master, ok := GetMasterFromToken(token)
	if !ok || master == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación como usuario master para registrar oficiales",
		})
		return
	}

	// Verificar que el master esté activo
	if !master.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario master inactivo",
		})
		return
	}

	// Verificar que tenga permiso de rrhh
	tienePermiso := false
	for _, permiso := range master.Permisos {
		if permiso == "rrhh" {
			tienePermiso = true
			break
		}
	}
	if !tienePermiso {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No tiene permisos para acceder al módulo RRHH",
		})
		return
	}

	var oficial models.Oficial
	if err := json.NewDecoder(r.Body).Decode(&oficial); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if oficial.Credencial == "" {
		http.Error(w, "La credencial es obligatoria", http.StatusBadRequest)
		return
	}

	if oficial.Cedula == "" {
		http.Error(w, "La cédula es obligatoria", http.StatusBadRequest)
		return
	}

	if oficial.Contraseña == "" || len(oficial.Contraseña) < 6 {
		http.Error(w, "La contraseña debe tener al menos 6 caracteres", http.StatusBadRequest)
		return
	}

	if !validarRango(oficial.Rango) {
		http.Error(w, "Rango inválido", http.StatusBadRequest)
		return
	}

	if oficial.FechaGraduacion == "" {
		http.Error(w, "La fecha de graduación es obligatoria", http.StatusBadRequest)
		return
	}

	// Verificar credencial única
	_, err := database.ObtenerOficialPorCredencial(oficial.Credencial)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "La credencial ya está registrada",
		})
		return
	}

	// Verificar cédula única
	_, err = database.ObtenerOficialPorCedula(oficial.Cedula)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "La cédula ya está registrada",
		})
		return
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(oficial.Contraseña), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al procesar la contraseña", http.StatusInternalServerError)
		return
	}
	oficial.Contraseña = string(hashedPassword)

	// Calcular antigüedad si no se proporciona
	if oficial.Antiguedad == 0 {
		oficial.Antiguedad = calcularAntiguedad(oficial.FechaGraduacion)
	}

	// Destacado es opcional - si viene vacío, dejarlo vacío (se asignará en otros módulos)
	if oficial.Destacado == "" {
		oficial.Destacado = ""
	}

	// Generar QR
	qrCode, err := generarQR(&oficial)
	if err != nil {
		log.Printf("Error al generar QR: %v", err)
		// Continuar sin QR si falla
	} else {
		oficial.QRCode = qrCode
	}

	// Crear oficial
	if err := database.CrearOficial(&oficial); err != nil {
		http.Error(w, "Error al crear el oficial: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// No retornar la contraseña
	oficial.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(oficial)
}

// GenerarQRHandler genera o retorna el QR del oficial
func GenerarQRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "ID de oficial requerido", http.StatusBadRequest)
		return
	}

	oficialID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	oficial, err := database.ObtenerOficialPorID(oficialID)
	if err != nil {
		http.Error(w, "Oficial no encontrado", http.StatusNotFound)
		return
	}

	// Si ya tiene QR, retornarlo
	if oficial.QRCode != "" {
		response := map[string]string{
			"qr_code": oficial.QRCode,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generar nuevo QR
	qrCode, err := generarQR(oficial)
	if err != nil {
		http.Error(w, "Error al generar QR", http.StatusInternalServerError)
		return
	}

	// Actualizar oficial con QR
	oficial.QRCode = qrCode
	if err := database.ActualizarOficial(oficial); err != nil {
		log.Printf("Error al actualizar QR: %v", err)
	}

	response := map[string]string{
		"qr_code": qrCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// VerificarQRHandler verifica el QR escaneado
func VerificarQRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener datos del QR de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Datos del QR requeridos", http.StatusBadRequest)
		return
	}

	qrData := pathParts[len(pathParts)-1]

	// Decodificar JSON del QR
	var datosQR map[string]interface{}
	if err := json.Unmarshal([]byte(qrData), &datosQR); err != nil {
		// Si no es JSON, buscar por ID o credencial directamente
		// Intentar como ID primero
		if id, err := primitive.ObjectIDFromHex(qrData); err == nil {
			oficial, err := database.ObtenerOficialPorID(id)
			if err == nil {
				// Retornar información sin datos sensibles
				response := map[string]interface{}{
					"id":               oficial.ID.Hex(),
					"credencial":       oficial.Credencial,
					"primer_nombre":    oficial.PrimerNombre,
					"segundo_nombre":   oficial.SegundoNombre,
					"primer_apellido":  oficial.PrimerApellido,
					"segundo_apellido": oficial.SegundoApellido,
					"rango":            oficial.Rango,
					"destacado":        oficial.Destacado,
					"antiguedad":       oficial.Antiguedad,
					"fecha_graduacion":  oficial.FechaGraduacion,
					"foto_cara":        oficial.FotoCara,
					"foto_carnet":      oficial.FotoCarnet,
					"activo":           oficial.Activo,
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		// Intentar como credencial
		oficial, err := database.ObtenerOficialPorCredencial(qrData)
		if err == nil {
			response := map[string]interface{}{
				"id":               oficial.ID.Hex(),
				"credencial":       oficial.Credencial,
				"primer_nombre":    oficial.PrimerNombre,
				"segundo_nombre":   oficial.SegundoNombre,
				"primer_apellido":  oficial.PrimerApellido,
				"segundo_apellido": oficial.SegundoApellido,
				"rango":            oficial.Rango,
				"destacado":        oficial.Destacado,
				"antiguedad":       oficial.Antiguedad,
				"fecha_graduacion":  oficial.FechaGraduacion,
				"foto_cara":        oficial.FotoCara,
				"foto_carnet":      oficial.FotoCarnet,
				"activo":           oficial.Activo,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
		http.Error(w, "QR no válido o oficial no encontrado", http.StatusNotFound)
		return
	}

	// Si es JSON, buscar por ID o credencial
	var oficial *models.Oficial
	var err error

	if idStr, ok := datosQR["id"].(string); ok {
		if id, err := primitive.ObjectIDFromHex(idStr); err == nil {
			oficial, err = database.ObtenerOficialPorID(id)
		}
	} else if credencial, ok := datosQR["credencial"].(string); ok {
		oficial, err = database.ObtenerOficialPorCredencial(credencial)
	}

	if err != nil || oficial == nil {
		http.Error(w, "QR no válido o oficial no encontrado", http.StatusNotFound)
		return
	}

	// Retornar información sin datos sensibles
	response := map[string]interface{}{
		"id":               oficial.ID.Hex(),
		"credencial":       oficial.Credencial,
		"primer_nombre":    oficial.PrimerNombre,
		"segundo_nombre":   oficial.SegundoNombre,
		"primer_apellido":  oficial.PrimerApellido,
		"segundo_apellido": oficial.SegundoApellido,
		"rango":            oficial.Rango,
		"destacado":        oficial.Destacado,
		"antiguedad":       oficial.Antiguedad,
		"fecha_graduacion": oficial.FechaGraduacion,
		"foto_cara":        oficial.FotoCara,
		"foto_carnet":      oficial.FotoCarnet,
		"activo":           oficial.Activo,
		// ❌ NO incluir: parientes, licencia_conducir, carnet_medico, contraseña, cédula completa
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarOficialesHandler lista oficiales con paginación
func ListarOficialesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parámetros de paginación
	page := 1
	limit := 10
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Filtros
	rango := r.URL.Query().Get("rango")
	estado := r.URL.Query().Get("estado")

	oficiales, total, err := database.ListarOficiales(page, limit, rango, estado)
	if err != nil {
		http.Error(w, "Error al listar oficiales", http.StatusInternalServerError)
		return
	}

	// Ocultar contraseñas
	for i := range oficiales {
		oficiales[i].Contraseña = ""
	}

	response := map[string]interface{}{
		"oficiales": oficiales,
		"total":     total,
		"page":      page,
		"limit":     limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AscensosPendientesHandler lista oficiales con ascensos pendientes
func AscensosPendientesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	oficiales, err := database.ObtenerOficialesConAscensosPendientes()
	if err != nil {
		http.Error(w, "Error al obtener ascensos pendientes", http.StatusInternalServerError)
		return
	}

	// Ocultar contraseñas
	for i := range oficiales {
		oficiales[i].Contraseña = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oficiales)
}

// AprobarAscensoHandler aprueba el ascenso de un oficial
func AprobarAscensoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ID de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "ID de oficial requerido", http.StatusBadRequest)
		return
	}

	oficialID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	oficial, err := database.ObtenerOficialPorID(oficialID)
	if err != nil {
		http.Error(w, "Oficial no encontrado", http.StatusNotFound)
		return
	}

	// Calcular nueva antigüedad basada en fecha de graduación
	if oficial.FechaGraduacion != "" {
		oficial.Antiguedad = calcularAntiguedad(oficial.FechaGraduacion)
	}

	// Actualizar oficial
	if err := database.ActualizarOficial(oficial); err != nil {
		http.Error(w, "Error al actualizar oficial", http.StatusInternalServerError)
		return
	}

	// No retornar la contraseña
	oficial.Contraseña = ""

	response := map[string]interface{}{
		"mensaje":  "Ascenso aprobado correctamente",
		"oficial":  oficial,
		"antiguedad_actualizada": oficial.Antiguedad,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

