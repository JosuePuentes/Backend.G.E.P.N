# Instrucciones para el Backend - Sistema Master

Este documento contiene todas las instrucciones para el sistema de usuarios master con JWT y permisos por módulos.

## 1. Instalar Dependencias

```bash
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v4
```

O ejecutar:
```bash
go mod tidy
```

## 2. Modelo UsuarioMaster

```go
type UsuarioMaster struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Usuario       string             `bson:"usuario" json:"usuario"` // Único
    Nombre        string             `bson:"nombre" json:"nombre"`
    Email         string             `bson:"email" json:"email"` // Único
    Contraseña    string             `bson:"contraseña" json:"-"` // Hash
    Permisos      []string           `bson:"permisos" json:"permisos"` // ["rrhh", "policial", etc.]
    Activo        bool               `bson:"activo" json:"activo"`
    CreadoPor     string             `bson:"creado_por" json:"creado_por"`
    FechaCreacion time.Time          `bson:"fecha_creacion" json:"fecha_creacion"`
    UltimoAcceso  *time.Time         `bson:"ultimo_acceso,omitempty" json:"ultimo_acceso,omitempty"`
}
```

## 3. Módulos Disponibles

```go
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
```

## 4. Endpoints Implementados

### POST /api/master/login
Login para usuarios master. Retorna JWT token.

**Request:**
```json
{
  "usuario": "admin",
  "contraseña": "Admin123!"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "master": {
    "id": "...",
    "usuario": "admin",
    "nombre": "Administrador",
    "email": "admin@gepn.gob.ve",
    "permisos": ["rrhh", "policial", "denuncias", ...],
    "activo": true
  },
  "mensaje": "Login exitoso"
}
```

### POST /api/master/crear-usuario (Protegida)
Crea nuevo usuario master. Requiere autenticación.

**Headers:**
```
Authorization: <jwt_token>
```

**Request:**
```json
{
  "usuario": "jperez",
  "nombre": "Juan Pérez",
  "email": "juan.perez@gepn.gob.ve",
  "contraseña": "password123",
  "permisos": ["rrhh", "policial"]
}
```

**Validaciones:**
- Usuario único
- Email único
- Contraseña mínimo 6 caracteres
- Permisos válidos (debe ser uno de los módulos disponibles)

### GET /api/master/usuarios (Protegida)
Lista todos los usuarios master.

**Headers:**
```
Authorization: <jwt_token>
```

**Response:**
```json
[
  {
    "id": "...",
    "usuario": "admin",
    "nombre": "Administrador",
    "email": "admin@gepn.gob.ve",
    "permisos": ["rrhh", "policial", ...],
    "activo": true,
    "creado_por": "sistema",
    "fecha_creacion": "2025-01-27T10:00:00Z"
  }
]
```

### PUT /api/master/usuarios/permisos/:usuarioId (Protegida)
Actualiza permisos de un usuario.

**Headers:**
```
Authorization: <jwt_token>
```

**Request:**
```json
{
  "permisos": ["rrhh", "policial"]
}
```

### PUT /api/master/usuarios/activar/:usuarioId (Protegida)
Activa/desactiva un usuario.

**Headers:**
```
Authorization: <jwt_token>
```

**Request:**
```json
{
  "activo": true
}
```

### GET /api/master/modulos
Retorna lista de módulos disponibles (público).

**Response:**
```json
{
  "modulos": [
    "rrhh",
    "policial",
    "denuncias",
    "detenidos",
    "minutas",
    "buscados",
    "verificacion",
    "panico"
  ]
}
```

### GET /api/master/verificar (Protegida)
Verifica el token y retorna información del master.

**Headers:**
```
Authorization: <jwt_token>
```

## 5. Usuario Master Inicial

El sistema crea automáticamente un usuario admin al iniciar:

- **Usuario:** `admin`
- **Contraseña:** `Admin123!` (⚠️ CAMBIAR EN PRODUCCIÓN)
- **Email:** `admin@gepn.gob.ve`
- **Permisos:** Todos los módulos disponibles
- **Activo:** `true`

Este usuario se crea automáticamente en `main.go` al iniciar el servidor.

## 6. Sistema de Permisos

### Verificación de Permisos

Los endpoints protegidos verifican que el usuario tenga el permiso necesario:

```go
// Ejemplo en RegistrarOficialHandler
tienePermiso := false
for _, permiso := range master.Permisos {
    if permiso == "rrhh" {
        tienePermiso = true
        break
    }
}
if !tienePermiso {
    // Retornar 403 Forbidden
}
```

### Endpoints Protegidos por Módulo

