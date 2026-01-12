# 🚀 INSTRUCCIONES PARA DESPLEGAR EL BACKEND EN PRODUCCIÓN

## 📋 Objetivo
Desplegar el backend de GEPN (Go + MongoDB) en un servidor con HTTPS para que la app móvil pueda conectarse y descargarse.

## 🎯 Stack Tecnológico
- **Backend**: Go (Golang) 1.21
- **Base de Datos**: MongoDB Atlas (Free Tier)
- **Hosting**: Render.com (Free Tier con HTTPS automático)
- **Autenticación**: JWT ya implementado
- **CORS**: Ya configurado

---

## ✅ OPCIÓN 1: RENDER.COM (RECOMENDADO - GRATIS)

### Paso 1: Crear cuenta en Render
1. Ve a https://render.com
2. Registrate con GitHub (recomendado)
3. Conecta tu repositorio de GitHub

### Paso 2: Configurar MongoDB Atlas (Base de datos GRATUITA)

#### 2.1. Crear cuenta en MongoDB Atlas
1. Ve a: https://www.mongodb.com/cloud/atlas/register
2. Regístrate con tu email o Google
3. Completa el formulario de registro

#### 2.2. Crear cluster gratuito
1. Selecciona **"Free Tier"** (M0 Sandbox)
2. **Provider**: AWS
3. **Region**: Elige la más cercana (ejemplo: US East 1 - Virginia)
4. **Cluster Name**: `gepn-cluster` (o el nombre que prefieras)
5. Click en **"Create Cluster"**
6. Espera 3-5 minutos mientras se crea el cluster

#### 2.3. Crear usuario de base de datos
1. En el menú izquierdo: **"Database Access"**
2. Click en **"Add New Database User"**
3. Configuración:
   - **Authentication Method**: Password
   - **Username**: `gepn_user`
   - **Password**: Genera una contraseña segura (Guárdala, la necesitarás)
   - **Database User Privileges**: Select "Read and write to any database"
4. Click en **"Add User"**

#### 2.4. Permitir acceso desde cualquier IP
1. En el menú izquierdo: **"Network Access"**
2. Click en **"Add IP Address"**
3. Click en **"Allow Access from Anywhere"**
4. IP Address: `0.0.0.0/0` (se llena automáticamente)
5. Click en **"Confirm"**

⚠️ **IMPORTANTE**: En producción real deberías restringir las IPs, pero para desarrollo y despliegue en Render necesitamos acceso desde cualquier IP.

#### 2.5. Obtener Connection String
1. Regresa a **"Database"** (menú izquierdo)
2. En tu cluster, click en **"Connect"**
3. Selecciona **"Connect your application"**
4. **Driver**: Go, **Version**: 1.13 or later
5. Copia el **Connection String**, se ve así:
   ```
   mongodb+srv://gepn_user:<password>@gepn-cluster.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```
6. **REEMPLAZA** `<password>` con la contraseña que creaste en el paso 2.3
7. **AGREGA** el nombre de la base de datos al final:
   ```
   mongodb+srv://gepn_user:TU_PASSWORD_AQUI@gepn-cluster.xxxxx.mongodb.net/gepn?retryWrites=true&w=majority
   ```

💾 **Guarda este Connection String**, lo necesitarás en el siguiente paso.

### Paso 3: Desplegar en Render

#### 3.1. Crear cuenta en Render
1. Ve a: https://render.com
2. Click en **"Get Started"**
3. **Sign Up with GitHub** (recomendado para conectar tu repositorio)
4. Autoriza Render para acceder a tus repositorios

#### 3.2. Crear Web Service
1. En el Dashboard de Render, click en **"New +"**
2. Selecciona **"Web Service"**
3. Conecta tu repositorio:
   - Si no aparece tu repositorio, click en **"Configure account"** 
   - Autoriza acceso al repositorio del backend Go
4. Selecciona el repositorio del backend GEPN

#### 3.3. Configurar el servicio
Llena el formulario con estos datos:

**Basic Settings:**
- **Name**: `gepn-backend` (o el nombre que prefieras)
- **Region**: Oregon (US West) o el más cercano
- **Branch**: `main` (o la rama principal de tu proyecto)
- **Root Directory**: Déjalo vacío (a menos que el backend esté en una subcarpeta)

**Build Settings:**
- **Runtime**: Selecciona **"Docker"** 
  (Tu proyecto tiene Dockerfile, Render lo detectará automáticamente)
- **Build Command**: Se usa automáticamente el Dockerfile
- **Start Command**: Se usa automáticamente el Dockerfile

**Instance Settings:**
- **Instance Type**: Selecciona **"Free"**

#### 3.4. Agregar Variables de Entorno
Scroll hacia abajo hasta **"Environment Variables"** y agrega estas variables:

| Key | Value |
|-----|-------|
| `MONGODB_URI` | `mongodb+srv://gepn_user:TU_PASSWORD@gepn-cluster.xxxxx.mongodb.net/gepn?retryWrites=true&w=majority` |
| `MONGODB_DB_NAME` | `gepn` |
| `JWT_SECRET` | `tu_secreto_super_seguro_cambiar_esto_12345` |
| `PORT` | `8080` |
| `GO_ENV` | `production` |

