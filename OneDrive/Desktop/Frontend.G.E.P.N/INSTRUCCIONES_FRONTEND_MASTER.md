# Instrucciones para el Frontend - Sistema Master

Este documento contiene las instrucciones para implementar el login y gestión de usuarios master en el frontend.

## 1. Función de Login Master

En `src/services/apiService.ts`, agregar:

```typescript
export const loginMaster = async (
  usuario: string,
  contraseña: string,
): Promise<{success: boolean; token?: string; master?: any; error?: string}> => {
  try {
    const response = await api.post('/api/master/login', {
      usuario,
      contraseña,
    });

    if (response.data && response.data.token) {
      // Guardar token en AsyncStorage
      await AsyncStorage.setItem('masterToken', response.data.token);
      await AsyncStorage.setItem('masterUser', JSON.stringify(response.data.master));
      
      return {
        success: true,
        token: response.data.token,
        master: response.data.master,
      };
    }
    
    return {
      success: false,
      error: response.data?.error || 'Error desconocido',
    };
  } catch (error: any) {
    console.error('Error en login master:', error);
    
    // Manejar error de respuesta
    if (error.response) {
      return {
        success: false,
        error: error.response.data?.error || 'Error al iniciar sesión',
      };
    }
    
    return {
      success: false,
      error: 'Error de conexión',
    };
  }
};
```

## 2. Respuesta del Backend

El endpoint `POST /api/master/login` retorna:

**Éxito (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "master": {
    "id": "...",
    "usuario": "admin",
    "nombre": "Administrador",
    "email": "admin@gepn.gob.ve",
    "permisos": ["rrhh", "policial", "denuncias", ...],
    "activo": true,
    "creado_por": "sistema",
    "fecha_creacion": "2025-12-22T06:42:06Z"
  },
  "mensaje": "Login exitoso"
}
```

**Error (401):**
```json
{
  "error": "Usuario o contraseña incorrectos"
}
```

## 3. Ejemplo de Uso en Componente

```typescript
import { loginMaster } from '../services/apiService';

const handleLogin = async () => {
  if (!usuario.trim() || !contraseña.trim()) {
    Alert.alert('Error', 'Por favor completa todos los campos');
    return;
  }

  setLoading(true);
  
  const result = await loginMaster(usuario, contraseña);
  
  if (result.success && result.token) {
    // Login exitoso
    Alert.alert('Éxito', 'Login exitoso');
    // Navegar a la pantalla principal del master
    navigation.replace('MasterDashboard');
  } else {
    // Mostrar error
    Alert.alert('Error', result.error || 'Error al iniciar sesión');
  }
  
  setLoading(false);
};
```

## 4. Verificar Token Master

```typescript
export const verificarMaster = async (): Promise<boolean> => {
  try {
    const token = await AsyncStorage.getItem('masterToken');
    if (!token) {
      return false;
    }

    const response = await api.get('/api/master/verificar', {
      headers: {
        Authorization: token,
      },
    });

    if (response.data && response.data.usuario) {
      // Token válido, actualizar usuario en storage
      await AsyncStorage.setItem('masterUser', JSON.stringify(response.data));
      return true;
    }
    
    return false;
  } catch (error) {
    // Token inválido o expirado
    await AsyncStorage.removeItem('masterToken');
    await AsyncStorage.removeItem('masterUser');
    return false;
  }
};
```

## 5. Obtener Token para Requests Protegidos

```typescript
export const getMasterToken = async (): Promise<string | null> => {
  try {
    const token = await AsyncStorage.getItem('masterToken');
    return token;
  } catch (error) {
    return null;
  }
};

