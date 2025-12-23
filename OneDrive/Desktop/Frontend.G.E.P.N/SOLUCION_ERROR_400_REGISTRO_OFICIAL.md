# Solución Error 400 - Registro de Oficial

## 🔍 Problema Identificado

**Error:** `Request failed with status code 400`

**Status Code:** 400 Bad Request

**Causa:** El backend está rechazando la petición porque falta algún campo obligatorio o hay un error de validación.

## ✅ Cambios Realizados en el Backend

He mejorado el handler para que ahora retorne mensajes de error más claros en formato JSON:

```json
{
  "success": false,
  "error": "Mensaje de error específico"
}
```

## 🔧 Cómo Identificar el Error Específico

### Paso 1: Ver la Respuesta del Backend

En el frontend, cuando recibas el error 400, necesitas capturar y mostrar la respuesta del backend:

```typescript
try {
  const response = await axios.post(
    'https://backend-g-e-p-n.onrender.com/api/rrhh/registrar-oficial',
    datosOficial,
    {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    }
  );
} catch (error) {
  // Capturar el error específico
  if (error.response) {
    // El servidor respondió con un error
    console.error('❌ Error del servidor:', error.response.data);
    console.error('❌ Status:', error.response.status);
    console.error('❌ Mensaje:', error.response.data.error || error.response.data.message);
    
    // Mostrar el error al usuario
    Alert.alert('Error', error.response.data.error || error.response.data.message || 'Error al registrar oficial');
  } else if (error.request) {
    // La petición se hizo pero no hubo respuesta
    console.error('❌ Sin respuesta del servidor');
    Alert.alert('Error', 'No se pudo conectar con el servidor');
  } else {
    // Error al configurar la petición
    console.error('❌ Error:', error.message);
    Alert.alert('Error', 'Error de configuración: ' + error.message);
  }
}
```

### Paso 2: Verificar los Logs del Backend

En los logs de Render, ahora verás mensajes más específicos:

```
❌ Validación fallida: Credencial vacía
❌ Validación fallida: Cédula vacía
❌ Validación fallida: Contraseña inválida (longitud: X)
❌ Validación fallida: Rango inválido: X
❌ Validación fallida: Fecha de graduación vacía
```

## 📋 Campos Obligatorios que Deben Estar Presentes

Asegúrate de que estos campos estén presentes y no estén vacíos:

1. ✅ **credencial** - No puede estar vacío
2. ✅ **cedula** - No puede estar vacío
3. ✅ **contraseña** - Mínimo 6 caracteres
4. ✅ **rango** - Debe ser uno de los rangos válidos
5. ✅ **fecha_graduacion** - No puede estar vacío (formato: YYYY-MM-DD)
6. ✅ **primer_nombre** - Recomendado
7. ✅ **primer_apellido** - Recomendado

## 🔍 Errores Comunes y Soluciones

### Error: "La credencial es obligatoria"

**Causa:** El campo `credencial` está vacío o no se está enviando.

**Solución:**
```typescript
// Verificar que credencial tenga valor
if (!credencial || credencial.trim() === '') {
  Alert.alert('Error', 'La credencial es obligatoria');
  return;
}
```

### Error: "La cédula es obligatoria"

**Causa:** El campo `cedula` está vacío o no se está enviando.

**Solución:**
```typescript
// Verificar que cedula tenga valor
if (!cedula || cedula.trim() === '') {
  Alert.alert('Error', 'La cédula es obligatoria');
  return;
}
```

### Error: "La contraseña debe tener al menos 6 caracteres"

**Causa:** El campo `contraseña` tiene menos de 6 caracteres.

**Solución:**
```typescript
// Verificar longitud de contraseña
if (!contraseña || contraseña.length < 6) {
  Alert.alert('Error', 'La contraseña debe tener al menos 6 caracteres');
  return;
}
```

### Error: "Rango inválido"

**Causa:** El campo `rango` no es uno de los rangos válidos.

**Rangos válidos:**
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
- Subcomisario
- Comisario General de Brigada
- Comisario General de División
- Comisario General Inspector
- Comisario General en Jefe

**Solución:**
```typescript
// Verificar que el rango sea válido
const rangosValidos = [
  'Oficial', 'Primer Oficial', 'Oficial Jefe',
  'Inspector', 'Primer Inspector', 'Inspector Jefe',
  'Comisario', 'Primer Comisario', 'Comisario Jefe',
  'Comisario General', 'Comisario Mayor', 'Comisario Superior',
  'Subcomisario', 'Comisario General de Brigada',
  'Comisario General de División', 'Comisario General Inspector',
  'Comisario General en Jefe'
];

if (!rango || !rangosValidos.includes(rango)) {
  Alert.alert('Error', 'Por favor seleccione un rango válido');
  return;
}
```

### Error: "La fecha de graduación es obligatoria"

**Causa:** El campo `fecha_graduacion` está vacío o no se está enviando.

