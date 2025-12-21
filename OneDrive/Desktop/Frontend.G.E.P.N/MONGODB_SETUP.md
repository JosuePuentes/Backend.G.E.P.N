# Configuración de MongoDB

## Variables de Entorno

Para conectar tu backend con MongoDB, necesitas configurar las siguientes variables de entorno:

### En Render.com

1. Ve a tu servicio en Render
2. Navega a la sección "Environment"
3. Agrega las siguientes variables:

```
MONGODB_URI=mongodb+srv://Drocolven2019:TU_CONTRASEÑA@drocolven2019.eof9ilx.mongodb.net/?appName=Drocolven2019
MONGODB_DB_NAME=gepn
```

**Importante:** Reemplaza `TU_CONTRASEÑA` con la contraseña real de tu base de datos MongoDB.

### En Local (Desarrollo)

Crea un archivo `.env` en la raíz del proyecto (no lo subas a Git):

```
MONGODB_URI=mongodb+srv://Drocolven2019:TU_CONTRASEÑA@drocolven2019.eof9ilx.mongodb.net/?appName=Drocolven2019
MONGODB_DB_NAME=gepn
PORT=8080
```

## Colecciones de MongoDB

El sistema creará automáticamente las siguientes colecciones:

- `usuarios` - Usuarios del sistema (policiales)
- `detenidos` - Registros de detenidos
- `minutas` - Minutas digitales
- `busquedas` - Historial de búsquedas de cédulas
- `mas_buscados` - Lista de personas más buscadas
- `panico` - Alertas de pánico

## Datos Iniciales

El sistema inicializa automáticamente:

1. **Usuarios por defecto:**
   - Credencial: `POL001`, PIN: `123456`
   - Credencial: `POL002`, PIN: `654321`

2. **Más buscados de ejemplo:**
   - Cédula: `1234567890`
   - Cédula: `0987654321`

## Verificación

Para verificar que la conexión funciona:

1. Inicia el servidor: `go run main.go`
2. Deberías ver: `✅ Conectado a MongoDB exitosamente`
3. Prueba el endpoint: `GET /health`

## Solución de Problemas

### Error: "MONGODB_URI no está configurada"
- Asegúrate de configurar la variable de entorno `MONGODB_URI` en Render

### Error: "Authentication failed"
- Verifica que la contraseña en la URI sea correcta
- Asegúrate de que el usuario tenga permisos en la base de datos

### Error: "Connection timeout"
- Verifica que la IP de Render esté en la whitelist de MongoDB Atlas
- En MongoDB Atlas, ve a "Network Access" y agrega `0.0.0.0/0` para permitir todas las IPs (solo para desarrollo)