// Usar en requests protegidos
const token = await getMasterToken();
const response = await api.get('/api/master/usuarios', {
  headers: {
    Authorization: token,
  },
});
```

## 6. Listar Módulos Disponibles

```typescript
export const obtenerModulos = async (): Promise<string[]> => {
  try {
    const response = await api.get('/api/master/modulos');
    return response.data.modulos || [];
  } catch (error) {
    console.error('Error al obtener módulos:', error);
    return [];
  }
};
```

## 7. Crear Usuario Master

```typescript
export const crearUsuarioMaster = async (
  usuario: string,
  nombre: string,
  email: string,
  contraseña: string,
  permisos: string[],
): Promise<{success: boolean; error?: string}> => {
  try {
    const token = await getMasterToken();
    if (!token) {
      return {success: false, error: 'No autenticado'};
    }

    const response = await api.post(
      '/api/master/crear-usuario',
      {
        usuario,
        nombre,
        email,
        contraseña,
        permisos,
      },
      {
        headers: {
          Authorization: token,
        },
      },
    );

    if (response.status === 201) {
      return {success: true};
    }
    
    return {success: false, error: 'Error al crear usuario'};
  } catch (error: any) {
    return {
      success: false,
      error: error.response?.data?.error || 'Error al crear usuario',
    };
  }
};
```

## 8. Listar Usuarios Master

```typescript
export const listarUsuariosMaster = async (): Promise<any[]> => {
  try {
    const token = await getMasterToken();
    if (!token) {
      return [];
    }

    const response = await api.get('/api/master/usuarios', {
      headers: {
        Authorization: token,
      },
    });

    return response.data || [];
  } catch (error) {
    console.error('Error al listar usuarios master:', error);
    return [];
  }
};
```

## 9. Actualizar Permisos

```typescript
export const actualizarPermisos = async (
  usuarioId: string,
  permisos: string[],
): Promise<{success: boolean; error?: string}> => {
  try {
    const token = await getMasterToken();
    if (!token) {
      return {success: false, error: 'No autenticado'};
    }

    const response = await api.put(
      `/api/master/usuarios/permisos/${usuarioId}`,
      {permisos},
      {
        headers: {
          Authorization: token,
        },
      },
    );

    if (response.status === 200) {
      return {success: true};
    }
    
    return {success: false, error: 'Error al actualizar permisos'};
  } catch (error: any) {
    return {
      success: false,
      error: error.response?.data?.error || 'Error al actualizar permisos',
    };
  }
};
```

## 10. Activar/Desactivar Usuario

```typescript
export const activarUsuarioMaster = async (
  usuarioId: string,
  activo: boolean,
): Promise<{success: boolean; error?: string}> => {
  try {
    const token = await getMasterToken();
    if (!token) {
      return {success: false, error: 'No autenticado'};
    }

    const response = await api.put(
      `/api/master/usuarios/activar/${usuarioId}`,
      {activo},
      {
        headers: {
          Authorization: token,
        },
      },
    );

    if (response.status === 200) {
      return {success: true};
    }
    
    return {success: false, error: 'Error al actualizar estado'};
  } catch (error: any) {
    return {
      success: false,
      error: error.response?.data?.error || 'Error al actualizar estado',
    };
  }
};
```

## 11. Verificar Permisos del Usuario Actual

```typescript
export const tienePermiso = async (modulo: string): Promise<boolean> => {
  try {
    const masterUser = await AsyncStorage.getItem('masterUser');
    if (!masterUser) {
      return false;
    }

    const master = JSON.parse(masterUser);
    return master.permisos && master.permisos.includes(modulo);
  } catch (error) {
    return false;
  }
};
```

## 12. Cerrar Sesión Master

```typescript
export const logoutMaster = async (): Promise<void> => {
  await AsyncStorage.removeItem('masterToken');
  await AsyncStorage.removeItem('masterUser');
};
```

## 13. Interceptor de Axios para Agregar Token Automáticamente

En `src/services/apiService.ts`:

```typescript
// Interceptor para agregar token automáticamente a requests protegidos
api.interceptors.request.use(
  async (config) => {
    // Solo agregar token a rutas de master
    if (config.url?.includes('/api/master/') && !config.url?.includes('/login') && !config.url?.includes('/modulos')) {
      const token = await AsyncStorage.getItem('masterToken');
      if (token) {
        config.headers.Authorization = token;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Interceptor para manejar errores 401
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // Token expirado o inválido
      await AsyncStorage.removeItem('masterToken');
      await AsyncStorage.removeItem('masterUser');
      // Redirigir a login
      // navigation.navigate('MasterLogin');
    }
    return Promise.reject(error);
  },
);
```

## 14. Estructura de Respuesta Esperada

### Login Exitoso
El frontend debe verificar:
- `response.status === 200`
- `response.data.token` existe
- `response.data.master` existe

### Login Fallido
El frontend debe verificar:
- `response.status === 401` o `response.status === 403`
- `response.data.error` contiene el mensaje de error

## 15. Ejemplo Completo de Pantalla de Login

```typescript
import React, {useState} from 'react';
import {View, TextInput, TouchableOpacity, Text, Alert} from 'react-native';
import {loginMaster} from '../services/apiService';

const MasterLoginScreen = ({navigation}) => {
  const [usuario, setUsuario] = useState('');
  const [contraseña, setContraseña] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    if (!usuario.trim() || !contraseña.trim()) {
      Alert.alert('Error', 'Por favor completa todos los campos');
      return;
    }

    setLoading(true);

    try {
      const result = await loginMaster(usuario, contraseña);

      if (result.success && result.token) {
        Alert.alert('Éxito', 'Login exitoso');
        navigation.replace('MasterDashboard');
      } else {
        Alert.alert('Error', result.error || 'Error al iniciar sesión');
      }
    } catch (error) {
      Alert.alert('Error', 'Error de conexión');
    } finally {
      setLoading(false);
    }
  };

  return (
    <View>
      <TextInput
        placeholder="Usuario"
        value={usuario}
        onChangeText={setUsuario}
        autoCapitalize="none"
      />
      <TextInput
        placeholder="Contraseña"
        value={contraseña}
        onChangeText={setContraseña}
        secureTextEntry
      />
      <TouchableOpacity onPress={handleLogin} disabled={loading}>
        <Text>{loading ? 'Cargando...' : 'Iniciar Sesión'}</Text>
      </TouchableOpacity>
    </View>
  );
};
```

## 16. Credenciales por Defecto

**Usuario:** `admin`  
**Contraseña:** `Admin123!` (A mayúscula, resto minúsculas, números y `!`)

## 17. Notas Importantes

1. El token JWT expira después de 24 horas
2. El token debe incluirse en el header `Authorization` para todas las rutas protegidas
3. Las rutas públicas son:
   - `POST /api/master/login`
   - `GET /api/master/modulos`
4. Todas las demás rutas requieren autenticación
5. El campo `permisos` es un array de strings con los módulos a los que tiene acceso

---

**Estado**: ✅ Listo para implementar
**Fecha**: 2025-12-22