**Solución:**
```typescript
// Verificar que fechaGraduacion tenga valor
if (!fechaGraduacion || fechaGraduacion.trim() === '') {
  Alert.alert('Error', 'La fecha de graduación es obligatoria');
  return;
}

// Verificar formato YYYY-MM-DD
const fechaRegex = /^\d{4}-\d{2}-\d{2}$/;
if (!fechaRegex.test(fechaGraduacion)) {
  Alert.alert('Error', 'La fecha de graduación debe estar en formato YYYY-MM-DD');
  return;
}
```

### Error: "Error al decodificar la petición"

**Causa:** El JSON enviado está mal formado o hay un problema con el formato de los datos.

**Solución:**
```typescript
// Verificar que los datos se estén serializando correctamente
const datosOficial = {
  primer_nombre: primerNombre || '',
  segundo_nombre: segundoNombre || '',
  primer_apellido: primerApellido || '',
  segundo_apellido: segundoApellido || '',
  cedula: cedula || '',
  contraseña: contraseña || '',
  fecha_nacimiento: fechaNacimiento || '',
  estatura: parseFloat(estatura) || 0,
  color_piel: colorPiel || '',
  tipo_sangre: tipoSangre || '',
  ciudad_nacimiento: ciudadNacimiento || '',
  credencial: credencial || '',
  rango: rango || '',
  destacado: destacado || '',
  fecha_graduacion: fechaGraduacion || '',
  estado: estado || '',
  municipio: municipio || '',
  parroquia: parroquia || '',
  foto_cara: fotoCara || '',
};

// Verificar antes de enviar
console.log('📤 Datos a enviar:', JSON.stringify(datosOficial, null, 2));
```

## 🧪 Código de Ejemplo Completo para el Frontend

```typescript
const handleRegistrarOficial = async () => {
  // Validar campos obligatorios antes de enviar
  if (!credencial || credencial.trim() === '') {
    Alert.alert('Error', 'La credencial es obligatoria');
    return;
  }

  if (!cedula || cedula.trim() === '') {
    Alert.alert('Error', 'La cédula es obligatoria');
    return;
  }

  if (!contraseña || contraseña.length < 6) {
    Alert.alert('Error', 'La contraseña debe tener al menos 6 caracteres');
    return;
  }

  if (!rango || rango.trim() === '') {
    Alert.alert('Error', 'El rango es obligatorio');
    return;
  }

  if (!fechaGraduacion || fechaGraduacion.trim() === '') {
    Alert.alert('Error', 'La fecha de graduación es obligatoria');
    return;
  }

  // Construir objeto de datos
  const datosOficial = {
    primer_nombre: primerNombre || '',
    segundo_nombre: segundoNombre || '',
    primer_apellido: primerApellido || '',
    segundo_apellido: segundoApellido || '',
    cedula: cedula.trim(),
    contraseña: contraseña,
    fecha_nacimiento: fechaNacimiento || '',
    estatura: estatura ? parseFloat(estatura) : 0,
    color_piel: colorPiel || '',
    tipo_sangre: tipoSangre || '',
    ciudad_nacimiento: ciudadNacimiento || '',
    credencial: credencial.trim(),
    rango: rango,
    destacado: destacado || '',
    fecha_graduacion: fechaGraduacion,
    estado: estado || '',
    municipio: municipio || '',
    parroquia: parroquia || '',
    foto_cara: fotoCara || '',
  };

  try {
    setLoading(true);
    
    const response = await axios.post(
      'https://backend-g-e-p-n.onrender.com/api/rrhh/registrar-oficial',
      datosOficial,
      {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token,
        },
      }
    );

    if (response.data.success) {
      Alert.alert('Éxito', 'Oficial registrado correctamente');
      // Limpiar formulario o navegar
    }
  } catch (error) {
    console.error('❌ Error completo:', error);
    
    if (error.response) {
      // El servidor respondió con un error
      const errorMessage = error.response.data?.error || 
                          error.response.data?.message || 
                          'Error al registrar oficial';
      
      console.error('❌ Error del servidor:', error.response.data);
      console.error('❌ Status:', error.response.status);
      
      Alert.alert('Error', errorMessage);
    } else if (error.request) {
      // La petición se hizo pero no hubo respuesta
      console.error('❌ Sin respuesta del servidor');
      Alert.alert('Error', 'No se pudo conectar con el servidor');
    } else {
      // Error al configurar la petición
      console.error('❌ Error:', error.message);
      Alert.alert('Error', 'Error: ' + error.message);
    }
  } finally {
    setLoading(false);
  }
};
```

## 📊 Verificar en los Logs de Render

Después de hacer el deploy, cuando intentes registrar un oficial, verás en los logs de Render mensajes como:

```
📝 Intento de registro - Credencial: POL001, Cédula: V-12345678, Rango: Oficial, FechaGraduacion: 2015-06-15
```

O si hay un error:

```
❌ Validación fallida: Credencial vacía
❌ Validación fallida: Rango inválido: X
```

## ✅ Próximos Pasos

1. **Actualizar el frontend** para capturar y mostrar el mensaje de error específico del backend
2. **Verificar los logs de Render** después de intentar registrar para ver qué validación está fallando
3. **Asegurarse de que todos los campos obligatorios** estén presentes antes de enviar

---

**Última actualización:** 2025-12-23
**Versión:** 1.0.0

