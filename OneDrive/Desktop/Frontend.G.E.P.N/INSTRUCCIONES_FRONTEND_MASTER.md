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

## 18. Pantalla de Dashboard Master - Mostrar Módulos

Ejemplo completo de cómo mostrar todos los módulos en la pantalla del master:

```typescript
import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  ScrollView,
  TouchableOpacity,
  StyleSheet,
  Alert,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import {obtenerModulos, tienePermiso, logoutMaster} from '../services/apiService';

// Mapeo de módulos a nombres y rutas
const modulosInfo = {
  rrhh: {
    nombre: 'Recursos Humanos',
    icono: '👥',
    ruta: 'RRHHDashboard',
    descripcion: 'Gestionar oficiales y personal',
  },
  policial: {
    nombre: 'Módulo Policial',
    icono: '👮',
    ruta: 'PolicialDashboard',
    descripcion: 'Gestión de guardias y operaciones',
  },
  denuncias: {
    nombre: 'Denuncias',
    icono: '📋',
    ruta: 'DenunciasDashboard',
    descripcion: 'Gestionar denuncias ciudadanas',
  },
  detenidos: {
    nombre: 'Detenidos',
    icono: '🔒',
    ruta: 'DetenidosDashboard',
    descripcion: 'Registro de detenidos',
  },
  minutas: {
    nombre: 'Minutas Digitales',
    icono: '📝',
    ruta: 'MinutasDashboard',
    descripcion: 'Crear y gestionar minutas',
  },
  buscados: {
    nombre: 'Más Buscados',
    icono: '🔍',
    ruta: 'BuscadosDashboard',
    descripcion: 'Lista de personas buscadas',
  },
  verificacion: {
    nombre: 'Verificación de Cédulas',
    icono: '🆔',
    ruta: 'VerificacionDashboard',
    descripcion: 'Verificar cédulas de identidad',
  },
  panico: {
    nombre: 'Botón de Pánico',
    icono: '🚨',
    ruta: 'PanicoDashboard',
    descripcion: 'Alertas y emergencias',
  },
};

const MasterDashboardScreen = ({navigation}) => {
  const [modulos, setModulos] = useState<string[]>([]);
  const [permisosUsuario, setPermisosUsuario] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    cargarDatos();
  }, []);

  const cargarDatos = async () => {
    try {
      // Obtener módulos disponibles del servidor
      const modulosDisponibles = await obtenerModulos();
      setModulos(modulosDisponibles);

      // Obtener permisos del usuario actual
      const masterUser = await AsyncStorage.getItem('masterUser');
      if (masterUser) {
        const master = JSON.parse(masterUser);
        setPermisosUsuario(master.permisos || []);
      }
    } catch (error) {
      console.error('Error al cargar datos:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleModuloPress = async (modulo: string) => {
    // Verificar si el usuario tiene permiso para este módulo
    const tieneAcceso = permisosUsuario.includes(modulo);
    
    if (!tieneAcceso) {
      Alert.alert(
        'Acceso Denegado',
        'No tienes permisos para acceder a este módulo',
      );
      return;
    }

    const info = modulosInfo[modulo];
    if (info) {
      navigation.navigate(info.ruta);
    }
  };

  const handleLogout = async () => {
    Alert.alert(
      'Cerrar Sesión',
      '¿Estás seguro de que deseas cerrar sesión?',
      [
        {text: 'Cancelar', style: 'cancel'},
        {
          text: 'Cerrar Sesión',
          style: 'destructive',
          onPress: async () => {
            await logoutMaster();
            navigation.replace('MasterLogin');
          },
        },
      ],
    );
  };

  if (loading) {
    return (
      <View style={styles.container}>
        <Text>Cargando...</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.title}>Panel de Control Master</Text>
        <TouchableOpacity onPress={handleLogout} style={styles.logoutButton}>
          <Text style={styles.logoutText}>Cerrar Sesión</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.modulosContainer}>
        <Text style={styles.sectionTitle}>Módulos Disponibles</Text>
        
        {modulos.map((modulo) => {
          const info = modulosInfo[modulo];
          const tieneAcceso = permisosUsuario.includes(modulo);
          
          if (!info) return null;

          return (
            <TouchableOpacity
              key={modulo}
              style={[
                styles.moduloCard,
                !tieneAcceso && styles.moduloCardDisabled,
              ]}
              onPress={() => handleModuloPress(modulo)}
              disabled={!tieneAcceso}>
              <View style={styles.moduloContent}>
                <Text style={styles.moduloIcon}>{info.icono}</Text>
                <View style={styles.moduloInfo}>
                  <Text style={styles.moduloNombre}>{info.nombre}</Text>
                  <Text style={styles.moduloDescripcion}>
                    {info.descripcion}
                  </Text>
                  {!tieneAcceso && (
                    <Text style={styles.sinPermiso}>
                      Sin acceso a este módulo
                    </Text>
                  )}
                </View>
                {tieneAcceso && (
                  <Text style={styles.arrow}>→</Text>
                )}
              </View>
            </TouchableOpacity>
          );
        })}
      </View>

      <View style={styles.permisosContainer}>
        <Text style={styles.sectionTitle}>Tus Permisos</Text>
        <View style={styles.permisosList}>
          {permisosUsuario.length > 0 ? (
            permisosUsuario.map((permiso) => {
              const info = modulosInfo[permiso];
              return (
                <View key={permiso} style={styles.permisoTag}>
                  <Text style={styles.permisoText}>
                    {info ? info.icono + ' ' + info.nombre : permiso}
                  </Text>
                </View>
              );
            })
          ) : (
            <Text style={styles.sinPermisos}>
              No tienes permisos asignados
            </Text>
          )}
        </View>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#00247D',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#FFFFFF',
  },
  logoutButton: {
    padding: 10,
    backgroundColor: '#FF3B30',
    borderRadius: 8,
  },
  logoutText: {
    color: '#FFFFFF',
    fontWeight: '600',
  },
  modulosContainer: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 15,
    color: '#333',
  },
  moduloCard: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 20,
    marginBottom: 15,
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
    borderLeftWidth: 4,
    borderLeftColor: '#00247D',
  },
  moduloCardDisabled: {
    opacity: 0.5,
    borderLeftColor: '#CCCCCC',
  },
  moduloContent: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  moduloIcon: {
    fontSize: 40,
    marginRight: 15,
  },
  moduloInfo: {
    flex: 1,
  },
  moduloNombre: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 5,
  },
  moduloDescripcion: {
    fontSize: 14,
    color: '#666',
  },
  sinPermiso: {
    fontSize: 12,
    color: '#FF3B30',
    marginTop: 5,
    fontStyle: 'italic',
  },
  arrow: {
    fontSize: 24,
    color: '#00247D',
  },
  permisosContainer: {
    padding: 20,
    backgroundColor: '#FFFFFF',
    marginTop: 20,
  },
  permisosList: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginTop: 10,
  },
  permisoTag: {
    backgroundColor: '#E3F2FD',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 20,
    marginRight: 10,
    marginBottom: 10,
  },
  permisoText: {
    color: '#00247D',
    fontSize: 12,
    fontWeight: '600',
  },
  sinPermisos: {
    color: '#999',
    fontStyle: 'italic',
  },
});

export default MasterDashboardScreen;
```