- **RRHH:** `/api/rrhh/*` - Requiere permiso `"rrhh"`
- **Policial:** `/api/policial/*` - Requiere permiso `"policial"`
- **Denuncias:** `/api/denuncia/*` - Requiere permiso `"denuncias"`
- **Detenidos:** `/api/detenidos/*` - Requiere permiso `"detenidos"`
- **Minutas:** `/api/minutas/*` - Requiere permiso `"minutas"`
- **Buscados:** `/api/mas-buscados` - Requiere permiso `"buscados"`
- **Verificación:** `/api/buscar/cedula` - Requiere permiso `"verificacion"`
- **Pánico:** `/api/panico/*` - Requiere permiso `"panico"`

## 7. JWT Tokens

### Generación de Token

```go
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
```

### Verificación de Token

```go
claims := &MasterClaims{}
tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
    return jwtSecret, nil
})
```

### JWT Secret

El JWT secret se obtiene de la variable de entorno `JWT_SECRET`. Si no está configurada, usa un valor por defecto (⚠️ cambiar en producción).

```go
var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "gepn-secret-key-change-in-production"
    }
    return secret
}
```

## 8. Índices MongoDB

Los siguientes índices se crean automáticamente:

```javascript
db.usuarios_master.createIndex({ "usuario": 1 }, { unique: true })
db.usuarios_master.createIndex({ "email": 1 }, { unique: true })
db.usuarios_master.createIndex({ "permisos": 1 })
```

## 9. Flujo de Trabajo

1. **Sistema inicia** → Crea usuario admin automáticamente
2. **Admin accede** → Login en `/api/master/login` con usuario `admin` y contraseña `Admin123!`
3. **Admin crea usuarios** → Usa `/api/master/crear-usuario` para crear otros usuarios master
4. **Admin asigna permisos** → Usa `/api/master/usuarios/permisos/:id` para asignar permisos
5. **Usuario con permiso rrhh** → Puede acceder a `/api/rrhh/*` y crear oficiales
6. **Oficiales creados** → Acceden con `/api/policial/login` usando su credencial y contraseña

## 10. Seguridad

### Hash de Contraseñas
- Usa bcrypt con `bcrypt.DefaultCost` (10 salt rounds)
- Las contraseñas nunca se retornan en las respuestas

### JWT Tokens
- Tokens expiran después de 24 horas
- Secret key debe configurarse en variable de entorno
- Tokens incluyen información del usuario y permisos

### Validación de Permisos
- Cada endpoint protegido verifica que el usuario tenga el permiso necesario
- Si no tiene permiso, retorna `403 Forbidden`

### Rate Limiting
- ⚠️ Implementar rate limiting en producción
- Especialmente en endpoints de login

### Logs de Auditoría
- ⚠️ Implementar logs de auditoría para cambios en usuarios master
- Registrar quién creó/modificó cada usuario

## 11. Funciones de Repositorio

Las siguientes funciones están disponibles en `database/repositories.go`:

- `CrearUsuarioMaster(master *models.UsuarioMaster) error`
- `ObtenerUsuarioMasterPorUsuario(usuario string) (*models.UsuarioMaster, error)`
- `ObtenerUsuarioMasterPorID(id primitive.ObjectID) (*models.UsuarioMaster, error)`
- `ObtenerUsuarioMasterPorEmail(email string) (*models.UsuarioMaster, error)`
- `ActualizarUltimoAccesoMaster(masterID primitive.ObjectID) error`
- `ListarUsuariosMaster() ([]models.UsuarioMaster, error)`
- `ActualizarPermisosMaster(masterID primitive.ObjectID, permisos []string) error`
- `ActualizarEstadoMaster(masterID primitive.ObjectID, activo bool) error`

## 12. Variables de Entorno

Configurar las siguientes variables de entorno:

```bash
JWT_SECRET=tu-secret-key-super-seguro-aqui
MONGODB_URI=mongodb+srv://...
MONGODB_DB_NAME=gepn
PORT=8080
```

## 13. Ejemplo de Uso

### 1. Login como Admin
```bash
curl -X POST http://localhost:8080/api/master/login \
  -H "Content-Type: application/json" \
  -d '{
    "usuario": "admin",
    "contraseña": "Admin123!"
  }'
```

### 2. Crear Usuario Master
```bash
curl -X POST http://localhost:8080/api/master/crear-usuario \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_jwt>" \
  -d '{
    "usuario": "jperez",
    "nombre": "Juan Pérez",
    "email": "juan.perez@gepn.gob.ve",
    "contraseña": "password123",
    "permisos": ["rrhh", "policial"]
  }'
```

### 3. Registrar Oficial (requiere permiso rrhh)
```bash
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_jwt>" \
  -d '{
    "primer_nombre": "Juan",
    "credencial": "POL001",
    ...
  }'
```

## 14. Notas Importantes

1. **Cambiar contraseña del admin** en producción
2. **Configurar JWT_SECRET** en variable de entorno
3. **Implementar rate limiting** en producción
4. **Implementar logs de auditoría** para cambios en usuarios
5. **Validar permisos** en cada endpoint protegido
6. **Tokens expiran** después de 24 horas

---

**Estado**: ✅ Implementación completa
**Fecha**: 2025-01-27
**Versión**: 2.0.0