⚠️ **MUY IMPORTANTE**: 
- Reemplaza `TU_PASSWORD` con la contraseña de MongoDB del Paso 2.3
- Reemplaza `JWT_SECRET` con un valor completamente aleatorio y seguro (mínimo 32 caracteres)
- Ejemplo de JWT_SECRET seguro: `a8f3k9m2p5q7w1e4r6t8y0u3i5o7p9s1d3f5`

#### 3.5. Crear el servicio
1. Click en **"Create Web Service"**
2. Render comenzará a construir y desplegar tu aplicación
3. Verás los logs en tiempo real

#### 3.6. Esperar el despliegue
- El primer despliegue toma **5-10 minutos**
- Verás mensajes como:
  - `Building...`
  - `Pushing...`
  - `Starting service...`
- Cuando veas `✅ Live` significa que está funcionando
- Tu URL será algo como: `https://gepn-backend.onrender.com`

📝 **Copia esta URL**, la necesitarás para la app móvil.

### Paso 4: Esperar el despliegue
- El despliegue tarda 5-10 minutos
- Verás los logs en tiempo real
- Cuando termine, te dará una URL como:
  ```
  https://gepn-backend.onrender.com
  ```

### Paso 4: Verificar que el backend funciona correctamente

#### 4.1. Health Check (Verificar que el servidor está vivo)

**En tu navegador**, abre:
```
https://tu-backend.onrender.com/health
```

**Respuesta esperada:**
```json
{
  "status": "healthy",
  "message": "GEPN Backend is running"
}
```

✅ Si ves esto, el servidor está funcionando correctamente.

#### 4.2. Probar Registro de Ciudadano

**Usando curl** (Terminal/PowerShell):
```bash
curl -X POST https://tu-backend.onrender.com/api/ciudadano/registro \
  -H "Content-Type: application/json" \
  -d "{\"nombre\":\"Usuario Prueba\",\"cedula\":\"V-99999999\",\"telefono\":\"0412-9999999\",\"password\":\"test123\"}"
```

**Usando Postman:**
- Método: `POST`
- URL: `https://tu-backend.onrender.com/api/ciudadano/registro`
- Headers: `Content-Type: application/json`
- Body (JSON):
```json
{
  "nombre": "Usuario Prueba",
  "cedula": "V-99999999",
  "telefono": "0412-9999999",
  "password": "test123"
}
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "ciudadano": {
    "id": "...",
    "nombre": "Usuario Prueba",
    "cedula": "V-99999999",
    "telefono": "0412-9999999"
  }
}
```

#### 4.3. Probar Login de Ciudadano

**Usando curl:**
```bash
curl -X POST https://tu-backend.onrender.com/api/ciudadano/login \
  -H "Content-Type: application/json" \
  -d "{\"cedula\":\"V-99999999\",\"password\":\"test123\"}"
```

**Usando Postman:**
- Método: `POST`
- URL: `https://tu-backend.onrender.com/api/ciudadano/login`
- Headers: `Content-Type: application/json`
- Body (JSON):
```json
{
  "cedula": "V-99999999",
  "password": "test123"
}
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "ciudadano": {
    "id": "...",
    "nombre": "Usuario Prueba",
    "cedula": "V-99999999",
    "telefono": "0412-9999999"
  }
}
```

#### 4.4. Probar Crear Denuncia (Requiere autenticación)

Primero, copia el `token` que obtuviste en el login anterior.

**Usando curl:**
```bash
curl -X POST https://tu-backend.onrender.com/api/denuncia/crear \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN_AQUI" \
  -d "{\"motivo\":\"Robo\",\"descripcion\":\"Me robaron el celular\",\"ubicacion\":\"Caracas\"}"
```

**Usando Postman:**
- Método: `POST`
- URL: `https://tu-backend.onrender.com/api/denuncia/crear`
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer TU_TOKEN_AQUI`
- Body (JSON):
```json
{
  "motivo": "Robo",
  "descripcion": "Me robaron el celular",
  "ubicacion": "Caracas"
}
```

**Respuesta esperada (201 Created):**
```json
{
  "message": "Denuncia creada exitosamente",
  "denuncia": {
    "id": "...",
    "motivo": "Robo",
    "descripcion": "Me robaron el celular",
    "ubicacion": "Caracas",
    "fecha": "2026-01-12T..."
  }
}
```

#### 4.5. Probar Login Policial

**Usando curl:**
```bash
curl -X POST https://tu-backend.onrender.com/api/policial/login \
  -H "Content-Type: application/json" \
  -d "{\"credencial\":\"POL001\",\"pin\":\"123456\"}"
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "policial": {
    "credencial": "POL001",
    "nombre": "...",
    "rango": "..."
  }
}
```

#### 4.6. Probar CORS desde el navegador

Abre la **Consola del Navegador** (F12) y ejecuta:

```javascript
fetch('https://tu-backend.onrender.com/health')
  .then(res => res.json())
  .then(data => console.log('✅ CORS funciona:', data))
  .catch(err => console.error('❌ Error CORS:', err));