## 19. Mapeo de Módulos a Pantallas

Cada módulo debe tener su propia pantalla de dashboard:

```typescript
// Ejemplo de estructura de rutas
const MasterStack = () => {
  return (
    <Stack.Navigator>
      <Stack.Screen name="MasterLogin" component={MasterLoginScreen} />
      <Stack.Screen name="MasterDashboard" component={MasterDashboardScreen} />
      
      {/* Pantallas de cada módulo */}
      <Stack.Screen name="RRHHDashboard" component={RRHHDashboardScreen} />
      <Stack.Screen name="PolicialDashboard" component={PolicialDashboardScreen} />
      <Stack.Screen name="DenunciasDashboard" component={DenunciasDashboardScreen} />
      <Stack.Screen name="DetenidosDashboard" component={DetenidosDashboardScreen} />
      <Stack.Screen name="MinutasDashboard" component={MinutasDashboardScreen} />
      <Stack.Screen name="BuscadosDashboard" component={BuscadosDashboardScreen} />
      <Stack.Screen name="VerificacionDashboard" component={VerificacionDashboardScreen} />
      <Stack.Screen name="PanicoDashboard" component={PanicoDashboardScreen} />
    </Stack.Navigator>
  );
};
```

## 20. Información de Cada Módulo

### RRHH (Recursos Humanos)
- **Ruta:** `/api/rrhh/*`
- **Funcionalidades:**
  - Registrar oficiales
  - Listar oficiales
  - Generar QR codes
  - Gestionar ascensos
  - Verificar QR

### Policial
- **Ruta:** `/api/policial/*`
- **Funcionalidades:**
  - Login de oficiales
  - Finalizar guardias
  - Ver guardias activas

### Denuncias
- **Ruta:** `/api/denuncia/*`
- **Funcionalidades:**
  - Ver denuncias
  - Gestionar estado de denuncias
  - Estadísticas

### Detenidos
- **Ruta:** `/api/detenidos/*`
- **Funcionalidades:**
  - Registrar detenidos
  - Listar detenidos
  - Actualizar estado

### Minutas
- **Ruta:** `/api/minutas/*`
- **Funcionalidades:**
  - Crear minutas
  - Listar minutas
  - Ver detalles

### Buscados
- **Ruta:** `/api/mas-buscados`
- **Funcionalidades:**
  - Ver lista de más buscados
  - Agregar/eliminar buscados

### Verificación
- **Ruta:** `/api/buscar/cedula`
- **Funcionalidades:**
  - Verificar cédulas
  - Historial de búsquedas

### Pánico
- **Ruta:** `/api/panico/*`
- **Funcionalidades:**
  - Ver alertas de pánico
  - Gestionar alertas activas

---

**Estado**: ✅ Listo para implementar
**Fecha**: 2025-12-22

