# Opciones para el Frontend - Módulo RRHH

Este documento contiene todas las opciones disponibles para los campos del formulario de registro de oficiales en el módulo RRHH.

## 1. Color de Piel

Las opciones válidas para el campo "Color de Piel" son:

- **Blanco**
- **Negro**
- **Moreno**
- **Trigueño**
- **Mestizo**
- **Amarillo** (para personas de origen asiático)
- **Otro**

**Nota:** No se debe incluir "Indígena" como opción de color de piel, ya que se refiere a etnia, no a color de piel.

## 2. Tipo de Sangre

Todos los tipos de sangre del sistema ABO con factor Rh:

### Tipo O
- **O+** (O Positivo)
- **O-** (O Negativo)

### Tipo A
- **A+** (A Positivo)
- **A-** (A Negativo)

### Tipo B
- **B+** (B Positivo)
- **B-** (B Negativo)

### Tipo AB
- **AB+** (AB Positivo)
- **AB-** (AB Negativo)

## 3. Ciudades de Nacimiento (Venezuela)

Lista de ciudades principales de Venezuela por estado:

### Amazonas
- Puerto Ayacucho

### Anzoátegui
- Barcelona
- Puerto La Cruz
- El Tigre
- Anaco
- Cantaura
- Puerto Píritu

### Apure
- San Fernando de Apure
- Guasdualito
- Elorza
- Achaguas

### Aragua
- Maracay
- Turmero
- La Victoria
- Cagua
- Villa de Cura
- El Consejo
- San Mateo

### Barinas
- Barinas
- Ciudad Bolivia
- Socopó
- Barinitas

### Bolívar
- Ciudad Bolívar
- Ciudad Guayana
- Upata
- El Callao
- Tumeremo
- Caicara del Orinoco

### Carabobo
- Valencia
- Puerto Cabello
- Guacara
- Mariara
- San Joaquín
- Bejuma
- Morón

### Cojedes
- San Carlos
- Tinaquillo
- El Baúl
- Las Vegas

### Delta Amacuro
- Tucupita
- Curiapo

### Distrito Capital
- Caracas

### Falcón
- Coro
- Punto Fijo
- La Vela
- Churuguara
- Dabajuro

### Guárico
- San Juan de los Morros
- Calabozo
- Valle de la Pascua
- Zaraza
- Altagracia de Orituco

### Lara
- Barquisimeto
- Carora
- Duaca
- El Tocuyo
- Quíbor
- Cabudare

### Mérida
- Mérida
- El Vigía
- Ejido
- Tovar
- Bailadores

### Miranda
- Los Teques
- Guarenas
- Guatire
- Santa Teresa del Tuy
- Ocumare del Tuy
- Charallave
- Cúa
- San Antonio de los Altos

### Monagas
- Maturín
- Punta de Mata
- Caripito
- Caripe

### Nueva Esparta
- La Asunción
- Porlamar
- Juan Griego
- Pampatar

### Portuguesa
- Guanare
- Acarigua
- Araure
- Turén
- Ospino

### Sucre
- Cumaná
- Carúpano
- Güiria
- Irapa
- Araya

### Táchira
- San Cristóbal
- Táriba
- La Fría
- Rubio
- Colón

### Trujillo
- Valera
- Trujillo
- Boconó
- Betijoque
- Carache

### Vargas
- La Guaira
- Maiquetía
- Catia La Mar
- Caraballeda
- Macuto

### Yaracuy
- San Felipe
- Yaritagua
- Chivacoa
- Nirgua
- Aroa

### Zulia
- Maracaibo
- Cabimas
- Ciudad Ojeda
- San Francisco
- La Villa del Rosario
- Machiques
- Santa Bárbara del Zulia

## 4. Rangos

Todos los rangos válidos del sistema (en orden jerárquico):