```

✅ Si ves el mensaje de éxito, CORS está configurado correctamente.
❌ Si ves un error de CORS, revisa el middleware en `middleware/middleware.go`

---

## ✅ OPCIÓN 2: RAILWAY.APP (ALTERNATIVA)

Railway es otra excelente opción, similar a Render pero con un enfoque diferente.

### Paso 1: Configurar MongoDB Atlas
**Importante**: Railway ya no ofrece MongoDB gratuito, así que necesitas usar MongoDB Atlas igual que en la Opción 1 (ver Paso 2 arriba).

### Paso 2: Crear cuenta en Railway
1. Ve a https://railway.app
2. Click en **"Start a New Project"**
3. **Login with GitHub** (recomendado)
4. Verifica tu cuenta (requiere verificación, pero no te cobran)

### Paso 3: Desplegar el Backend
1. En Railway, click en **"New Project"**
2. Selecciona **"Deploy from GitHub repo"**
3. Conecta y autoriza tu repositorio
4. Selecciona el repositorio del backend GEPN
5. Railway detectará automáticamente que es Go con Docker

### Paso 4: Configurar Variables de Entorno
1. Click en tu servicio
2. Ve a la pestaña **"Variables"**
3. Agrega estas variables:

```
MONGODB_URI=mongodb+srv://gepn_user:TU_PASSWORD@gepn-cluster.xxxxx.mongodb.net/gepn?retryWrites=true&w=majority
MONGODB_DB_NAME=gepn
JWT_SECRET=tu_secreto_super_seguro_12345
PORT=8080
GO_ENV=production
```

### Paso 5: Generar dominio público
1. Ve a la pestaña **"Settings"**
2. Scroll hasta **"Networking"**
3. Click en **"Generate Domain"**
4. Railway te dará una URL como:
   ```
   https://gepn-backend-production.up.railway.app
   ```

### Paso 6: Verificar
Abre en tu navegador:
```
https://tu-url-railway.app/health
```

### Costos Railway
- **Trial**: $5 de crédito gratis al mes
- **Developer**: $5/mes de suscripción + uso
- Estimado para este proyecto: ~$5-10/mes

---

## ✅ OPCIÓN 3: GOOGLE CLOUD RUN (PROFESIONAL)

Google Cloud Run es ideal para aplicaciones containerizadas con Docker (como tu proyecto).

### Requisitos previos
- Cuenta de Google Cloud (300 USD de crédito gratis para nuevos usuarios)
- gcloud CLI instalado: https://cloud.google.com/sdk/docs/install

### Paso 1: Configurar MongoDB Atlas
Usa MongoDB Atlas (igual que Opción 1, ver Paso 2 arriba).

### Paso 2: Configurar Google Cloud

```bash
# 1. Autenticarse en Google Cloud
gcloud auth login

# 2. Crear proyecto (o usar uno existente)
gcloud projects create gepn-backend-2026 --name="GEPN Backend"

# 3. Configurar proyecto activo
gcloud config set project gepn-backend-2026

# 4. Habilitar APIs necesarias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com

# 5. Configurar región por defecto
gcloud config set run/region us-central1
```

### Paso 3: Desplegar desde el código fuente

Desde el directorio del proyecto:

```bash
# Cloud Run construirá la imagen Docker automáticamente
gcloud run deploy gepn-backend \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars "MONGODB_URI=mongodb+srv://gepn_user:PASSWORD@cluster.mongodb.net/gepn,MONGODB_DB_NAME=gepn,JWT_SECRET=tu_secreto_seguro,GO_ENV=production"
```

### Paso 4: Obtener la URL
Después del despliegue, Cloud Run te dará una URL como:
```
https://gepn-backend-xxxxx-uc.a.run.app
```

### Paso 5: Configurar dominio personalizado (Opcional)
```bash
# Mapear un dominio personalizado
gcloud run domain-mappings create \
  --service gepn-backend \
  --domain api.tudominio.com \
  --region us-central1
```

### Costos Google Cloud Run
- **Free Tier**: 
  - 2 millones de requests/mes gratis
  - 360,000 GB-segundos gratis
- **Después del Free Tier**: Pay-per-use
- **Estimado para este proyecto**: $0-5/mes (dentro del free tier normalmente)

---

## 📱 CONFIGURAR LA APP MÓVIL

Una vez que tengas el backend desplegado con su URL pública con HTTPS, necesitas actualizar la configuración en tu app móvil.

### Para React Native:

Busca el archivo `src/config/api.ts`, `src/services/api.ts` o similar:

```typescript
// ❌ ANTES (desarrollo local):
const API_BASE_URL = 'http://localhost:8080';

// ✅ DESPUÉS (producción):
const API_BASE_URL = 'https://gepn-backend.onrender.com'; // O tu URL

// O mejor aún, usar variable de entorno:
const API_BASE_URL = __DEV__ 
  ? 'http://localhost:8080'  // Desarrollo
  : 'https://gepn-backend.onrender.com';  // Producción

