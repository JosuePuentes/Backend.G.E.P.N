# Instrucciones para el Backend - Sistema RRHH

Este documento contiene todas las instrucciones para implementar el sistema de Recursos Humanos (RRHH) en el backend.

## 🔐 Sistema de Usuarios Master

### Registro de Usuario Master (Público)

**Endpoint:** `POST /api/rrhh/master/registro`

Este endpoint permite registrar usuarios master que tendrán acceso para gestionar el sistema RRHH.

**Request Body:**
```json
{
  "nombre": "Juan Pérez",
  "email": "juan.perez@gepn.gob.ve",
  "usuario": "jperez",
  "contraseña": "password123"
}
```

**Validaciones:**
- Nombre obligatorio
- Email obligatorio y único
- Usuario obligatorio y único
- Contraseña mínimo 6 caracteres (se hashea con bcrypt)

**Respuesta:**
```json
{
  "mensaje": "Usuario master registrado correctamente",
  "master": {
    "id": "...",
    "nombre": "Juan Pérez",
    "email": "juan.perez@gepn.gob.ve",
    "usuario": "jperez",
    "rol": "master",
    "activo": true,
    "fecha_registro": "2025-01-27T10:00:00Z"
  }
}
```

### Login de Usuario Master

**Endpoint:** `POST /api/rrhh/master/login`

**Request Body:**
```json
{
  "usuario": "jperez",
  "contraseña": "password123"
}
```

**Respuesta:**
```json
{
  "token": "20250127100000-master-token",
  "master": {
    "id": "...",
    "nombre": "Juan Pérez",
    "email": "juan.perez@gepn.gob.ve",
    "usuario": "jperez",
    "rol": "master",
    "activo": true
  },
  "mensaje": "Login exitoso"
}
```

### Verificar Token Master

**Endpoint:** `GET /api/rrhh/master/verificar`

**Headers:**
```
Authorization: <token>
```

**Respuesta:**
```json
{
  "id": "...",
  "nombre": "Juan Pérez",
  "email": "juan.perez@gepn.gob.ve",
  "usuario": "jperez",
  "rol": "master",
  "activo": true
}
```

### ⚠️ Importante: Protección de Endpoints

**Todos los endpoints de RRHH requieren autenticación como master:**

- `POST /api/rrhh/registrar-oficial` - **Requiere token de master**
- `GET /api/rrhh/listar-oficiales` - Puede ser público o requerir autenticación
- `GET /api/rrhh/ascensos-pendientes` - **Requiere token de master**
- `POST /api/rrhh/aprobar-ascenso/:oficialId` - **Requiere token de master**

Para usar estos endpoints, incluir el header:
```
Authorization: <token_obtenido_del_login>
```

---

## 1. Instalar Dependencias

```bash
go get github.com/skip2/go-qrcode
go get golang.org/x/crypto/bcrypt
```

O ejecutar:
```bash
go mod tidy
```

## 2. Modelos de Datos

### Pariente
```go
type Pariente struct {
    Nombre          string `bson:"nombre" json:"nombre"`
    Cedula          string `bson:"cedula" json:"cedula"`
    FechaNacimiento string `bson:"fecha_nacimiento,omitempty" json:"fecha_nacimiento,omitempty"`
}
```

### Parientes
```go
type Parientes struct {
    Padre  *Pariente   `bson:"padre,omitempty" json:"padre,omitempty"`
    Madre  *Pariente   `bson:"madre,omitempty" json:"madre,omitempty"`
    Esposa *Pariente   `bson:"esposa,omitempty" json:"esposa,omitempty"`
    Hijos  []Pariente  `bson:"hijos,omitempty" json:"hijos,omitempty"`
}
```

### Oficial
```go
type Oficial struct {
    ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PrimerNombre      string             `bson:"primer_nombre" json:"primer_nombre"`
    SegundoNombre     string             `bson:"segundo_nombre" json:"segundo_nombre"`
    PrimerApellido    string             `bson:"primer_apellido" json:"primer_apellido"`
    SegundoApellido   string             `bson:"segundo_apellido" json:"segundo_apellido"`
    Cedula            string             `bson:"cedula" json:"cedula"`
    Contraseña        string             `bson:"contraseña" json:"-"`
    FechaNacimiento   string             `bson:"fecha_nacimiento" json:"fecha_nacimiento"`
    Estatura          float64            `bson:"estatura" json:"estatura"`
    ColorPiel         string             `bson:"color_piel" json:"color_piel"`
    TipoSangre        string             `bson:"tipo_sangre" json:"tipo_sangre"`
    CiudadNacimiento  string             `bson:"ciudad_nacimiento" json:"ciudad_nacimiento"`
    Credencial        string             `bson:"credencial" json:"credencial"`
    Rango             string             `bson:"rango" json:"rango"`
    Destacado         string             `bson:"destacado" json:"destacado"`
    FechaGraduacion   string             `bson:"fecha_graduacion" json:"fecha_graduacion"`
    Antiguedad        float64            `bson:"antiguedad" json:"antiguedad"`
    Estado            string             `bson:"estado" json:"estado"`
    Municipio         string             `bson:"municipio" json:"municipio"`
    Parroquia         string             `bson:"parroquia" json:"parroquia"`
    LicenciaConducir  string             `bson:"licencia_conducir,omitempty" json:"licencia_conducir,omitempty"`
    CarnetMedico      string             `bson:"carnet_medico,omitempty" json:"carnet_medico,omitempty"`
    FotoCara          string             `bson:"foto_cara" json:"foto_cara"`
    FotoCarnet        string             `bson:"foto_carnet,omitempty" json:"foto_carnet,omitempty"`
    QRCode            string             `bson:"qr_code" json:"qr_code,omitempty"`
    FechaRegistro     time.Time          `bson:"fecha_registro" json:"fecha_registro"`
    Activo            bool               `bson:"activo" json:"activo"`
    Parientes         *Parientes         `bson:"parientes,omitempty" json:"parientes,omitempty"`
}
```

