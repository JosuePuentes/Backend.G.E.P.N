package models

import "time"

// Usuario representa un usuario del sistema (policial)
type Usuario struct {
	ID           int       `json:"id"`
	Credencial   string    `json:"credencial"`
	PIN          string    `json:"-"` // No se expone en JSON
	Nombre       string    `json:"nombre"`
	Rango        string    `json:"rango"`
	Activo       bool      `json:"activo"`
	FechaCreacion time.Time `json:"fecha_creacion"`
}

// Detenido representa un registro de detenido
type Detenido struct {
	ID              int       `json:"id"`
	Cedula          string    `json:"cedula"`
	Nombre          string    `json:"nombre"`
	Apellido        string    `json:"apellido"`
	FechaNacimiento string    `json:"fecha_nacimiento"`
	Direccion       string    `json:"direccion"`
	Motivo          string    `json:"motivo"`
	Ubicacion       string    `json:"ubicacion"`
	Latitud         float64   `json:"latitud"`
	Longitud        float64   `json:"longitud"`
	OficialID       int       `json:"oficial_id"`
	FechaDetencion  time.Time `json:"fecha_detencion"`
	Estado          string    `json:"estado"` // "detenido", "liberado", "trasladado"
}

// Minuta representa una minuta digital
type Minuta struct {
	ID            int       `json:"id"`
	Titulo        string    `json:"titulo"`
	Descripcion   string    `json:"descripcion"`
	Tipo          string    `json:"tipo"` // "patrullaje", "incidente", "reporte"
	Ubicacion     string    `json:"ubicacion"`
	Latitud      float64   `json:"latitud"`
	Longitud      float64   `json:"longitud"`
	OficialID     int       `json:"oficial_id"`
	FechaCreacion time.Time `json:"fecha_creacion"`
	Archivos      []string  `json:"archivos,omitempty"`
}

// BusquedaCedula representa una búsqueda de cédula
type BusquedaCedula struct {
	Cedula      string    `json:"cedula"`
	Nombre      string    `json:"nombre,omitempty"`
	Apellido    string    `json:"apellido,omitempty"`
	Resultado   string    `json:"resultado"` // "encontrado", "no_encontrado", "buscado"
	FechaBusqueda time.Time `json:"fecha_busqueda"`
	OficialID   int       `json:"oficial_id"`
}

// MasBuscado representa una persona más buscada
type MasBuscado struct {
	ID          int      `json:"id"`
	Cedula      string   `json:"cedula"`
	Nombre      string   `json:"nombre"`
	Apellido    string   `json:"apellido"`
	Foto        string   `json:"foto,omitempty"`
	Motivo      string   `json:"motivo"`
	Prioridad   string   `json:"prioridad"` // "alta", "media", "baja"
	VecesBuscado int     `json:"veces_buscado"`
}

// Panico representa una alerta de pánico
type Panico struct {
	ID            int       `json:"id"`
	OficialID     int       `json:"oficial_id"`
	Latitud       float64   `json:"latitud"`
	Longitud      float64   `json:"longitud"`
	Ubicacion     string    `json:"ubicacion"`
	FechaActivacion time.Time `json:"fecha_activacion"`
	Estado        string    `json:"estado"` // "activo", "atendido", "cancelado"
}

// LoginRequest representa la petición de login
type LoginRequest struct {
	Credencial string `json:"credencial"`
	PIN        string `json:"pin"`
}

// LoginResponse representa la respuesta de login
type LoginResponse struct {
	Token   string  `json:"token"`
	Usuario Usuario `json:"usuario"`
}