export default API_BASE_URL;
```

**Crear el cliente API:**

```typescript
import axios from 'axios';
import API_BASE_URL from '../config/api';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Interceptor para agregar el token JWT
apiClient.interceptors.request.use(
  async (config) => {
    const token = await AsyncStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default apiClient;
```

### Para Flutter:

Busca el archivo `lib/config/api_config.dart` o similar:

```dart
// ❌ ANTES (desarrollo local):
class ApiConfig {
  static const String baseUrl = 'http://localhost:8080';
}

// ✅ DESPUÉS (producción):
class ApiConfig {
  static const String baseUrl = 'https://gepn-backend.onrender.com';
  
  // O con detección de modo:
  static const bool isProduction = bool.fromEnvironment('dart.vm.product');
  static String get baseUrl => isProduction 
    ? 'https://gepn-backend.onrender.com'  // Producción
    : 'http://localhost:8080';  // Desarrollo
}
```

### Para Expo (React Native):

Crea un archivo `app.config.js`:

```javascript
export default {
  name: 'GEPN App',
  extra: {
    apiUrl: process.env.API_URL || 'https://gepn-backend.onrender.com',
  },
};
```

Luego en tu código:

```typescript
import Constants from 'expo-constants';

const API_BASE_URL = Constants.expoConfig?.extra?.apiUrl;
```

### Endpoints disponibles:

Una vez configurado, tu app podrá usar estos endpoints:

**Autenticación:**
- `POST /api/ciudadano/registro` - Registro de ciudadano
- `POST /api/ciudadano/login` - Login de ciudadano
- `POST /api/policial/login` - Login de policial
- `POST /api/master/login` - Login de master

**Denuncias (requieren token):**
- `POST /api/denuncia/crear` - Crear denuncia
- `GET /api/denuncia/mis-denuncias` - Mis denuncias
- `GET /api/denuncia/obtener?id=123` - Detalle de denuncia

**Policiales (requieren token):**
- `POST /api/detenidos` - Crear registro de detenido
- `GET /api/detenidos/listar` - Listar detenidos
- `POST /api/minutas` - Crear minuta
- `GET /api/minutas/listar` - Listar minutas
- `POST /api/buscar/cedula` - Buscar cédula
- `POST /api/panico/activar` - Activar botón de pánico

**Otros:**
- `GET /health` - Health check

---

## 🔍 VERIFICACIÓN COMPLETA DEL BACKEND

Antes de entregar la URL al equipo de la app móvil, verifica que todo funcione correctamente.

### 1. ✅ Health Check
```bash
curl https://tu-backend.onrender.com/health
```

**Respuesta esperada (200 OK):**
```json
{
  "status": "healthy",
  "message": "GEPN Backend is running"
}
```

### 2. ✅ Registro de Ciudadano
```bash
curl -X POST https://tu-backend.onrender.com/api/ciudadano/registro \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Test User","cedula":"V-88888888","telefono":"0412-8888888","password":"test123"}'
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGc...",
  "ciudadano": {
    "id": "...",
    "nombre": "Test User",
    "cedula": "V-88888888"
  }
}
```

### 3. ✅ Login de Ciudadano
```bash
curl -X POST https://tu-backend.onrender.com/api/ciudadano/login \
  -H "Content-Type: application/json" \
  -d '{"cedula":"V-88888888","password":"test123"}'
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGc...",
  "ciudadano": {...}
}
```

### 4. ✅ Login Policial
```bash
curl -X POST https://tu-backend.onrender.com/api/policial/login \
  -H "Content-Type: application/json" \
  -d '{"credencial":"POL001","pin":"123456"}'
```

**Respuesta esperada (200 OK):**
```json
{
  "token": "eyJhbGc...",
  "policial": {
    "credencial": "POL001",
    "nombre": "...",
    "rango": "..."
  }
}
```

### 5. ✅ Verificar CORS desde navegador

Abre la **Consola del Navegador** (F12 → Console) y ejecuta:

```javascript
fetch('https://tu-backend.onrender.com/health')
  .then(res => res.json())
  .then(data => console.log('✅ CORS OK:', data))
  .catch(err => console.error('❌ CORS Error:', err));
```

Si ves `✅ CORS OK`, el CORS está configurado correctamente.

### 6. ✅ Verificar HTTPS

La URL **DEBE** empezar con `https://` (no `http://`).

```bash
# ✅ Correcto
https://gepn-backend.onrender.com

# ❌ Incorrecto (iOS no funcionará)
http://gepn-backend.onrender.com
```

### 7. ✅ Verificar MongoDB

Revisa los logs de Render para ver si la conexión a MongoDB fue exitosa:

En el dashboard de Render → Logs, deberías ver:
```
🔌 Conectando a MongoDB...
✅ Conectado a MongoDB exitosamente
📦 Inicializando datos por defecto...
👤 Inicializando usuario admin...
🚀 Servidor GEPN iniciado en el puerto 8080
```

### 8. ✅ Verificar desde Postman

Importa esta colección en Postman para probar todos los endpoints:

**Collection: GEPN Backend Tests**

1. **Health Check**
   - GET `{{baseUrl}}/health`

2. **Registro Ciudadano**
   - POST `{{baseUrl}}/api/ciudadano/registro`
   - Body: `{"nombre":"...", "cedula":"...", "telefono":"...", "password":"..."}`

3. **Login Ciudadano**
   - POST `{{baseUrl}}/api/ciudadano/login`
   - Body: `{"cedula":"...", "password":"..."}`

4. **Login Policial**
   - POST `{{baseUrl}}/api/policial/login`
   - Body: `{"credencial":"POL001", "pin":"123456"}`

**Variable de Postman:**
- `baseUrl`: `https://tu-backend.onrender.com`

---

## ✅ CHECKLIST FINAL

Antes de entregar al equipo de la app móvil, verifica:

