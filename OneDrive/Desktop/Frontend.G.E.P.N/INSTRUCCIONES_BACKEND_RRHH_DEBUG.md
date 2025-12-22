# Instrucciones de Debug - Backend RRHH

Este documento contiene instrucciones para verificar y depurar el endpoint de registro de oficiales.

## ✅ Verificaciones del Backend

### 1. Endpoint Verificado

**Endpoint:** `POST /api/rrhh/registrar-oficial`

**Ubicación:** `routes/routes.go` línea 63

```go
mux.HandleFunc("/api/rrhh/registrar-oficial", handlers.RegistrarOficialHandler)
```

✅ **Estado:** El endpoint está configurado correctamente.

### 2. Handler Implementado

**Handler:** `handlers.RegistrarOficialHandler`

**Ubicación:** `handlers/rrhh.go` línea 99

**Respuesta Exitosa:**
```json
{
  "success": true,
  "message": "Oficial registrado correctamente",
  "oficial": {
    "id": "...",
    "credencial": "POL001",
    "qr_code": "data:image/png;base64,...",
    ...
  }
}
```

**Status Code:** 201 (Created)

✅ **Estado:** El handler está implementado y retorna el formato correcto.

### 3. CORS Configurado

**Ubicación:** `middleware/middleware.go` línea 32

**Configuración:**
```go
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

**Aplicado en:** `main.go` línea 45

```go
handler := middleware.CORSMiddleware(middleware.LoggingMiddleware(mux))
```

✅ **Estado:** CORS está configurado correctamente.

### 4. MongoDB Conectado

**Ubicación:** `main.go` línea 17

```go
if err := database.Connect(); err != nil {
	log.Fatalf("❌ Error al conectar a MongoDB: %v", err)
}
```

**Verificación:**
- El servidor debe mostrar: `🔌 Conectando a MongoDB...`
- Si hay error, se mostrará: `❌ Error al conectar a MongoDB: ...`
- Si conecta correctamente, no habrá error

✅ **Estado:** MongoDB se conecta al iniciar el servidor.

### 5. Logging de Peticiones

**Ubicación:** `middleware/middleware.go` línea 11

**Formato de logs:**
```
POST /api/rrhh/registrar-oficial 127.0.0.1:xxxxx 201 123.456ms
```

✅ **Estado:** Todas las peticiones se registran en los logs.

## 🔍 Cómo Verificar que Todo Funciona

### Paso 1: Iniciar el Backend

```bash
go run main.go
```

**Logs esperados:**
```
🔌 Conectando a MongoDB...
📦 Inicializando datos por defecto...
👤 Inicializando usuario admin...
🚀 Servidor GEPN iniciado en el puerto 8080
📍 Health check disponible en: http://localhost:8080/health
```

### Paso 2: Verificar Health Check

```bash
curl http://localhost:8080/health
```

**Respuesta esperada:**
```json
{"status": "ok"}
```

### Paso 3: Obtener Token de Master

Primero necesitas hacer login como master:

```bash
curl -X POST http://localhost:8080/api/master/login \
  -H "Content-Type: application/json" \
  -d '{
    "usuario": "admin",
    "contraseña": "Admin123!"
  }'
```

**Respuesta esperada:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "master": {
    "id": "...",
    "usuario": "admin",
    "permisos": ["rrhh", "policial", ...]
  },
  "mensaje": "Login exitoso"
}
```

**Nota:** Guarda el token para usarlo en el siguiente paso.

### Paso 4: Probar Registro de Oficial

```bash
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <TOKEN_OBTENIDO_EN_PASO_3>" \
  -d '{
    "primer_nombre": "Juan",
    "primer_apellido": "Pérez",
    "cedula": "V-12345678",
    "contraseña": "password123",
    "fecha_nacimiento": "1990-05-15",
    "estatura": 175.5,
    "color_piel": "Moreno",
    "tipo_sangre": "O+",
    "ciudad_nacimiento": "Caracas",
    "credencial": "POL-12345",
    "rango": "Oficial",
    "destacado": "",
    "fecha_graduacion": "2015-06-15",
    "estado": "Distrito Capital",
    "municipio": "Libertador",
    "parroquia": "Catedral",
    "foto_cara": "data:image/png;base64,iVBORw0KGgo..."
  }'
```

**Respuesta exitosa esperada:**
```json
{
  "success": true,
  "message": "Oficial registrado correctamente",
  "oficial": {
    "id": "...",
    "credencial": "POL-12345",
    "qr_code": "data:image/png;base64,...",
    "primer_nombre": "Juan",
    "primer_apellido": "Pérez",
    ...
  }
}
```

**Status Code:** 201

### Paso 5: Verificar Logs del Backend

En la consola del backend deberías ver:

```
POST /api/rrhh/registrar-oficial 127.0.0.1:xxxxx 201 123.456ms
```

## ❌ Posibles Errores y Soluciones

### Error 404 - Endpoint no encontrado

**Causa:** La ruta no está configurada o la URL es incorrecta.