## 3. Endpoints Implementados

### ⚠️ POST /api/rrhh/registrar-oficial
**REQUIERE AUTENTICACIÓN COMO MASTER**

Registra un nuevo oficial en el sistema. Solo usuarios master pueden registrar oficiales.

**Headers:**
```
Authorization: <token_master>
```

**Request Body:**
```json
{
  "primer_nombre": "Juan",
  "segundo_nombre": "Carlos",
  "primer_apellido": "Pérez",
  "segundo_apellido": "González",
  "cedula": "12345678",
  "contraseña": "password123",
  "fecha_nacimiento": "1990-01-15",
  "estatura": 1.75,
  "color_piel": "Moreno",
  "tipo_sangre": "O+",
  "ciudad_nacimiento": "Caracas",
  "credencial": "POL001",
  "rango": "Oficial",
  "destacado": "Comando Metropolitano",
  "fecha_graduacion": "2015-06-15",
  "estado": "Distrito Capital",
  "municipio": "Libertador",
  "parroquia": "El Recreo",
  "foto_cara": "base64...",
  "parientes": {
    "padre": {
      "nombre": "José Pérez",
      "cedula": "87654321",
      "fecha_nacimiento": "1960-05-20"
    },
    "madre": {
      "nombre": "María González",
      "cedula": "87654322"
    }
  }
}
```

**Validaciones:**
- Credencial única
- Cédula única
- Contraseña mínimo 6 caracteres
- Rango válido
- Fecha de graduación obligatoria

**Respuesta:**
```json
{
  "id": "...",
  "credencial": "POL001",
  "qr_code": "data:image/png;base64,...",
  ...
}
```

### GET /api/rrhh/generar-qr/:oficialId
Genera o retorna el QR del oficial.

**Respuesta:**
```json
{
  "qr_code": "data:image/png;base64,..."
}
```

### GET /api/rrhh/verificar-qr/:qrData
Verifica el QR escaneado y retorna información del oficial (sin datos sensibles).

**Respuesta:**
```json
{
  "id": "...",
  "credencial": "POL001",
  "primer_nombre": "Juan",
  "rango": "Oficial",
  "foto_cara": "...",
  "antiguedad": 8.5,
  ...
}
```

**⚠️ NO incluye:**
- Parientes
- Licencia de conducir
- Carnet médico
- Contraseña
- Cédula completa

### GET /api/rrhh/listar-oficiales
Lista oficiales con paginación.

**Query Parameters:**
- `page` (opcional): Número de página (default: 1)
- `limit` (opcional): Cantidad por página (default: 10, max: 100)
- `rango` (opcional): Filtrar por rango
- `estado` (opcional): Filtrar por estado

**Ejemplo:**
```
GET /api/rrhh/listar-oficiales?page=1&limit=20&rango=Oficial
```

**Respuesta:**
```json
{
  "oficiales": [...],
  "total": 150,
  "page": 1,
  "limit": 20
}
```

### GET /api/rrhh/ascensos-pendientes
Lista oficiales con ascensos pendientes (cada 4 años desde fecha de graduación).

**Respuesta:**
```json
[
  {
    "id": "...",
    "credencial": "POL001",
    "nombre_completo": "Juan Carlos Pérez González",
    "rango": "Oficial",
    "antiguedad": 8.5,
    "fecha_graduacion": "2015-06-15",
    ...
  }
]
```

### POST /api/rrhh/aprobar-ascenso/:oficialId
Aprueba el ascenso de un oficial, actualizando su antigüedad.

**Respuesta:**
```json
{
  "mensaje": "Ascenso aprobado correctamente",
  "oficial": {...},
  "antiguedad_actualizada": 8.5
}
```

## 4. Modificación del Login Policial

El endpoint `/api/policial/login` ahora:
1. Busca por credencial en la colección `oficiales`
2. Verifica contraseña con bcrypt
3. Verifica que esté activo
4. Crea/actualiza guardia con GPS
5. Retorna token JWT

**Request:**
```json
{
  "credencial": "POL001",
  "pin": "password123",
  "latitud": 10.4969,
  "longitud": -66.8983
}
```

## 5. Generación de QR

