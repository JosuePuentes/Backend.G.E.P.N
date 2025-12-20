# GEPN - Backend API

Backend en Go para la aplicación móvil de la Policía Nacional.

## 🚀 Inicio Rápido

```bash
# Ejecutar el servidor
go run main.go

# El servidor estará disponible en http://localhost:8080
```

## 📡 Endpoints de la API

### Públicos

- `GET /` - Home
- `GET /health` - Health check
- `GET /ciudadano` - Portal ciudadano

### Autenticación

- `POST /api/policial/login` - Login de policiales
  ```json
  {
    "credencial": "POL001",
    "pin": "123456"
  }
  ```

### Protegidos (requieren token en header `Authorization`)

#### Detenidos
- `POST /api/detenidos` - Crear registro de detenido
- `GET /api/detenidos/listar` - Listar todos los detenidos
- `GET /api/detenidos/obtener?id=1` - Obtener detenido por ID

#### Minutas
- `POST /api/minutas` - Crear minuta digital
- `GET /api/minutas/listar` - Listar todas las minutas
- `GET /api/minutas/obtener?id=1` - Obtener minuta por ID

#### Búsqueda
- `POST /api/buscar/cedula` - Buscar cédula
- `GET /api/mas-buscados` - Listar los más buscados

#### Pánico
- `POST /api/panico/activar` - Activar botón de pánico
- `GET /api/panico/alertas` - Listar alertas de pánico

## 🔐 Usuarios de Prueba

- Credencial: `POL001`, PIN: `123456`
- Credencial: `POL002`, PIN: `654321`

## 📱 Instrucciones para el Frontend Móvil

### Tecnología Recomendada: React Native

Para crear una app que funcione en Android (APK) e iOS (iPhone 12-17), usa **React Native**.

### Pasos para crear el frontend:

1. **Instalar React Native CLI:**
```bash
npm install -g react-native-cli
```

2. **Crear el proyecto:**
```bash
npx react-native init GEPNApp --template react-native-template-typescript
cd GEPNApp
```

3. **Instalar dependencias necesarias:**
```bash
npm install @react-navigation/native @react-navigation/stack
npm install react-native-screens react-native-safe-area-context
npm install @react-native-community/geolocation
npm install axios
npm install @react-native-async-storage/async-storage
```

4. **Estructura de pantallas sugerida:**
```
src/
├── screens/
│   ├── HomeScreen.tsx          # Pantalla inicial con botón login
│   ├── CiudadanoScreen.tsx     # Pantalla /ciudadano
│   ├── LoginPolicialScreen.tsx # Login con credenciales y PIN
│   ├── DashboardScreen.tsx     # Dashboard con 4 botones
│   ├── DetenidosScreen.tsx     # Registro de detenidos
│   ├── MinutasScreen.tsx       # Minutas digitales
│   ├── BusquedaScreen.tsx      # Buscador de cédulas
│   └── MasBuscadosScreen.tsx   # Los más buscados
├── components/
│   └── PanicButton.tsx         # Botón de pánico
├── services/
│   └── api.ts                  # Cliente API
└── navigation/
    └── AppNavigator.tsx        # Navegación
```

### Características a implementar:

1. **Home Screen:**
   - Botón "Iniciar Sesión" que navega a `/ciudadano`
   - Ruta oculta `/policial` para login de policías

2. **Login Policial:**
   - Campo de credencial
   - Campo de PIN (6 dígitos, numérico)
   - Solicitar permisos de GPS al hacer login

3. **Dashboard (después del login):**
   - 4 botones/símbolos en grid:
     - Registro de Detenidos
     - Minutas Digitales
     - Buscador de Cédulas
     - Los Más Buscados
   - Botón de pánico rojo abajo en el centro
   - El botón de pánico requiere mantener presionado 5 segundos

4. **GPS:**
   - Solicitar permisos al iniciar sesión
   - Usar `@react-native-community/geolocation`
   - Enviar coordenadas en todas las peticiones que lo requieran

### Configuración para iOS (iPhone 12-17):

En `ios/GEPNApp/Info.plist` agregar:
```xml
<key>NSLocationWhenInUseUsageDescription</key>
<string>Necesitamos tu ubicación para los servicios policiales</string>
<key>NSLocationAlwaysUsageDescription</key>
<string>Necesitamos tu ubicación para los servicios policiales</string>
```

### Configuración para Android (APK):

En `android/app/src/main/AndroidManifest.xml` agregar:
```xml
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />
```

### Generar APK:

```bash
cd android
./gradlew assembleRelease
# El APK estará en: android/app/build/outputs/apk/release/app-release.apk
```

### Ejemplo de Cliente API (services/api.ts):

```typescript
import axios from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const API_URL = 'http://tu-backend-url.com'; // Cambiar por tu URL

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para agregar token
api.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('token');
  if (token) {
    config.headers.Authorization = token;
  }
  return config;
});

export const authService = {
  login: (credencial: string, pin: string) =>
    api.post('/api/policial/login', { credencial, pin }),
};

export const detenidosService = {
  crear: (data: any) => api.post('/api/detenidos', data),
  listar: () => api.get('/api/detenidos/listar'),
};

// ... más servicios
```

## 🌐 Variables de Entorno

El servidor usa la variable de entorno `PORT` (por defecto 8080).

Para producción (Render/Vercel):
```bash
export PORT=8080
```

## 📝 Notas

- El sistema de autenticación actual es básico (tokens en memoria)
- En producción, implementar JWT y base de datos
- Los datos se almacenan en memoria (se pierden al reiniciar)
- En producción, usar PostgreSQL, MySQL o MongoDB

