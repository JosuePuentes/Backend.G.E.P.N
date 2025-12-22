# Instrucciones para Date Picker (Calendario) - Módulo RRHH

Este documento contiene las instrucciones para implementar calendarios (date pickers) en los campos de fecha del formulario de registro de oficiales.

## 📅 Campos que Requieren Date Picker

### Campos Principales del Oficial

1. **Fecha de Nacimiento** (`fecha_nacimiento`)
   - Campo obligatorio
   - Formato esperado: `YYYY-MM-DD` (ejemplo: `1990-01-15`)

2. **Fecha de Graduación** (`fecha_graduacion`)
   - Campo obligatorio
   - Formato esperado: `YYYY-MM-DD` (ejemplo: `2015-06-15`)

### Campos de Parientes (Opcionales)

3. **Fecha de Nacimiento del Padre** (`parientes.padre.fecha_nacimiento`)
   - Campo opcional
   - Formato esperado: `YYYY-MM-DD`

4. **Fecha de Nacimiento de la Madre** (`parientes.madre.fecha_nacimiento`)
   - Campo opcional
   - Formato esperado: `YYYY-MM-DD`

5. **Fecha de Nacimiento de los Hijos** (`parientes.hijos[].fecha_nacimiento`)
   - Campo opcional (puede haber múltiples hijos)
   - Formato esperado: `YYYY-MM-DD`

## 🔧 Implementación para React Native

### Opción 1: Usar @react-native-community/datetimepicker (Recomendado)

#### Instalación

```bash
npm install @react-native-community/datetimepicker
# Para iOS
cd ios && pod install && cd ..
```

#### Ejemplo de Implementación

```typescript
import React, { useState } from 'react';
import { View, Text, TouchableOpacity, Platform } from 'react-native';
import DateTimePicker from '@react-native-community/datetimepicker';

interface DatePickerProps {
  label: string;
  value: string; // Formato: YYYY-MM-DD
  onChange: (date: string) => void;
  maximumDate?: Date;
  minimumDate?: Date;
  required?: boolean;
}

const DatePickerField: React.FC<DatePickerProps> = ({
  label,
  value,
  onChange,
  maximumDate,
  minimumDate,
  required = false,
}) => {
  const [show, setShow] = useState(false);
  const [date, setDate] = useState<Date>(
    value ? new Date(value) : new Date()
  );

  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const formatDisplayDate = (dateString: string): string => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-VE', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const onDateChange = (event: any, selectedDate?: Date) => {
    if (Platform.OS === 'android') {
      setShow(false);
    }
    
    if (selectedDate) {
      setDate(selectedDate);
      onChange(formatDate(selectedDate));
    }
  };

  return (
    <View style={{ marginBottom: 16 }}>
      <Text style={{ marginBottom: 8, fontWeight: 'bold' }}>
        {label} {required && <Text style={{ color: 'red' }}>*</Text>}
      </Text>
      
      <TouchableOpacity
        onPress={() => setShow(true)}
        style={{
          borderWidth: 1,
          borderColor: '#ccc',
          borderRadius: 8,
          padding: 12,
          backgroundColor: '#fff',
        }}
      >
        <Text style={{ color: value ? '#000' : '#999' }}>
          {value ? formatDisplayDate(value) : `Seleccionar ${label.toLowerCase()}`}
        </Text>
      </TouchableOpacity>

      {show && (
        <DateTimePicker
          value={date}
          mode="date"
          display={Platform.OS === 'ios' ? 'spinner' : 'default'}
          onChange={onDateChange}
          maximumDate={maximumDate}
          minimumDate={minimumDate}
          locale="es-VE"
        />
      )}

      {Platform.OS === 'ios' && show && (
        <View style={{ flexDirection: 'row', justifyContent: 'flex-end', marginTop: 8 }}>
          <TouchableOpacity
            onPress={() => setShow(false)}
            style={{ padding: 10, marginRight: 10 }}
          >
            <Text style={{ color: '#007AFF' }}>Cancelar</Text>
          </TouchableOpacity>
          <TouchableOpacity
            onPress={() => {
              onChange(formatDate(date));
              setShow(false);
            }}
            style={{ padding: 10 }}
          >
            <Text style={{ color: '#007AFF', fontWeight: 'bold' }}>Confirmar</Text>
          </TouchableOpacity>
        </View>
      )}
    </View>
  );
};

export default DatePickerField;
```