### Infraestructura
- [ ] **MongoDB Atlas** configurado y funcionando
- [ ] **Usuario de base de datos** creado con permisos correctos
- [ ] **Network Access** configurado (0.0.0.0/0 permitido)
- [ ] **Connection String** correcto y guardado

### Despliegue
- [ ] **Backend desplegado** en Render/Railway/Cloud Run
- [ ] **URL pública** disponible (ejemplo: `https://gepn-backend.onrender.com`)
- [ ] **HTTPS activo** (la URL debe empezar con `https://`)
- [ ] **Variables de entorno** configuradas:
  - `MONGODB_URI`
  - `MONGODB_DB_NAME`
  - `JWT_SECRET`
  - `PORT`
  - `GO_ENV`

### Funcionalidad
- [ ] **Health check** funcionando: `GET /health`
- [ ] **Registro de ciudadano** funcionando: `POST /api/ciudadano/registro`
- [ ] **Login de ciudadano** funcionando: `POST /api/ciudadano/login`
- [ ] **Login policial** funcionando: `POST /api/policial/login`
- [ ] **Crear denuncia** funcionando: `POST /api/denuncia/crear` (con token)
- [ ] **CORS configurado** correctamente (sin errores desde navegador)

### Base de Datos
- [ ] **Conexión a MongoDB** exitosa (revisar logs)
- [ ] **Datos iniciales** creados automáticamente
- [ ] **Usuario admin** creado (master/admin)
- [ ] **Colecciones** creadas correctamente

### Seguridad
- [ ] **JWT_SECRET** configurado con valor seguro y aleatorio
- [ ] **Contraseñas** hasheadas en la base de datos
- [ ] **HTTPS** activo (obligatorio para iOS)
- [ ] **CORS** permite peticiones de apps móviles

### Logs y Monitoreo
- [ ] **Logs de Render** muestran servidor iniciado correctamente
- [ ] **No hay errores** en los logs
- [ ] **Conexión a MongoDB** exitosa en los logs
- [ ] **Inicialización completada** sin errores

---

## 📞 INFORMACIÓN PARA ENTREGAR AL EQUIPO DE LA APP MÓVIL

Una vez completados todos los pasos y verificaciones, proporciona esta información al equipo de desarrollo de la app móvil:

---

### ✅ URL del Backend

```
https://gepn-backend.onrender.com
```
*(Reemplaza con tu URL real)*

---

### ✅ Estado del Backend

- **Servidor**: ✅ Corriendo
- **Base de datos**: ✅ MongoDB Atlas conectada
- **HTTPS**: ✅ Activo (obligatorio para iOS)
- **CORS**: ✅ Configurado para apps móviles
- **Última actualización**: [Fecha]

---

### ✅ Endpoints Disponibles

#### **Públicos (no requieren autenticación):**

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/health` | Health check del servidor |
| GET | `/` | Página principal |
| GET | `/ciudadano` | Portal ciudadano |
| POST | `/api/ciudadano/registro` | Registro de nuevo ciudadano |
| POST | `/api/ciudadano/login` | Login de ciudadano |
| POST | `/api/policial/login` | Login de policial |
| POST | `/api/master/login` | Login de master |

#### **Protegidos (requieren token JWT en header `Authorization: Bearer TOKEN`):**

**Denuncias (Ciudadanos):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/denuncia/crear` | Crear nueva denuncia |
| GET | `/api/denuncia/mis-denuncias` | Obtener mis denuncias |
| GET | `/api/denuncia/obtener?id=X` | Detalle de una denuncia |

**Detenidos (Policiales):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/detenidos` | Registrar detenido |
| GET | `/api/detenidos/listar` | Listar todos los detenidos |
| GET | `/api/detenidos/obtener?id=X` | Obtener detalle de detenido |

**Minutas (Policiales):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/minutas` | Crear minuta digital |
| GET | `/api/minutas/listar` | Listar todas las minutas |
| GET | `/api/minutas/obtener?id=X` | Obtener detalle de minuta |

**Búsqueda (Policiales):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/buscar/cedula` | Buscar persona por cédula |
| GET | `/api/mas-buscados` | Listar los más buscados |

**Pánico (Policiales):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/panico/activar` | Activar botón de pánico |
| GET | `/api/panico/alertas` | Listar alertas de pánico |

