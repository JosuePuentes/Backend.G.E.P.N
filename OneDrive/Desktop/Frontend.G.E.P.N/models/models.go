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
	EnGuardia    bool              `bson:"en_guardia" json:"en_guardia"`
	Latitud      float64           `bson:"latitud,omitempty" json:"latitud,omitempty"`
	Longitud     float64           `bson:"longitud,omitempty" json:"longitud,omitempty"`
}

// Guardia representa una guardia activa de un oficial
type Guardia struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OficialID     primitive.ObjectID `bson:"oficial_id" json:"oficial_id"`
	FechaInicio   time.Time          `bson:"fecha_inicio" json:"fecha_inicio"`
	FechaFin      *time.Time         `bson:"fecha_fin,omitempty" json:"fecha_fin,omitempty"`
	LatitudInicio float64            `bson:"latitud_inicio" json:"latitud_inicio"`
	LongitudInicio float64           `bson:"longitud_inicio" json:"longitud_inicio"`
	Activa        bool               `bson:"activa" json:"activa"`
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
	Credencial string  `json:"credencial"`
	PIN        string  `json:"pin"`
	Latitud    float64 `json:"latitud,omitempty"`
	Longitud   float64 `json:"longitud,omitempty"`
}

// LoginResponse representa la respuesta de login
type LoginResponse struct {
	Token   string  `json:"token"`
	Usuario Usuario `json:"usuario"`
}

// Ciudadano representa un ciudadano registrado
type Ciudadano struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre        string            `bson:"nombre" json:"nombre"`
	Cedula        string            `bson:"cedula" json:"cedula"`
	Telefono      string            `bson:"telefono" json:"telefono"`
	Contraseña    string            `bson:"contraseña" json:"-"` // No se expone en JSON
	FechaRegistro time.Time         `bson:"fecha_registro" json:"fecha_registro"`
	Activo        bool              `bson:"activo" json:"activo"`
}

// Denuncia representa una denuncia realizada por un ciudadano
type Denuncia struct {
	ID                        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CiudadanoID               primitive.ObjectID `bson:"ciudadano_id" json:"ciudadano_id"`
	NumeroDenuncia            string            `bson:"numero_denuncia" json:"numero_denuncia"`
	// Datos del denunciante
	NombreDenunciante         string            `bson:"nombre_denunciante" json:"nombre_denunciante"`
	CedulaDenunciante         string            `bson:"cedula_denunciante" json:"cedula_denunciante"`
	TelefonoDenunciante       string            `bson:"telefono_denunciante" json:"telefono_denunciante"`
	FechaNacimientoDenunciante string           `bson:"fecha_nacimiento_denunciante,omitempty" json:"fecha_nacimiento_denunciante,omitempty"`
	ParroquiaDenunciante      string            `bson:"parroquia_denunciante,omitempty" json:"parroquia_denunciante,omitempty"`
	// Datos de la denuncia
	Motivo                    string            `bson:"motivo" json:"motivo"`
	Hechos                    string            `bson:"hechos" json:"hechos"`
	// Datos del denunciado
	NombreDenunciado          string            `bson:"nombre_denunciado,omitempty" json:"nombre_denunciado,omitempty"`
	DireccionDenunciado       string            `bson:"direccion_denunciado,omitempty" json:"direccion_denunciado,omitempty"`
	EstadoDenunciado          string            `bson:"estado_denunciado,omitempty" json:"estado_denunciado,omitempty"`
	MunicipioDenunciado       string            `bson:"municipio_denunciado,omitempty" json:"municipio_denunciado,omitempty"`
	ParroquiaDenunciado       string            `bson:"parroquia_denunciado,omitempty" json:"parroquia_denunciado,omitempty"`
	// Metadatos
	FechaDenuncia             time.Time         `bson:"fecha_denuncia" json:"fecha_denuncia"`
	Estado                    string            `bson:"estado" json:"estado"` // "Pendiente", "En Proceso", "Resuelta", "Archivada"
}

// RegistroCiudadanoRequest representa la petición de registro
type RegistroCiudadanoRequest struct {
	Nombre     string `json:"nombre"`
	Cedula     string `json:"cedula"`
	Telefono   string `json:"telefono"`
	Contraseña string `json:"contraseña"`
}

// LoginCiudadanoRequest representa la petición de login
type LoginCiudadanoRequest struct {
	Cedula     string `json:"cedula"`
	Contraseña string `json:"contraseña"`
}

// CrearDenunciaRequest representa la petición de crear denuncia
type CrearDenunciaRequest struct {
	Denunciante struct {
		Nombre          string `json:"nombre"`
		Cedula          string `json:"cedula"`
		Telefono        string `json:"telefono"`
		FechaNacimiento string `json:"fechaNacimiento,omitempty"`
		Parroquia       string `json:"parroquia,omitempty"`
	} `json:"denunciante"`
	Denuncia struct {
		Motivo string `json:"motivo"`
		Hechos string `json:"hechos"`
	} `json:"denuncia"`
	Denunciado struct {
		Nombre    string `json:"nombre,omitempty"`
		Direccion string `json:"direccion,omitempty"`
		Estado    string `json:"estado,omitempty"`
		Municipio string `json:"municipio,omitempty"`
		Parroquia string `json:"parroquia,omitempty"`
	} `json:"denunciado,omitempty"`
}

