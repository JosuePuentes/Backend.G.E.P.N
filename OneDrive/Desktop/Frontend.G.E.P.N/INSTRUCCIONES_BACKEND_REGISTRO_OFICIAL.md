# Instrucciones de Debug - Registro de Oficial en RRHH

Este documento contiene instrucciones para verificar y depurar el endpoint de registro de oficiales.

## 🔍 Problema Reportado

**Síntoma:** Al hacer click en "Registrar Oficial" en el módulo RRHH, no pasa nada (solo "pestaña").

**Posibles causas:**
1. El frontend no está enviando la petición
2. Error en el frontend que no se muestra
3. El endpoint no está recibiendo la petición
4. Error de autenticación
5. Error de validación

## ✅ Verificación del Backend

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

**Requisitos:**
- Método: `POST`
- Header: `Authorization: <token_master>`
- Content-Type: `application/json`
- El usuario master debe tener permiso "rrhh"

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

### 3. Validaciones del Handler

El handler valida:
1. ✅ Método HTTP debe ser POST
2. ✅ Token de autorización presente
3. ✅ Usuario master activo
4. ✅ Permiso "rrhh" en el usuario master
5. ✅ Credencial no vacía
6. ✅ Cédula no vacía
7. ✅ Contraseña mínimo 6 caracteres
8. ✅ Rango válido
9. ✅ Fecha de graduación no vacía
10. ✅ Credencial única
11. ✅ Cédula única

## 🔧 Cómo Verificar que el Backend Funciona

### Paso 1: Verificar que el Servidor Esté Corriendo

En los logs de Render deberías ver:
```
🚀 Servidor GEPN iniciado en el puerto 10000
📍 Health check disponible en: http://localhost:10000/health
```

### Paso 2: Probar el Endpoint con curl

**Nota:** Reemplaza `<TOKEN>` con un token válido de master.

```bash
curl -X POST https://backend-g-e-p-n.onrender.com/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <TOKEN>" \
  -d '{
    "primer_nombre": "Juan",
    "segundo_nombre": "Carlos",
    "primer_apellido": "Pérez",
    "segundo_apellido": "González",
    "cedula": "V-12345678",
    "contraseña": "password123",
    "fecha_nacimiento": "1990-05-15",
    "estatura": 175.5,
    "color_piel": "Moreno",
    "tipo_sangre": "O+",
    "ciudad_nacimiento": "Caracas",
    "credencial": "POL-TEST-001",
    "rango": "Oficial",
    "destacado": "",
    "fecha_graduacion": "2015-06-15",
    "estado": "Distrito Capital",
    "municipio": "Libertador",
    "parroquia": "Catedral",
    "foto_cara": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
  }'
```

### Paso 3: Verificar los Logs del Backend

Cuando se hace una petición, deberías ver en los logs:

```
POST /api/rrhh/registrar-oficial <IP> 201 <tiempo>
```

O si hay un error:

```
POST /api/rrhh/registrar-oficial <IP> 400 <tiempo>
POST /api/rrhh/registrar-oficial <IP> 401 <tiempo>
POST /api/rrhh/registrar-oficial <IP> 403 <tiempo>
POST /api/rrhh/registrar-oficial <IP> 409 <tiempo>
POST /api/rrhh/registrar-oficial <IP> 500 <tiempo>
```

## ❌ Errores Comunes y Soluciones

### Error 401 - No autorizado

**Causa:** No se envió el token o el token es inválido.

**Respuesta:**
```json
{
  "error": "Se requiere autenticación como usuario master para registrar oficiales"
}
```

**Solución:**
1. Verificar que el frontend esté enviando el header `Authorization`
2. Verificar que el token sea válido
3. Hacer login nuevamente para obtener un token fresco

### Error 403 - Sin permisos

**Causa:** El usuario master no tiene el permiso "rrhh".

**Respuesta:**
```json
{
  "error": "No tiene permisos para acceder al módulo RRHH"
}
```

**Solución:**
1. Verificar que el usuario master tenga el permiso "rrhh" en su array de permisos
2. Actualizar los permisos usando el endpoint `/api/master/usuarios/permisos/:id`

### Error 400 - Validación fallida

**Causas posibles:**
- Credencial vacía
- Cédula vacía
- Contraseña menor a 6 caracteres
- Rango inválido
- Fecha de graduación vacía

**Solución:**
1. Verificar que todos los campos obligatorios estén presentes
2. Verificar que el rango sea válido
3. Verificar que la contraseña tenga al menos 6 caracteres

### Error 409 - Credencial o cédula duplicada

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
1. Usar una credencial o cédula diferente
2. O eliminar el oficial existente de la base de datos

### Error 500 - Error interno

**Causa:** Error al conectar con MongoDB o error al crear el oficial.

**Solución:**
1. Verificar que MongoDB esté conectado (ver logs del servidor)
2. Revisar los logs del backend para ver el error específico
3. Verificar que la base de datos esté accesible

## 📋 Checklist para el Frontend

El frontend debe verificar:

- [ ] El botón "Registrar Oficial" tiene un `onPress` o `onClick` configurado
- [ ] Se está capturando el evento correctamente
- [ ] Se está construyendo el objeto con todos los campos requeridos
- [ ] Se está enviando la petición POST a `/api/rrhh/registrar-oficial`
- [ ] Se está incluyendo el header `Authorization` con el token
- [ ] Se está incluyendo el header `Content-Type: application/json`
- [ ] Se está manejando la respuesta (éxito o error)
- [ ] Se están mostrando mensajes de error al usuario
- [ ] Se está validando el formulario antes de enviar

## 🔍 Debug del Frontend

### Verificar en la Consola del Navegador (F12)

1. **Console Tab:**
   - Buscar errores de JavaScript
   - Buscar errores de red
   - Verificar que la función de registro se esté llamando

2. **Network Tab:**
   - Filtrar por "registrar-oficial"
   - Verificar que se esté haciendo la petición POST
   - Ver el Status Code de la respuesta
   - Ver los Headers enviados
   - Ver el Payload (body) enviado
   - Ver la Response recibida

### Ejemplo de Código Frontend Correcto

```typescript
const registrarOficial = async () => {
  // Validar formulario
  if (!validarFormulario()) {
    Alert.alert('Error', 'Por favor complete todos los campos requeridos');
    return;
  }

  // Construir objeto de datos
  const datosOficial = {
    primer_nombre: primerNombre,
    segundo_nombre: segundoNombre,
    primer_apellido: primerApellido,
    segundo_apellido: segundoApellido,
    cedula: cedula,
    contraseña: contraseña,
    fecha_nacimiento: fechaNacimiento, // Formato: YYYY-MM-DD
    estatura: parseFloat(estatura),
    color_piel: colorPiel,
    tipo_sangre: tipoSangre,
    ciudad_nacimiento: ciudadNacimiento,
    credencial: credencial,
    rango: rango,
    destacado: destacado || "",
    fecha_graduacion: fechaGraduacion, // Formato: YYYY-MM-DD
    estado: estado,
    municipio: municipio,
    parroquia: parroquia,
    foto_cara: fotoCara, // Base64
  };

  try {
    setLoading(true);
    
    const response = await fetch('https://backend-g-e-p-n.onrender.com/api/rrhh/registrar-oficial', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token, // Token del master
      },
      body: JSON.stringify(datosOficial),
    });

    const data = await response.json();

    if (response.ok && data.success) {
      Alert.alert('Éxito', 'Oficial registrado correctamente');
      // Limpiar formulario o navegar
      limpiarFormulario();
    } else {
      Alert.alert('Error', data.error || 'Error al registrar oficial');
    }
  } catch (error) {
    console.error('Error:', error);
    Alert.alert('Error', 'Error de conexión. Verifique su internet.');
  } finally {
    setLoading(false);
  }
};
```

## 📊 Verificar en la Base de Datos

### Conectar a MongoDB y Verificar

```javascript
// En MongoDB Compass o mongo shell
use gepn

// Ver todos los oficiales
db.oficiales.find().pretty()

// Contar oficiales
db.oficiales.countDocuments()

// Buscar un oficial específico
db.oficiales.findOne({ credencial: "POL-TEST-001" })

// Verificar índices
db.oficiales.getIndexes()
```

## 🚀 Instrucciones para Verificar el Backend

### 1. Verificar que el Endpoint Esté Activo

```bash
# Health check
curl https://backend-g-e-p-n.onrender.com/health

# Debe retornar: {"status":"ok"}
```

### 2. Obtener Token de Master

```bash
curl -X POST https://backend-g-e-p-n.onrender.com/api/master/login \
  -H "Content-Type: application/json" \
  -d '{
    "usuario": "admin",
    "contraseña": "Admin123!"
  }'
```

**Guardar el token** de la respuesta.

### 3. Probar Registro de Oficial

```bash
curl -X POST https://backend-g-e-p-n.onrender.com/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <TOKEN_OBTENIDO>" \
  -d '{
    "primer_nombre": "Test",
    "primer_apellido": "Usuario",
    "cedula": "V-TEST-001",
    "contraseña": "test123",
    "fecha_nacimiento": "1990-01-01",
    "estatura": 175,
    "color_piel": "Moreno",
    "tipo_sangre": "O+",
    "ciudad_nacimiento": "Caracas",
    "credencial": "POL-TEST-001",
    "rango": "Oficial",
    "destacado": "",
    "fecha_graduacion": "2015-01-01",
    "estado": "Distrito Capital",
    "municipio": "Libertador",
    "parroquia": "Catedral",
    "foto_cara": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
  }'
```

### 4. Verificar en los Logs

En los logs de Render deberías ver:
```
POST /api/rrhh/registrar-oficial <IP> 201 <tiempo>
```

## 📝 Información para Compartir si el Problema Persiste

Si el problema persiste, comparte:

1. **Logs de la consola del navegador** (F12 → Console)
2. **Petición en Network Tab** (F12 → Network → buscar "registrar-oficial")
   - Status Code
   - Request Headers
   - Request Payload
   - Response
3. **Logs del backend en Render** cuando intentas registrar
4. **Código del botón/componente** que maneja el registro

Con esta información podremos identificar el problema exacto.

---

**Última actualización:** 2025-01-27
**Versión:** 1.0.0