**Solución:**
1. Verifica que el servidor esté corriendo
2. Verifica la URL: debe ser exactamente `/api/rrhh/registrar-oficial`
3. Verifica que el método sea `POST`

### Error 401 - No autorizado

**Causa:** No se envió el token o el token es inválido.

**Respuesta:**
```json
{
  "error": "Se requiere autenticación como usuario master para registrar oficiales"
}
```

**Solución:**
1. Haz login como master primero
2. Incluye el token en el header: `Authorization: <token>`
3. Verifica que el token no haya expirado

### Error 403 - Sin permisos

**Causa:** El usuario master no tiene el permiso "rrhh".

**Respuesta:**
```json
{
  "error": "No tiene permisos para acceder al módulo RRHH"
}
```

**Solución:**
1. Verifica que el usuario master tenga el permiso "rrhh" en su array de permisos
2. Puedes actualizar los permisos usando el endpoint `/api/master/usuarios/permisos/:id`

### Error 409 - Credencial o cédula duplicada

**Causa:** La credencial o cédula ya está registrada.

**Respuesta:**
```json
{
  "error": "La credencial ya está registrada"
}
```

o

```json
{
  "error": "La cédula ya está registrada"
}
```

**Solución:**
1. Usa una credencial o cédula diferente
2. O elimina el oficial existente de la base de datos

### Error 400 - Validación fallida

**Causas posibles:**
- Credencial vacía
- Cédula vacía
- Contraseña menor a 6 caracteres
- Rango inválido
- Fecha de graduación vacía

**Solución:**
1. Verifica que todos los campos obligatorios estén presentes
2. Verifica que el rango sea uno de los válidos (ver `OPCIONES_FRONTEND_RRHH.md`)
3. Verifica que la contraseña tenga al menos 6 caracteres

### Error 500 - Error interno del servidor

**Causa:** Error al conectar con MongoDB o error al crear el oficial.

**Solución:**
1. Verifica que MongoDB esté corriendo
2. Verifica la conexión a MongoDB en los logs del servidor
3. Revisa los logs del backend para ver el error específico

### Error CORS

**Causa:** El frontend está en un origen diferente y CORS no está configurado.

**Solución:**
1. Verifica que el middleware CORS esté aplicado (ya está configurado)
2. Verifica que el frontend esté enviando el header `Content-Type: application/json`
3. Si el error persiste, verifica que el servidor esté corriendo

## 📋 Checklist de Verificación

Antes de reportar un problema, verifica:

- [ ] El servidor está corriendo (`go run main.go`)
- [ ] MongoDB está conectado (ver logs del servidor)
- [ ] El endpoint existe: `/api/rrhh/registrar-oficial`
- [ ] El método es `POST`
- [ ] Se incluye el header `Authorization` con el token de master
- [ ] El token es válido y no ha expirado
- [ ] El usuario master tiene el permiso "rrhh"
- [ ] Todos los campos obligatorios están presentes
- [ ] La credencial y cédula son únicas
- [ ] El rango es válido
- [ ] La contraseña tiene al menos 6 caracteres
- [ ] CORS está configurado (ya está configurado en el código)

## 🔧 Comandos Útiles para Debug

### Ver logs en tiempo real

```bash
# En Windows PowerShell
Get-Content -Path "logs.txt" -Wait -ErrorAction SilentlyContinue

# O simplemente observa la consola donde corre el servidor
```

### Verificar conexión a MongoDB

```bash
# Verificar que MongoDB esté corriendo
mongosh
# O
mongo
```

### Probar con Postman

1. Método: `POST`
2. URL: `http://localhost:8080/api/rrhh/registrar-oficial`
3. Headers:
   - `Content-Type: application/json`
   - `Authorization: <token>`
4. Body (raw JSON):
```json
{
  "primer_nombre": "Juan",
  "primer_apellido": "Pérez",
  "cedula": "V-12345678",
  "contraseña": "password123",
  "fecha_nacimiento": "1990-05-15",
  "estatura": 175.5,
  "color_piel": "Moreno",
  "tipo_sangre": "O+",
  "ciudad_nacimiento": "Caracas",
  "credencial": "POL-12345",
  "rango": "Oficial",
  "destacado": "",
  "fecha_graduacion": "2015-06-15",
  "estado": "Distrito Capital",
  "municipio": "Libertador",
  "parroquia": "Catedral",
  "foto_cara": "data:image/png;base64,iVBORw0KGgo..."
}
```

## 📝 Información para Compartir si el Problema Persiste

Si el problema persiste después de verificar todo lo anterior, comparte:

1. **Logs de la consola del backend** (cuando intentas registrar)
2. **Status Code** de la respuesta (200, 201, 400, 401, 403, 409, 500, etc.)
3. **Respuesta completa del backend** (JSON completo)
4. **Logs de la consola del navegador** (F12 → Console)
5. **Request Headers** (F12 → Network → Headers)
6. **Request Payload** (F12 → Network → Payload)

Con esta información podremos identificar el problema exacto.

---

**Última actualización:** 2025-01-27
**Versión:** 1.0.0