**RRHH (Master con permiso):**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/rrhh/registrar-oficial` | Registrar nuevo oficial |
| GET | `/api/rrhh/listar-oficiales` | Listar todos los oficiales |
| GET | `/api/rrhh/generar-qr/:credencial` | Generar QR para oficial |

---

### ✅ Formato de Autenticación

**Header requerido para endpoints protegidos:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Ejemplo en JavaScript/TypeScript:**
```typescript
const token = await AsyncStorage.getItem('token');
const response = await fetch('https://gepn-backend.onrender.com/api/denuncia/crear', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({...})
});
```

---

### ✅ Usuarios de Prueba

**Para probar la app durante desarrollo:**

**Ciudadano:**
- Cédula: `V-12345678`
- Contraseña: `test123`

**Policial:**
- Credencial: `POL001`
- PIN: `123456`

**Master/Admin:**
- Usuario: `admin`
- Contraseña: `admin123`

---

### ✅ Formato de Respuestas

**Éxito (200-201):**
```json
{
  "token": "eyJhbGc...",
  "data": {...}
}
```

**Error (400-500):**
```json
{
  "error": "Descripción del error"
}
```

---

### ✅ Configuración Recomendada para la App

**React Native:**
```typescript
// config/api.ts
export const API_BASE_URL = 'https://gepn-backend.onrender.com';
export const API_TIMEOUT = 10000; // 10 segundos
```

**Flutter:**
```dart
// lib/config/api_config.dart
class ApiConfig {
  static const String baseUrl = 'https://gepn-backend.onrender.com';
  static const Duration timeout = Duration(seconds: 10);
}
```

---

### ⚠️ Notas Importantes

1. **HTTPS es obligatorio**: iOS bloquea conexiones HTTP por defecto
2. **Timeout recomendado**: 10-15 segundos (Render free tier puede tardar en responder si está dormido)
3. **Token JWT**: Guardar en AsyncStorage (React Native) o SharedPreferences (Flutter)
4. **Manejo de errores**: Implementar retry logic para errores de red
5. **Render Free Tier**: El servidor se duerme después de 15 minutos de inactividad, la primera petición después puede tardar 30-60 segundos

---

### 📧 Contacto

Si tienen problemas o preguntas sobre el backend:
- Revisar logs en: https://dashboard.render.com
- Verificar health check: https://gepn-backend.onrender.com/health
- Reportar errores con: código de estado HTTP + mensaje de error completo

---

## 🆘 SOLUCIÓN DE PROBLEMAS COMUNES

### ❌ Error: "Failed to connect to MongoDB"

**Síntomas**: El servidor no inicia, logs muestran error de conexión a MongoDB

**Soluciones**:
1. Verifica que la `MONGODB_URI` esté correcta en las variables de entorno
2. Asegúrate de haber reemplazado `<password>` con tu contraseña real
3. Verifica que la contraseña no tenga caracteres especiales sin codificar
   - Si tiene caracteres como `@`, `#`, `%`, debes codificarlos:
   - `@` → `%40`
   - `#` → `%23`
   - `%` → `%25`
4. Verifica en MongoDB Atlas → Network Access que `0.0.0.0/0` está permitido
5. Verifica que el usuario de MongoDB tenga permisos de lectura/escritura

**Ejemplo de URI correcta**:
```
mongodb+srv://gepn_user:MiPass%40word123@cluster.mongodb.net/gepn?retryWrites=true&w=majority
```

---

### ❌ Error: "CORS policy has blocked the request"

**Síntomas**: La app móvil o navegador muestra error de CORS

**Soluciones**:
1. El middleware CORS ya está configurado en `middleware/middleware.go`
2. Verifica que el código desplegado es el más reciente
3. Verifica que el middleware se está aplicando correctamente en `main.go`:
   ```go
   handler := middleware.CORSMiddleware(middleware.LoggingMiddleware(mux))
   ```
4. Para apps móviles, CORS generalmente no es problema (las apps no envían `Origin` header)
5. Si es desde navegador web, verifica que el origen esté permitido

---

### ⏱️ Backend muy lento en Render (Plan Gratuito)

**Síntomas**: La primera petición tarda 30-60 segundos en responder

**Causa**: El plan gratuito de Render duerme el servicio después de 15 minutos de inactividad

**Soluciones**:
1. **Esperar**: La primera petición despertará el servidor, las siguientes serán rápidas
2. **Ping service**: Crear un servicio que haga peticiones cada 10 minutos
3. **Upgrade**: Pagar plan Starter de Render ($7/mes) que nunca se duerme
4. **Usar Railway**: Tiene comportamiento similar pero límites diferentes
5. **Implementar en la app**:
   ```typescript
   // Mostrar loading mientras el servidor despierta
   const [isWaking, setIsWaking] = useState(false);
   
   const makeRequest = async () => {
     setIsWaking(true);
     try {
       const response = await fetch(API_URL, { timeout: 60000 });
       // ...
     } finally {
       setIsWaking(false);
     }
   };
   ```

---

### ❌ Error: 502 Bad Gateway

**Síntomas**: Render muestra "502 Bad Gateway"

**Causas comunes**:
1. El servidor está iniciando (espera 1-2 minutos)
2. El servidor crasheó durante el inicio
3. El puerto no está configurado correctamente

**Soluciones**:
1. Revisa los logs en Render Dashboard
2. Verifica que `PORT=8080` esté en las variables de entorno
3. Verifica que el servidor escucha en el puerto correcto:
   ```go
   port := os.Getenv("PORT")
   if port == "" {
       port = "8080"
   }
   ```
4. Espera 2-3 minutos para que el servidor termine de iniciar

---

### ❌ Error: "Invalid token" o "Unauthorized"

**Síntomas**: Endpoints protegidos devuelven 401 Unauthorized

**Soluciones**:
1. Verifica que el token JWT está en el header:
   ```
   Authorization: Bearer eyJhbGc...
   ```
2. Verifica que el token no ha expirado (validez: 7 días por defecto)
3. Verifica que `JWT_SECRET` es el mismo en todas las variables de entorno
4. No uses tokens generados en desarrollo local con el servidor de producción
5. Obtén un nuevo token haciendo login nuevamente

---

### ❌ Error: "Cannot read environment variables"

**Síntomas**: El servidor inicia pero no puede conectarse a MongoDB