El QR contiene **SOLO información pública**:

```go
datosQR := map[string]interface{}{
    "id": oficial.ID.Hex(),
    "credencial": oficial.Credencial,
    "nombre_completo": nombreCompleto,
    "rango": oficial.Rango,
    "foto_cara": oficial.FotoCara,
    "foto_carnet": oficial.FotoCarnet,
    "destacado": oficial.Destacado,
    "antiguedad": oficial.Antiguedad,
    "fecha_graduacion": oficial.FechaGraduacion,
    "fecha_registro": oficial.FechaRegistro.Format(time.RFC3339),
}
```

**❌ NO incluir:**
- Parientes
- Licencia de conducir
- Carnet médico
- Contraseña
- Cédula completa

## 6. Sistema de Ascensos Automáticos

- Se calcula la antigüedad desde la fecha de graduación
- Cada 4 años = 1 ascenso
- Los ascensos se verifican automáticamente
- Deben aprobarse manualmente mediante el endpoint `/api/rrhh/aprobar-ascenso/:oficialId`

**Cálculo de antigüedad:**
```go
func calcularAntiguedad(fechaGraduacion string) float64 {
    fecha, _ := time.Parse("2006-01-02", fechaGraduacion)
    ahora := time.Now()
    antiguedad := ahora.Sub(fecha).Hours() / 24 / 365.25
    return antiguedad
}
```

## 7. Índices MongoDB

Los siguientes índices se crean automáticamente:

```javascript
db.oficiales.createIndex({ "credencial": 1 }, { unique: true })
db.oficiales.createIndex({ "cedula": 1 }, { unique: true })
db.oficiales.createIndex({ "rango": 1 })
db.oficiales.createIndex({ "estado": 1 })
```

## 8. Validaciones Importantes

### Credencial
- Debe ser única
- No puede estar vacía

### Cédula
- Debe ser única
- No puede estar vacía

### Contraseña
- Mínimo 6 caracteres
- Se hashea con bcrypt (salt rounds >= 10)

### Rango
Debe ser uno de los siguientes:
- Oficial
- Primer Oficial
- Oficial Jefe
- Inspector
- Primer Inspector
- Inspector Jefe
- Comisario
- Primer Comisario
- Comisario Jefe
- Comisario General
- Comisario Mayor
- Comisario Superior

### Fecha de Graduación
- Obligatoria
- Formato: `YYYY-MM-DD`
- Se usa para calcular antigüedad automáticamente

### Antigüedad
- Se calcula automáticamente si no se proporciona
- Basada en la fecha de graduación

## 9. Seguridad

### Autenticación
- Los endpoints `/api/rrhh/*` requieren autenticación (implementar middleware)
- El login usa bcrypt para verificar contraseñas

### Hash de Contraseñas
```go
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(contraseña), bcrypt.DefaultCost)
```

### Validación y Sanitización
- Todos los datos se validan antes de guardar
- Las contraseñas nunca se retornan en las respuestas
- Los datos sensibles no se incluyen en el QR

### Rate Limiting
- Implementar rate limiting en producción
- Especialmente en endpoints de login y registro

## 10. Ejemplo de Handler Completo

Ver el archivo `handlers/rrhh.go` para la implementación completa de todos los handlers.

## 11. Funciones de Repositorio

Las siguientes funciones están disponibles en `database/repositories.go`:

- `CrearOficial(oficial *models.Oficial) error`
- `ObtenerOficialPorCredencial(credencial string) (*models.Oficial, error)`
- `ObtenerOficialPorID(id primitive.ObjectID) (*models.Oficial, error)`
- `ObtenerOficialPorCedula(cedula string) (*models.Oficial, error)`
- `ActualizarOficial(oficial *models.Oficial) error`
- `ListarOficiales(page, limit int, rango, estado string) ([]models.Oficial, int64, error)`
- `ObtenerOficialesConAscensosPendientes() ([]models.Oficial, error)`
- `VerificarOficialPorQR(qrData string) (*models.Oficial, error)`

## 12. Notas Importantes

1. **Migración de Usuarios**: Los usuarios existentes en la colección `usuarios` seguirán funcionando, pero el login ahora busca primero en `oficiales`.

2. **Compatibilidad**: El login retorna un objeto `Usuario` para mantener compatibilidad con el frontend existente.

3. **QR Code**: El QR se genera automáticamente al registrar un oficial, pero puede regenerarse con el endpoint correspondiente.

4. **Ascensos**: El sistema calcula automáticamente los ascensos pendientes, pero deben aprobarse manualmente.

5. **Datos Sensibles**: Nunca incluir información sensible en el QR o en respuestas públicas.

## 13. Próximos Pasos

1. Implementar middleware de autenticación para endpoints de RRHH
2. Agregar validación de permisos (solo administradores pueden registrar oficiales)
3. Implementar JWT tokens en lugar de tokens simples
4. Agregar logs de auditoría para cambios en oficiales
5. Implementar rate limiting
6. Agregar tests unitarios

---

**Estado**: ✅ Implementación completa
**Fecha**: 2025-01-27
**Versión**: 1.0.0