#### Uso en el Formulario

```typescript
import DatePickerField from './components/DatePickerField';

const RegistrarOficialScreen = () => {
  const [fechaNacimiento, setFechaNacimiento] = useState('');
  const [fechaGraduacion, setFechaGraduacion] = useState('');
  const [fechaNacimientoPadre, setFechaNacimientoPadre] = useState('');

  // Fecha máxima: hoy (para fecha de nacimiento)
  const fechaMaxima = new Date();
  
  // Fecha mínima: hace 100 años (para fecha de nacimiento)
  const fechaMinimaNacimiento = new Date();
  fechaMinimaNacimiento.setFullYear(fechaMinimaNacimiento.getFullYear() - 100);

  // Fecha mínima para graduación: hace 50 años
  const fechaMinimaGraduacion = new Date();
  fechaMinimaGraduacion.setFullYear(fechaMinimaGraduacion.getFullYear() - 50);

  return (
    <View>
      {/* Fecha de Nacimiento */}
      <DatePickerField
        label="Fecha de Nacimiento"
        value={fechaNacimiento}
        onChange={setFechaNacimiento}
        maximumDate={fechaMaxima}
        minimumDate={fechaMinimaNacimiento}
        required={true}
      />

      {/* Fecha de Graduación */}
      <DatePickerField
        label="Fecha de Graduación"
        value={fechaGraduacion}
        onChange={setFechaGraduacion}
        maximumDate={fechaMaxima}
        minimumDate={fechaMinimaGraduacion}
        required={true}
      />

      {/* Fecha de Nacimiento del Padre (Opcional) */}
      <DatePickerField
        label="Fecha de Nacimiento del Padre"
        value={fechaNacimientoPadre}
        onChange={setFechaNacimientoPadre}
        maximumDate={fechaMaxima}
        minimumDate={fechaMinimaNacimiento}
        required={false}
      />
    </View>
  );
};
```

### Opción 2: Usar react-native-date-picker

#### Instalación

```bash
npm install react-native-date-picker
cd ios && pod install && cd ..
```

#### Ejemplo de Implementación

```typescript
import React, { useState } from 'react';
import { View, Text, TouchableOpacity } from 'react-native';
import DatePicker from 'react-native-date-picker';

const DatePickerField = ({ label, value, onChange, required = false }) => {
  const [open, setOpen] = useState(false);
  const [date, setDate] = useState(value ? new Date(value) : new Date());

  const formatDate = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  const formatDisplayDate = (dateString: string): string => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-VE', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  return (
    <View style={{ marginBottom: 16 }}>
      <Text style={{ marginBottom: 8, fontWeight: 'bold' }}>
        {label} {required && <Text style={{ color: 'red' }}>*</Text>}
      </Text>
      
      <TouchableOpacity
        onPress={() => setOpen(true)}
        style={{
          borderWidth: 1,
          borderColor: '#ccc',
          borderRadius: 8,
          padding: 12,
          backgroundColor: '#fff',
        }}
      >
        <Text style={{ color: value ? '#000' : '#999' }}>
          {value ? formatDisplayDate(value) : `Seleccionar ${label.toLowerCase()}`}
        </Text>
      </TouchableOpacity>

      <DatePicker
        modal
        open={open}
        date={date}
        mode="date"
        locale="es"
        onConfirm={(selectedDate) => {
          setOpen(false);
          setDate(selectedDate);
          onChange(formatDate(selectedDate));
        }}
        onCancel={() => {
          setOpen(false);
        }}
      />
    </View>
  );
};
```

## 📋 Validaciones Importantes

### Fecha de Nacimiento
- **Obligatorio**
- **Formato:** `YYYY-MM-DD`
- **Validación:** No puede ser mayor a la fecha actual
- **Límite razonable:** No más de 100 años en el pasado

### Fecha de Graduación
- **Obligatorio**
- **Formato:** `YYYY-MM-DD`
- **Validación:** 
  - No puede ser mayor a la fecha actual
  - Debe ser posterior a la fecha de nacimiento (al menos 18 años después)
  - Límite razonable: No más de 50 años en el pasado