**Soluciones**:
1. Verifica que todas las variables de entorno están configuradas en Render:
   - `MONGODB_URI`
   - `MONGODB_DB_NAME`
   - `JWT_SECRET`
   - `PORT`
2. En Render, ve a Environment → Variables y verifica que existen
3. Después de cambiar variables, redeploy el servicio
4. Verifica los logs para ver qué variable falta

---

### ❌ Error: "Database collection not found"

**Síntomas**: Error al crear/obtener documentos de MongoDB

**Soluciones**:
1. MongoDB Atlas crea colecciones automáticamente al insertar el primer documento
2. Verifica que `MONGODB_DB_NAME=gepn` está configurado
3. Verifica que la inicialización de datos se ejecutó:
   ```
   📦 Inicializando datos por defecto...
   ```
4. Revisa MongoDB Atlas → Browse Collections para ver las colecciones creadas

---

### ❌ Error: Build failed en Render/Railway

**Síntomas**: El despliegue falla durante la fase de build

**Soluciones**:
1. Verifica que `go.mod` y `go.sum` están en el repositorio
2. Verifica que el `Dockerfile` es correcto
3. Revisa los logs de build para ver el error específico
4. Asegúrate de que todas las dependencias están en `go.mod`:
   ```bash
   go mod tidy
   git add go.mod go.sum
   git commit -m "Update dependencies"
   git push
   ```

---

### ❌ La app no puede conectarse al backend

**Síntomas**: La app muestra "Network error" o "Connection failed"

**Soluciones**:
1. Verifica que la URL en la app es correcta y usa `https://`
2. Prueba la URL en el navegador: `https://tu-backend.onrender.com/health`
3. Verifica que el dispositivo/emulador tiene conexión a internet
4. En iOS, verifica que `Info.plist` no bloquea conexiones HTTPS
5. Desactiva temporalmente VPNs o proxies
6. Verifica que no hay firewall bloqueando las peticiones

---

### ⚠️ Render Free Tier: "Service unavailable"

**Síntomas**: Después de tiempo sin uso, el servicio no responde

**Causa**: Render free tier tiene límite de 750 horas/mes

**Soluciones**:
1. Verifica el uso en Render Dashboard
2. El servicio se resetea el 1ro de cada mes
3. Considera upgrade a plan de pago ($7/mes)
4. O usa Railway que tiene modelo de pricing diferente

---

### 📝 Cómo reportar un problema

Si ninguna solución funciona, proporciona:
1. **URL del backend**
2. **Endpoint que falla** (ejemplo: `/api/ciudadano/login`)
3. **Código de error** (ejemplo: 500, 502, 404)
4. **Mensaje de error** completo
5. **Screenshot de los logs** de Render
6. **Request completo** (headers, body)
7. **Respuesta completa** del servidor

---

## 💰 COSTOS ESTIMADOS

### Comparación de Opciones de Hosting:

| Servicio | Plan Gratuito | Limitaciones | Plan Pagado | Recomendado Para |
|----------|---------------|--------------|-------------|------------------|
| **Render** | ✅ Gratis forever | Se duerme después de 15 min sin uso, primera petición ~30-60s | $7/mes - Sin dormirse | Desarrollo y producción pequeña |
| **Railway** | ✅ $5 crédito/mes | 500 horas/mes, después paga por uso | $5/mes + uso (~$10 total) | Producción con tráfico moderado |
| **Google Cloud Run** | ✅ Free tier generoso | 2M requests/mes gratis | Pay-per-use (~$5-15/mes) | Producción profesional |
| **MongoDB Atlas** | ✅ 512MB gratis forever | Suficiente para ~100k documentos | $9/mes (2GB) | Todas las fases |

### Escenario 1: **100% GRATIS** (Desarrollo y MVP)
- **Hosting**: Render Free Tier
- **Base de datos**: MongoDB Atlas M0 (Free)
- **Total**: **$0/mes** 🎉

**Pros**:
- Sin costo
- Fácil de configurar
- HTTPS incluido

**Contras**:
- Servidor se duerme (primera petición lenta)
- Límite de 750 horas/mes en Render

---

### Escenario 2: **RECOMENDADO** (Producción con usuarios reales)
- **Hosting**: Render Starter
- **Base de datos**: MongoDB Atlas M0 (Free)
- **Total**: **$7/mes** 💵

**Pros**:
- Servidor nunca se duerme
- Respuesta rápida siempre
- Base de datos gratis

**Contras**:
- Costo mensual recurrente

---

### Escenario 3: **PROFESIONAL** (Alta demanda)
- **Hosting**: Google Cloud Run
- **Base de datos**: MongoDB Atlas M2 (2GB)
- **Total**: **~$15-20/mes** 💵💵

**Pros**:
- Escalabilidad automática
- Alta disponibilidad
- Mejor rendimiento

**Contras**:
- Más costoso
- Requiere más configuración

---

### Crecimiento Estimado de Costos:

| Usuarios Activos | Requests/día | Render | Railway | Cloud Run | MongoDB |
|------------------|--------------|--------|---------|-----------|---------|
| 0 - 100 | < 1,000 | Gratis | Gratis | Gratis | Gratis |
| 100 - 1,000 | 1,000 - 10,000 | $7/mes | $10/mes | $5-10/mes | Gratis |
| 1,000 - 10,000 | 10,000 - 100,000 | $7/mes | $20/mes | $10-20/mes | $9/mes |
| > 10,000 | > 100,000 | $25/mes+ | $50/mes+ | $30/mes+ | $25/mes+ |

