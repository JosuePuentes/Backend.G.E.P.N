# Instrucciones para el Backend - Sistema Centro de Coordinación

Este documento contiene las instrucciones para implementar el sistema de Centro de Coordinación con control de permisos.

## 1. Modelos de Datos

En `models/models.go`, agregar:

```go
// FuncionarioAsignado representa un funcionario asignado a una estación
type FuncionarioAsignado struct {
    FuncionarioID primitive.ObjectID `bson:"funcionario_id" json:"funcionario_id"`
    Nombre        string             `bson:"nombre" json:"nombre"`
    Credencial    string             `bson:"credencial" json:"credencial"`
    Rango         string             `bson:"rango" json:"rango"`
    FechaAsignacion time.Time         `bson:"fecha_asignacion" json:"fecha_asignacion"`
    Activo        bool               `bson:"activo" json:"activo"`
}

// Parte representa un parte de servicio
type Parte struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    NumeroParte     string             `bson:"numero_parte" json:"numero_parte"`
    EstacionID      primitive.ObjectID `bson:"estacion_id" json:"estacion_id"`
    FuncionarioID   primitive.ObjectID `bson:"funcionario_id" json:"funcionario_id"`
    TipoParte       string             `bson:"tipo_parte" json:"tipo_parte"` // "entrada", "salida", "incidente"
    Descripcion     string             `bson:"descripcion" json:"descripcion"`
    FechaHora       time.Time          `bson:"fecha_hora" json:"fecha_hora"`
    Ubicacion       string             `bson:"ubicacion,omitempty" json:"ubicacion,omitempty"`
    Latitud         float64            `bson:"latitud,omitempty" json:"latitud,omitempty"`
    Longitud        float64            `bson:"longitud,omitempty" json:"longitud,omitempty"`
    Estado          string             `bson:"estado" json:"estado"` // "activo", "cerrado", "cancelado"
    FechaCreacion   time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
}

// Estacion representa una estación policial
type Estacion struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Nombre          string             `bson:"nombre" json:"nombre"`
    Codigo          string             `bson:"codigo" json:"codigo"` // Único
    CentroID       primitive.ObjectID `bson:"centro_id" json:"centro_id"`
    Estado          string             `bson:"estado" json:"estado"`
    Municipio       string             `bson:"municipio" json:"municipio"`
    Parroquia       string             `bson:"parroquia" json:"parroquia"`
    Direccion       string             `bson:"direccion" json:"direccion"`
    Telefono        string             `bson:"telefono,omitempty" json:"telefono,omitempty"`
    FuncionariosAsignados []FuncionarioAsignado `bson:"funcionarios_asignados" json:"funcionarios_asignados"`
    Activa          bool               `bson:"activa" json:"activa"`
    FechaCreacion   time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
}

// Centro representa un centro de coordinación
type Centro struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Nombre          string             `bson:"nombre" json:"nombre"`
    Codigo          string             `bson:"codigo" json:"codigo"` // Único
    Estado          string             `bson:"estado" json:"estado"` // Estado/Región
    Municipio       string             `bson:"municipio,omitempty" json:"municipio,omitempty"`
    Direccion       string             `bson:"direccion" json:"direccion"`
    Telefono        string             `bson:"telefono,omitempty" json:"telefono,omitempty"`
    Responsable     string             `bson:"responsable,omitempty" json:"responsable,omitempty"`
    Activo          bool               `bson:"activo" json:"activo"`
    FechaCreacion   time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
}
```

## 2. Endpoints Requeridos

### POST /api/centro-coordinacion/centros
Crear nuevo centro de coordinación
- **Solo admin puede crear centros**
- Validar código único

### GET /api/centro-coordinacion/centros
Listar centros
- **Admin:** Ve todos los centros
- **RRHH Regional:** Solo ve centros de su estado/región

### POST /api/centro-coordinacion/estaciones
Crear nueva estación
- **Solo admin puede crear estaciones**
- Validar código único
- Asociar a un centro

### GET /api/centro-coordinacion/estaciones
Listar estaciones
- **Admin:** Ve todas las estaciones
- **RRHH Regional:** Solo ve estaciones de su estado/región
- Filtro opcional por centro_id

### POST /api/centro-coordinacion/estaciones/:estacionId/asignar-funcionario
Asignar funcionario a estación
- **Solo admin puede asignar**
- Validar que el funcionario exista

### GET /api/centro-coordinacion/estaciones/:estacionId/funcionarios
Listar funcionarios de una estación
- **Admin:** Ve todos
- **RRHH Regional:** Solo si la estación está en su estado

### POST /api/centro-coordinacion/partes
Crear parte de servicio
- Requiere autenticación
- Validar que el funcionario esté asignado a la estación

### GET /api/centro-coordinacion/partes
Listar partes
- **Admin:** Ve todos los partes
- **RRHH Regional:** Solo partes de su estado/región
- Filtros opcionales: estacion_id, funcionario_id, tipo_parte, fecha

## 3. Control de Permisos

### Usuario Admin
- Ve **TODOS** los centros, estaciones, funcionarios y partes
- Puede crear centros y estaciones
- Puede asignar funcionarios a cualquier estación

### Usuario RRHH Regional
- Solo ve datos de **su estado/región**
- No puede crear centros ni estaciones
- No puede asignar funcionarios
- Puede ver funcionarios asignados a estaciones de su estado
- Puede ver partes de su estado

## 4. Campo de Estado en UsuarioMaster

Agregar campo `Estado` al modelo `UsuarioMaster` para identificar la región:

```go
type UsuarioMaster struct {
    // ... campos existentes ...
    Estado          string             `bson:"estado,omitempty" json:"estado,omitempty"` // Para usuarios RRHH regionales
}
```

## 5. Verificación de Permisos

```go
// Verificar si es admin (tiene todos los permisos)
func esAdmin(master *models.UsuarioMaster) bool {
    return len(master.Permisos) == len(ModulosDisponibles) && 
           master.Usuario == "admin"
}

// Verificar acceso a estado/región
func tieneAccesoEstado(master *models.UsuarioMaster, estado string) bool {
    if esAdmin(master) {
        return true // Admin ve todo
    }
    return master.Estado == estado // RRHH regional solo su estado
}
```

## 6. Índices MongoDB

```javascript
db.centros.createIndex({ "codigo": 1 }, { unique: true })
db.centros.createIndex({ "estado": 1 })
db.estaciones.createIndex({ "codigo": 1 }, { unique: true })
db.estaciones.createIndex({ "centro_id": 1 })
db.estaciones.createIndex({ "estado": 1 })
db.partes.createIndex({ "estacion_id": 1 })
db.partes.createIndex({ "funcionario_id": 1 })
db.partes.createIndex({ "fecha_hora": -1 })
```

## 7. Notas Importantes

1. **Solo el admin puede crear centros y estaciones**
2. **Solo el admin puede asignar funcionarios**
3. **RRHH regional solo ve su estado/región**
4. **El campo `Estado` en UsuarioMaster identifica la región del usuario RRHH**
5. **Los funcionarios asignados deben existir en la colección `oficiales`**

---

**Estado**: ⚠️ Pendiente de implementación
**Fecha**: 2025-12-22