### Fechas de Parientes
- **Opcionales**
- **Formato:** `YYYY-MM-DD`
- **Validación:** No pueden ser mayores a la fecha actual

## 🔍 Función de Validación

```typescript
const validarFechas = (
  fechaNacimiento: string,
  fechaGraduacion: string
): { valido: boolean; error?: string } => {
  if (!fechaNacimiento || !fechaGraduacion) {
    return { valido: false, error: 'Las fechas son obligatorias' };
  }

  const fechaNac = new Date(fechaNacimiento);
  const fechaGrad = new Date(fechaGraduacion);
  const hoy = new Date();

  // Validar que no sean futuras
  if (fechaNac > hoy) {
    return { valido: false, error: 'La fecha de nacimiento no puede ser futura' };
  }

  if (fechaGrad > hoy) {
    return { valido: false, error: 'La fecha de graduación no puede ser futura' };
  }

  // Validar que la graduación sea posterior al nacimiento
  const edadMinima = 18; // Edad mínima para graduarse
  const añosDiferencia = (fechaGrad.getTime() - fechaNac.getTime()) / (1000 * 60 * 60 * 24 * 365.25);
  
  if (añosDiferencia < edadMinima) {
    return { 
      valido: false, 
      error: `La fecha de graduación debe ser al menos ${edadMinima} años después de la fecha de nacimiento` 
    };
  }

  return { valido: true };
};
```

## 📤 Formato para Enviar al Backend

El backend espera el formato `YYYY-MM-DD` (ISO 8601). Asegúrate de convertir la fecha seleccionada a este formato antes de enviarla:

```typescript
const enviarDatos = async () => {
  const datos = {
    primer_nombre: primerNombre,
    fecha_nacimiento: fechaNacimiento, // Formato: "1990-01-15"
    fecha_graduacion: fechaGraduacion, // Formato: "2015-06-15"
    // ... otros campos
  };

  // Validar antes de enviar
  const validacion = validarFechas(fechaNacimiento, fechaGraduacion);
  if (!validacion.valido) {
    Alert.alert('Error', validacion.error);
    return;
  }

  // Enviar al backend
  try {
    const response = await fetch('http://localhost:8080/api/rrhh/registrar-oficial', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
      body: JSON.stringify(datos),
    });
    // ... manejar respuesta
  } catch (error) {
    console.error('Error:', error);
  }
};
```

## 🎨 Estilos Recomendados

```typescript
const styles = StyleSheet.create({
  datePickerContainer: {
    marginBottom: 16,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
    marginBottom: 8,
    color: '#333',
  },
  datePickerButton: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    backgroundColor: '#fff',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  datePickerText: {
    fontSize: 16,
    color: '#333',
  },
  datePickerPlaceholder: {
    fontSize: 16,
    color: '#999',
  },
  required: {
    color: '#FF0000',
  },
});
```

## ✅ Checklist de Implementación

- [ ] Instalar la librería de date picker
- [ ] Crear componente DatePickerField reutilizable
- [ ] Implementar en campo "Fecha de Nacimiento"
- [ ] Implementar en campo "Fecha de Graduación"
- [ ] Implementar en campos de parientes (opcional)
- [ ] Agregar validaciones de fechas
- [ ] Formatear fechas a `YYYY-MM-DD` antes de enviar
- [ ] Probar en iOS y Android
- [ ] Agregar indicador visual de campo requerido
- [ ] Agregar mensajes de error de validación

## 📝 Notas Importantes

1. **Formato de Fecha:** El backend siempre espera `YYYY-MM-DD`. No uses otros formatos.

2. **Zona Horaria:** Asegúrate de manejar correctamente las zonas horarias. Usa UTC o la zona horaria local según corresponda.

3. **Validación en Frontend:** Valida las fechas antes de enviarlas al backend para mejorar la experiencia del usuario.

4. **Localización:** Configura el date picker para usar español (es-VE) para mejor experiencia del usuario.

5. **Plataformas:** Prueba en ambas plataformas (iOS y Android) ya que el comportamiento puede diferir.

---

**Última actualización:** 2025-01-27
**Versión:** 1.0.0