---

### 💡 Recomendación Final:

**Para empezar (Desarrollo/Testing)**:
- ✅ Render Free + MongoDB Atlas Free = **$0/mes**

**Para lanzar la app (Producción)**:
- ✅ Render Starter + MongoDB Atlas Free = **$7/mes**

**Para escalar (Muchos usuarios)**:
- ✅ Cloud Run + MongoDB Atlas M2 = **$15-20/mes**

---

### 🎯 Consejo de Ahorro:

1. **Comienza con el plan gratuito** para validar la app
2. **Upgrade a $7/mes** cuando tengas usuarios reales
3. **Escala según necesidad** basándote en métricas reales
4. **MongoDB gratis es suficiente** hasta ~50,000 usuarios activos

---

## 🎉 ¡LISTO! BACKEND EN PRODUCCIÓN

Una vez completados todos los pasos, tu backend estará funcionando en producción con:

✅ **URL pública con HTTPS**: `https://tu-backend.onrender.com`
✅ **Base de datos MongoDB**: Conectada y funcionando
✅ **CORS configurado**: Apps móviles pueden conectarse
✅ **JWT implementado**: Autenticación segura
✅ **Todos los endpoints**: Listos y probados

---

## 🚀 PRÓXIMOS PASOS

### 1. Para el Equipo del Backend:
- [ ] Completar todos los pasos de este documento
- [ ] Verificar que todos los endpoints funcionen
- [ ] Proporcionar la URL del backend al equipo frontend

### 2. Para el Equipo de la App Móvil:
- [ ] Actualizar la URL en la configuración de la app
- [ ] Probar el registro de usuario desde la app
- [ ] Probar el login desde la app
- [ ] Probar la creación de denuncias
- [ ] Compilar la app (APK para Android, IPA para iOS)
- [ ] Distribuir en Play Store / App Store o mediante link directo

### 3. Para Distribución de la App:

**Android (APK):**
```bash
# React Native
cd android
./gradlew assembleRelease
# APK en: android/app/build/outputs/apk/release/app-release.apk

# Flutter
flutter build apk --release
# APK en: build/app/outputs/flutter-apk/app-release.apk
```

**iOS (TestFlight):**
```bash
# React Native
cd ios
pod install
# Abrir Xcode y hacer Archive → Upload to App Store Connect

# Flutter  
flutter build ios --release
# Abrir Xcode y hacer Archive → Upload to App Store Connect
```

---

## 📋 RESUMEN PARA COMPARTIR

Copia y pega esto al equipo de la app móvil una vez que el backend esté listo:

```
🎉 ¡BACKEND GEPN LISTO PARA PRODUCCIÓN!

✅ URL del Backend: https://gepn-backend.onrender.com
✅ Base de datos: MongoDB Atlas (configurada)
✅ HTTPS: Activo
✅ CORS: Configurado
✅ Estado: Funcionando

📡 Endpoints principales:
- POST /api/ciudadano/registro
- POST /api/ciudadano/login  
- POST /api/policial/login
- POST /api/denuncia/crear
- GET /api/denuncia/mis-denuncias
- GET /health

🔑 Usuarios de prueba:
- Ciudadano: V-12345678 / test123
- Policial: POL001 / 123456
- Master: admin / admin123

📱 Actualizar en la app:
const API_BASE_URL = 'https://gepn-backend.onrender.com';

⚠️ Nota: Render free tier tarda 30-60s en responder la primera petición después de inactividad.

¡Pueden empezar a probar la app con el backend en producción!
```

---

## 📚 DOCUMENTACIÓN ADICIONAL

**Para el equipo técnico**:
- Logs de servidor: https://dashboard.render.com (requiere login)
- MongoDB Dashboard: https://cloud.mongodb.com (requiere login)
- Código fuente del backend: [URL del repositorio Git]

**Archivos importantes en el repositorio**:
- `main.go` - Punto de entrada del servidor
- `routes/routes.go` - Definición de todos los endpoints
- `handlers/` - Lógica de negocio de cada endpoint
- `middleware/middleware.go` - CORS y logging
- `database/database.go` - Conexión a MongoDB
- `Dockerfile` - Configuración de Docker
- `render.yaml` - Configuración de Render

---

## 🆘 SOPORTE Y MANTENIMIENTO

**Si algo falla**:
1. Revisar logs en Render Dashboard
2. Verificar `/health` endpoint
3. Consultar sección "Solución de Problemas" arriba
4. Verificar variables de entorno
5. Revisar MongoDB Atlas → Network Access

**Monitoreo recomendado**:
- Render Dashboard: Revisar logs y métricas
- MongoDB Atlas: Revisar uso de almacenamiento
- Implementar Sentry o similar para tracking de errores
- Configurar alertas en Render para downtime

---

## 📞 CONTACTOS

**Backend Team**: [Email/Slack]
**DevOps**: [Email/Slack]  
**App Mobile Team**: [Email/Slack]

---

## ✅ FIN DE LAS INSTRUCCIONES

Con esto, el backend está **100% listo para producción** y la app móvil puede ser distribuida a los usuarios finales. 🎊
