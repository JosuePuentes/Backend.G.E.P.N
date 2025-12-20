package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Usuario representa un usuario del sistema (policial)
type Usuario struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Credencial   string            `bson:"credencial" json:"credencial"`
	PIN          string            `bson:"pin" json:"-"` // No se expone en JSON
	Nombre       string            `bson:"nombre" json:"nombre"`
	Rango        string            `bson:"rango" json:"rango"`
	Activo       bool              `bson:"activo" json:"activo"`
	FechaCreacion time.Time        `bson:"fecha_creacion" json:"fecha_creacion"`
}

// Detenido representa un registro de detenido
type Detenido struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Cedula          string            `bson:"cedula" json:"cedula"`
	Nombre          string            `bson:"nombre" json:"nombre"`
	Apellido        string            `bson:"apellido" json:"apellido"`
	FechaNacimiento string            `bson:"fecha_nacimiento" json:"fecha_nacimiento"`
	Direccion       string            `bson:"direccion" json:"direccion"`
	Motivo          string            `bson:"motivo" json:"motivo"`
	Ubicacion       string            `bson:"ubicacion" json:"ubicacion"`
	Latitud         float64           `bson:"latitud" json:"latitud"`
	Longitud        float64           `bson:"longitud" json:"longitud"`
	OficialID       primitive.ObjectID `bson:"oficial_id" json:"oficial_id"`
	FechaDetencion  time.Time          `bson:"fecha_detencion" json:"fecha_detencion"`
	Estado          string             `bson:"estado" json:"estado"` // "detenido", "liberado", "trasladado"
}

// Minuta representa una minuta digital
type Minuta struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Titulo        string            `bson:"titulo" json:"titulo"`
	Descripcion   string            `bson:"descripcion" json:"descripcion"`
	Tipo          string            `bson:"tipo" json:"tipo"` // "patrullaje", "incidente", "reporte"
	Ubicacion     string            `bson:"ubicacion" json:"ubicacion"`
	Latitud      float64           `bson:"latitud" json:"latitud"`
	Longitud      float64           `bson:"longitud" json:"longitud"`
	OficialID     primitive.ObjectID `bson:"oficial_id" json:"oficial_id"`
	FechaCreacion time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
	Archivos      []string           `bson:"archivos,omitempty" json:"archivos,omitempty"`
}

// BusquedaCedula representa una búsqueda de cédula
type BusquedaCedula struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Cedula        string            `bson:"cedula" json:"cedula"`
	Nombre        string            `bson:"nombre,omitempty" json:"nombre,omitempty"`
	Apellido      string            `bson:"apellido,omitempty" json:"apellido,omitempty"`
	Resultado     string            `bson:"resultado" json:"resultado"` // "encontrado", "no_encontrado", "buscado"
	FechaBusqueda time.Time         `bson:"fecha_busqueda" json:"fecha_busqueda"`
	OficialID     primitive.ObjectID `bson:"oficial_id" json:"oficial_id"`
}

// MasBuscado representa una persona más buscada
type MasBuscado struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Cedula       string            `bson:"cedula" json:"cedula"`
	Nombre       string            `bson:"nombre" json:"nombre"`
	Apellido     string            `bson:"apellido" json:"apellido"`
	Foto         string            `bson:"foto,omitempty" json:"foto,omitempty"`
	Motivo       string            `bson:"motivo" json:"motivo"`
	Prioridad    string            `bson:"prioridad" json:"prioridad"` // "alta", "media", "baja"
	VecesBuscado int               `bson:"veces_buscado" json:"veces_buscado"`
}

// Panico representa una alerta de pánico
type Panico struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OficialID       primitive.ObjectID `bson:"oficial_id" json:"oficial_id"`
	Latitud         float64            `bson:"latitud" json:"latitud"`
	Longitud        float64            `bson:"longitud" json:"longitud"`
	Ubicacion       string             `bson:"ubicacion" json:"ubicacion"`
	FechaActivacion time.Time          `bson:"fecha_activacion" json:"fecha_activacion"`
	Estado          string             `bson:"estado" json:"estado"` // "activo", "atendido", "cancelado"
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

