# Problema: Contraseña con más de 6 caracteres rechazada

## 🔍 Problema Identificado

**Error:** `La contraseña debe tener al menos 6 caracteres`

**Contraseña enviada:** `"123456789a"` (10 caracteres)

**Status Code:** 400 Bad Request

## 📊 Análisis del Problema

Según los logs del frontend:
- La contraseña que se envía es: `"123456789a"` (10 caracteres)
- El backend rechaza con: "La contraseña debe tener al menos 6 caracteres"

Esto sugiere que:
1. El campo `contraseña` no está llegando correctamente al backend
2. El campo está llegando vacío o con menos de 6 caracteres
3. Hay un problema con el mapeo del JSON

## 🔧 Solución Implementada

He agregado logs detallados en el backend para ver exactamente qué está recibiendo:

```go
log.Printf("🔐 Contraseña recibida - Longitud: %d, Valor: [%s]", len(oficial.Contraseña), oficial.Contraseña)
```

## 📋 Pasos para Verificar

### 1. Verificar los Logs de Render

Después del deploy, cuando intentes registrar un oficial, verás en los logs de Render:

```
🔐 Contraseña recibida - Longitud: X, Valor: [valor recibido]
```

Esto te dirá:
- Si el campo está llegando vacío
- Si está llegando con menos caracteres de los esperados
- Si hay caracteres especiales que se están perdiendo

### 2. Verificar el Payload en el Frontend

En la consola del navegador (F12 → Network), verifica:
- **Request Payload**: ¿El campo `contraseña` está presente?
- **Valor exacto**: ¿Cuál es el valor exacto que se está enviando?

### 3. Posibles Causas

#### Causa 1: El campo se está perdiendo en el JSON

**Solución:** Verificar que el campo se llame exactamente `contraseña` (con la ñ):

```typescript
const datosOficial = {
  // ... otros campos
  contraseña: contraseña, // ← Verificar que el nombre del campo sea exacto
  // ...
};
```

#### Causa 2: El campo está siendo truncado

**Solución:** Verificar que no haya validaciones en el frontend que estén truncando la contraseña antes de enviarla.

#### Causa 3: Caracteres especiales

**Solución:** Si la contraseña tiene caracteres especiales, verificar que se estén codificando correctamente en el JSON.

## 🔍 Debug en el Frontend

Agrega este log justo antes de enviar la petición:

```typescript
console.log('🔐 Contraseña a enviar:', contraseña);
console.log('🔐 Longitud de contraseña:', contraseña.length);
console.log('🔐 Contraseña en objeto:', datosOficial.contraseña);
console.log('🔐 Longitud en objeto:', datosOficial.contraseña?.length);
console.log('📦 Objeto completo:', JSON.stringify(datosOficial));
```

## ✅ Verificación en el Backend

Después del deploy, los logs mostrarán:

```
📝 Intento de registro - Credencial: 24241240a, Cédula: 24241240, Rango: Oficial, FechaGraduacion: 2020-02-20
🔐 Contraseña recibida - Longitud: X, Valor: [valor]
```

Si la longitud es 0 o menor a 6, el problema está en cómo se está enviando desde el frontend.

## 🚀 Próximos Pasos

1. **Esperar el deploy** del backend con los nuevos logs
2. **Intentar registrar** un oficial nuevamente
3. **Revisar los logs de Render** para ver qué longitud tiene la contraseña recibida
4. **Comparar** con lo que se está enviando desde el frontend

## 📝 Nota Importante

El campo en el modelo de Go es `Contraseña` (con mayúscula y ñ), pero en JSON debe ser `contraseña` (minúscula). Go maneja esto automáticamente con las etiquetas `json:"contraseña"`, pero es importante verificar que el frontend esté usando el nombre correcto.

---

**Última actualización:** 2025-12-23
**Versión:** 1.0.0


