# Instrucciones de Actualización - Backend RRHH

Este documento contiene las instrucciones para aplicar los cambios realizados en el módulo RRHH.

## Cambios Realizados

### 1. Login Policial - Uso de Contraseña de RRHH

**Cambio:** El login policial ahora usa la contraseña registrada en RRHH.

**Archivos modificados:**
- `handlers/auth.go`: Actualizado para aceptar contraseña en campos "pin" o "contraseña"
- `models/models.go`: Agregado campo "contraseña" al LoginRequest (además de "pin" para compatibilidad)

**Comportamiento:**
- El login acepta la contraseña en el campo `pin` o `contraseña`
- La contraseña debe ser la misma registrada en RRHH para el oficial
- Se verifica con bcrypt contra la contraseña hasheada almacenada

### 2. Mensajes de Error Mejorados para Credenciales Duplicadas

**Cambio:** Los mensajes de error ahora retornan JSON con mensaje claro.

**Archivos modificados:**
- `handlers/rrhh.go`: Cambiado `http.Error` por respuesta JSON con status 409

**Mensajes:**
- Credencial duplicada: `{"error": "La credencial ya está registrada"}` (Status 409)
- Cédula duplicada: `{"error": "La cédula ya está registrada"}` (Status 409)

### 3. Rangos Adicionales

**Cambio:** Agregados más rangos válidos al sistema.

**Archivos modificados:**
- `handlers/rrhh.go`: Agregados rangos adicionales a `rangosValidos`

**Rangos agregados:**
- Subcomisario
- Comisario General de Brigada
- Comisario General de División
- Comisario General Inspector
- Comisario General en Jefe

**Rangos completos (17 total):**
1. Oficial
2. Primer Oficial
3. Oficial Jefe
4. Inspector
5. Primer Inspector
6. Inspector Jefe
7. Comisario
8. Primer Comisario
9. Comisario Jefe
10. Comisario General
11. Comisario Mayor
12. Comisario Superior
13. Subcomisario
14. Comisario General de Brigada
15. Comisario General de División
16. Comisario General Inspector
17. Comisario General en Jefe

### 4. Campo Destacado Opcional

**Cambio:** El campo "destacado" ahora es explícitamente opcional y se deja vacío por defecto.

**Archivos modificados:**
- `handlers/rrhh.go`: Agregada validación para permitir destacado vacío

**Comportamiento:**
- El campo puede estar vacío al registrar un oficial
- Se asignará posteriormente en otros módulos (Centro de Coordinación)
- No se valida como obligatorio

## Instrucciones para Aplicar los Cambios

### Paso 1: Verificar Archivos Modificados

Los siguientes archivos han sido modificados:
- `handlers/auth.go`
- `handlers/rrhh.go`
- `models/models.go`

### Paso 2: Compilar y Probar

```bash
# Compilar el proyecto
go build

# O ejecutar directamente
go run main.go
```

### Paso 3: Probar Endpoints

#### Probar Registro de Oficial con Credencial Duplicada

```bash
# Primera vez (debe funcionar)
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_master>" \
  -d '{
    "credencial": "POL001",
    "cedula": "12345678",
    "contraseña": "password123",
    "primer_nombre": "Juan",
    "primer_apellido": "Pérez",
    "rango": "Oficial",
    "fecha_graduacion": "2015-06-15",
    "foto_cara": "base64..."
  }'

# Segunda vez con misma credencial (debe retornar error 409)
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_master>" \
  -d '{
    "credencial": "POL001",
    "cedula": "87654321",
    "contraseña": "password456",
    "primer_nombre": "María",
    "primer_apellido": "González",
    "rango": "Oficial",
    "fecha_graduacion": "2016-06-15",
    "foto_cara": "base64..."
  }'
```

**Respuesta esperada (409):**
```json
{
  "error": "La credencial ya está registrada"
}
```

#### Probar Login Policial con Contraseña de RRHH

```bash
# Login con la contraseña registrada en RRHH
curl -X POST http://localhost:8080/api/policial/login \
  -H "Content-Type: application/json" \
  -d '{
    "credencial": "POL001",
    "pin": "password123",
    "latitud": 10.4969,
    "longitud": -66.8983
  }'
```

#### Probar Registro con Rango Superior

```bash
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_master>" \
  -d '{
    "credencial": "POL002",
    "cedula": "87654321",
    "contraseña": "password123",
    "primer_nombre": "Carlos",
    "primer_apellido": "Rodríguez",
    "rango": "Comisario General en Jefe",
    "fecha_graduacion": "2000-06-15",
    "foto_cara": "base64..."
  }'
```

#### Probar Registro con Destacado Vacío

```bash
curl -X POST http://localhost:8080/api/rrhh/registrar-oficial \
  -H "Content-Type: application/json" \
  -H "Authorization: <token_master>" \
  -d '{
    "credencial": "POL003",
    "cedula": "11223344",
    "contraseña": "password123",
    "primer_nombre": "Ana",
    "primer_apellido": "Martínez",
    "rango": "Inspector",
    "fecha_graduacion": "2018-06-15",
    "destacado": "",
    "foto_cara": "base64..."
  }'
```

## Validaciones del Frontend

El frontend debe implementar las siguientes validaciones y opciones:

1. **Color de Piel:** Ver `OPCIONES_FRONTEND_RRHH.md` para lista completa
   - No incluir "Indígena"
   - Opciones: Blanco, Negro, Moreno, Trigueño, Mestizo, Amarillo, Otro

2. **Tipo de Sangre:** Ver `OPCIONES_FRONTEND_RRHH.md` para lista completa
   - Todos los 8 tipos: O+, O-, A+, A-, B+, B-, AB+, AB-

3. **Ciudad de Nacimiento:** Ver `OPCIONES_FRONTEND_RRHH.md` para lista completa
   - Lista de ciudades de Venezuela por estado

4. **Rangos:** Ver `OPCIONES_FRONTEND_RRHH.md` para lista completa
   - Todos los 17 rangos disponibles

5. **Destacado:** Campo opcional, dejar vacío por defecto

6. **Mensajes de Error:**
   - Credencial duplicada: Mostrar "La credencial ya está registrada"
   - Cédula duplicada: Mostrar "La cédula ya está registrada"

## Documentación Adicional

Se ha creado el archivo `OPCIONES_FRONTEND_RRHH.md` con todas las opciones disponibles para el frontend:
- Colores de piel válidos
- Tipos de sangre completos
- Ciudades de Venezuela
- Rangos completos
- Validaciones y mensajes de error

## Notas Importantes

1. **Compatibilidad:** El login policial mantiene compatibilidad con el campo "pin" pero también acepta "contraseña"

2. **Seguridad:** Las contraseñas se siguen hasheando con bcrypt antes de almacenarse

3. **Validación:** Los rangos se validan contra la lista completa de rangos válidos

4. **Destacado:** Este campo se asignará en otros módulos, no en RRHH

5. **Mensajes de Error:** Todos los errores de credenciales duplicadas ahora retornan JSON con status 409

---

**Fecha de actualización:** 2025-01-27
**Versión:** 1.1.0