1. **Oficial**
2. **Primer Oficial**
3. **Oficial Jefe**
4. **Inspector**
5. **Primer Inspector**
6. **Inspector Jefe**
7. **Comisario**
8. **Primer Comisario**
9. **Comisario Jefe**
10. **Comisario General**
11. **Comisario Mayor**
12. **Comisario Superior**
13. **Subcomisario**
14. **Comisario General de Brigada**
15. **Comisario General de División**
16. **Comisario General Inspector**
17. **Comisario General en Jefe**

## 5. Campo Destacado

El campo "Destacado" es **opcional** y debe dejarse **vacío** al registrar un oficial en RRHH. Este campo se asignará posteriormente en otros módulos del sistema (como el módulo de Centro de Coordinación).

**Validación Frontend:**
- El campo debe permitir estar vacío
- No debe ser obligatorio
- No debe tener valor por defecto

## 6. Validaciones Importantes

### Credencial
- **Obligatorio**
- **Único** en el sistema
- Si ya está registrada, mostrar mensaje: **"La credencial ya está registrada"**

### Cédula
- **Obligatorio**
- **Única** en el sistema
- Si ya está registrada, mostrar mensaje: **"La cédula ya está registrada"**

### Contraseña
- **Obligatorio**
- Mínimo 6 caracteres
- Esta contraseña será la que se use para el login en el módulo policial
- El campo en el login policial puede llamarse "pin" o "contraseña", pero debe contener la contraseña registrada en RRHH

## 7. Ejemplo de Request para Registrar Oficial

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
  "destacado": "",
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

## 8. Mensajes de Error del Backend

### Credencial Duplicada
```json
{
  "error": "La credencial ya está registrada"
}
```
**Status Code:** 409 (Conflict)

### Cédula Duplicada
```json
{
  "error": "La cédula ya está registrada"
}
```
**Status Code:** 409 (Conflict)

### Rango Inválido
```json
{
  "error": "Rango inválido"
}
```
**Status Code:** 400 (Bad Request)

## 9. Campos de Fecha - Date Picker (Calendario)

### Campos que Requieren Date Picker

1. **Fecha de Nacimiento** (`fecha_nacimiento`)
   - Campo obligatorio
   - Formato: `YYYY-MM-DD` (ejemplo: `1990-01-15`)
   - Usar calendario para selección

2. **Fecha de Graduación** (`fecha_graduacion`)
   - Campo obligatorio
   - Formato: `YYYY-MM-DD` (ejemplo: `2015-06-15`)
   - Usar calendario para selección

3. **Fechas de Parientes** (opcionales)
   - Fecha de nacimiento del padre
   - Fecha de nacimiento de la madre
   - Fechas de nacimiento de los hijos
   - Formato: `YYYY-MM-DD`

### Validaciones de Fechas

- La fecha de nacimiento no puede ser mayor a la fecha actual
- La fecha de graduación no puede ser mayor a la fecha actual
- La fecha de graduación debe ser al menos 18 años después de la fecha de nacimiento

**Ver:** `INSTRUCCIONES_FRONTEND_DATE_PICKER_RRHH.md` para implementación completa del date picker.

## 10. Notas para el Frontend

1. **Login Policial:** El campo de contraseña en el login policial debe aceptar la contraseña registrada en RRHH. El backend acepta tanto el campo "pin" como "contraseña" en el request de login.

2. **Color de Piel:** No incluir "Indígena" como opción. Usar las opciones listadas arriba.

3. **Date Pickers:** Implementar calendarios para todos los campos de fecha. Ver `INSTRUCCIONES_FRONTEND_DATE_PICKER_RRHH.md` para detalles.

3. **Tipo de Sangre:** Mostrar todos los 8 tipos de sangre disponibles.

4. **Ciudad de Nacimiento:** Implementar un selector con todas las ciudades listadas, agrupadas por estado si es posible.

5. **Rangos:** Mostrar todos los rangos listados, no solo hasta "Inspector".

6. **Destacado:** Dejar el campo vacío por defecto y no hacerlo obligatorio.

---

**Última actualización:** 2025-01-27
**Versión:** 1.1.0

