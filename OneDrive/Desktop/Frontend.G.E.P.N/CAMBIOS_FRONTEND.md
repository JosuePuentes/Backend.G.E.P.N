# Cambios Necesarios en el Frontend

## Resumen
Este documento contiene todos los cambios que necesitas implementar en el frontend React Native para que funcione con el nuevo sistema de guardias policiales y notificaciones de pánico.

---

## 1. Actualizar `src/services/apiService.ts`

### Cambio en `loginPolicial`:
```typescript
// ANTES:
export const loginPolicial = async (
  credencial: string,
  pin: string,
): Promise<boolean> => {
  // ...
}

// DESPUÉS:
export const loginPolicial = async (
  credencial: string,
  pin: string,
  latitud?: number,
  longitud?: number,
): Promise<{success: boolean; token?: string; usuario?: any}> => {
  try {
    const body: any = {credencial, pin};
    if (latitud !== undefined && longitud !== undefined) {
      body.latitud = latitud;
      body.longitud = longitud;
    }

    const response = await api.post('/api/policial/login', body);

    if (response.data && response.data.token) {
      await AsyncStorage.setItem('authToken', response.data.token);
      await AsyncStorage.setItem('policial_user', JSON.stringify(response.data.usuario));
      await AsyncStorage.setItem('guardia_activa', 'true');
      return {success: true, token: response.data.token, usuario: response.data.usuario};
    }
    return {success: false};
  } catch (error) {
    console.error('Error en login:', error);
    return {success: false};
  }
};
```

### Agregar nueva función `finalizarGuardia`:
```typescript
export const finalizarGuardia = async (): Promise<boolean> => {
  try {
    const response = await api.post('/api/policial/finalizar-guardia');
    if (response.status === 200) {
      await AsyncStorage.removeItem('guardia_activa');
      return true;
    }
    return false;
  } catch (error) {
    console.error('Error al finalizar guardia:', error);
    return false;
  }
};
```

---

## 2. Actualizar `src/screens/HomeScreen.tsx`

### Agregar botón de "Acceso Policial" después del botón de "Realizar Denuncia":

```typescript
{/* Botón de Acceso Policial */}
<TouchableOpacity
  style={styles.policialButton}
  onPress={() => {
    navigation.navigate('LoginPolicial');
  }}
  activeOpacity={0.7}>
  <View style={styles.policialButtonContent}>
    <Text style={styles.policialIcon}>👮</Text>
    <Text style={styles.policialButtonText}>Acceso Policial</Text>
  </View>
</TouchableOpacity>
```

### Agregar estilos al final de `StyleSheet.create`:

```typescript
policialButton: {
  width: '100%',
  maxWidth: 400,
  backgroundColor: '#00247D',
  borderRadius: 16,
  paddingVertical: 25,
  paddingHorizontal: 30,
  marginBottom: 20,
  shadowColor: '#00247D',
  shadowOffset: {
    width: 0,
    height: 8,
  },
  shadowOpacity: 0.5,
  shadowRadius: 12,
  elevation: 15,
  borderWidth: 2,
  borderColor: '#0033A0',
},
policialButtonContent: {
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
},
policialIcon: {
  fontSize: 32,
  marginRight: 15,
},
policialButtonText: {
  color: '#FFFFFF',
  fontSize: 22,
  fontWeight: 'bold',
  letterSpacing: 1.5,
},
```

**Nota:** Este botón debe agregarse en AMBAS secciones del código (con ImageBackground y sin ImageBackground).

---

## 3. Actualizar `src/screens/LoginPolicialScreen.tsx`

### Cambiar la función `handleLogin`:

```typescript
const handleLogin = async () => {
  if (!credencial.trim() || !pin.trim()) {
    Alert.alert('Error', 'Por favor completa todos los campos');
    return;
  }

  if (pin.length !== 6 || !/^\d+$/.test(pin)) {
    Alert.alert('Error', 'El PIN debe tener 6 dígitos numéricos');
    return;
  }

  setLoading(true);

  try {
    // Solicitar permisos GPS
    const hasPermission = await requestLocationPermission();
    if (!hasPermission) {
      Alert.alert(
        'Permisos requeridos',
        'Se necesitan permisos de ubicación para iniciar guardia',
      );
      setLoading(false);
      return;
    }

    // Obtener ubicación GPS
    Geolocation.getCurrentPosition(
      async position => {
        const {latitude, longitude} = position.coords;
        
        // Realizar login con GPS
        const result = await loginPolicial(credencial, pin, latitude, longitude);
        if (result.success) {
          Alert.alert('Éxito', 'Guardia iniciada correctamente');
          navigation.replace('Dashboard');
        } else {
          Alert.alert('Error', 'Credenciales incorrectas');
        }
        setLoading(false);
      },
      error => {
        Alert.alert('Error', 'No se pudo obtener la ubicación. Intenta nuevamente.');
        setLoading(false);
      },
      {enableHighAccuracy: true, timeout: 15000, maximumAge: 10000},
    );
  } catch (error) {
    Alert.alert('Error', 'Error al iniciar sesión. Intenta nuevamente.');
    setLoading(false);
  }
};
```

### Agregar import de `Geolocation` al inicio del archivo (si no existe):

```typescript
// Importar Geolocation según la plataforma
let Geolocation: any;
if (Platform.OS === 'web') {
  Geolocation = {
    getCurrentPosition: (
      success: (position: any) => void,
      error: (error: any) => void,
      options: any,
    ) => {
      if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success, error, options);
      } else {
        error({message: 'Geolocation no está soportado'});
      }
    },
  };
} else {
  Geolocation = require('@react-native-community/geolocation').default;
}
```

---

## 4. Actualizar `src/screens/DashboardScreen.tsx`

### Cambiar los 4 botones del grid:

```typescript
<View style={styles.grid}>
  <TouchableOpacity
    style={styles.menuButton}
    onPress={() => handleMenuButton('Minutas Digitales')}>
    <Text style={styles.menuButtonIcon}>📝</Text>
    <Text style={styles.menuButtonText}>Minutas Digitales</Text>
  </TouchableOpacity>

  <TouchableOpacity
    style={styles.menuButton}
    onPress={() => handleMenuButton('Los Más Buscados')}>
    <Text style={styles.menuButtonIcon}>🔍</Text>
    <Text style={styles.menuButtonText}>Los Más Buscados</Text>
  </TouchableOpacity>

  <TouchableOpacity
    style={styles.menuButton}
    onPress={() => handleMenuButton('Verificación de Cédulas')}>
    <Text style={styles.menuButtonIcon}>🆔</Text>
    <Text style={styles.menuButtonText}>Verificación de Cédulas</Text>
  </TouchableOpacity>

  <TouchableOpacity
    style={styles.menuButton}
    onPress={() => handleMenuButton('Registro de Detenidos')}>
    <Text style={styles.menuButtonIcon}>👤</Text>
    <Text style={styles.menuButtonText}>Registro de Detenidos</Text>
  </TouchableOpacity>
</View>
```

### Actualizar el botón de pánico:

```typescript
<View style={styles.panicContainer}>
  <Animated.View style={{transform: [{scale: scaleAnim}]}}>
    <TouchableOpacity
      style={[
        styles.panicButton,
        panicPressed && styles.panicButtonPressed,
      ]}
      onPressIn={handlePanicPressIn}
      onPressOut={handlePanicPressOut}
      activeOpacity={0.8}>
      <Text style={styles.panicButtonIcon}>🚨</Text>
      <Text style={styles.panicButtonText}>
        {panicPressed ? 'Mantén presionado 5 segundos...' : 'Botón de Apoyo'}
      </Text>
    </TouchableOpacity>
  </Animated.View>
  
  <TouchableOpacity
    style={styles.finalizarButton}
    onPress={handleFinalizarGuardia}
    activeOpacity={0.8}>
    <Text style={styles.finalizarButtonText}>Finalizar Guardia</Text>
  </TouchableOpacity>
</View>
```

### Agregar función `handleFinalizarGuardia`:

```typescript
import {finalizarGuardia} from '../services/apiService';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Dentro del componente:
const handleFinalizarGuardia = async () => {
  Alert.alert(
    'Finalizar Guardia',
    '¿Estás seguro de que deseas finalizar tu guardia?',
    [
      {text: 'Cancelar', style: 'cancel'},
      {
        text: 'Finalizar',
        style: 'destructive',
        onPress: async () => {
          const success = await finalizarGuardia();
          if (success) {
            await AsyncStorage.removeItem('authToken');
            await AsyncStorage.removeItem('policial_user');
            Alert.alert('Éxito', 'Guardia finalizada correctamente');
            navigation.replace('Home');
          } else {
            Alert.alert('Error', 'No se pudo finalizar la guardia');
          }
        },
      },
    ],
  );
};
```

### Actualizar estilos:

```typescript
menuButtonIcon: {
  fontSize: 40,
  marginBottom: 10,
},
menuButtonText: {
  fontSize: 16,
  fontWeight: '600',
  color: '#D4AF37',
  textAlign: 'center',
},
panicButton: {
  backgroundColor: '#FF3B30',
  paddingVertical: 25,
  paddingHorizontal: 50,
  borderRadius: 16,
  minWidth: 250,
  alignItems: 'center',
  shadowColor: '#FF3B30',
  shadowOffset: {
    width: 0,
    height: 8,
  },
  shadowOpacity: 0.6,
  shadowRadius: 12,
  elevation: 15,
  borderWidth: 3,
  borderColor: '#FF6B60',
},
panicButtonPressed: {
  backgroundColor: '#CC2E24',
  borderColor: '#FF3B30',
},
panicButtonIcon: {
  fontSize: 48,
  marginBottom: 10,
},
panicButtonText: {
  color: '#fff',
  fontSize: 20,
  fontWeight: 'bold',
  textAlign: 'center',
},
finalizarButton: {
  marginTop: 20,
  backgroundColor: '#2a2a2a',
  paddingVertical: 15,
  paddingHorizontal: 30,
  borderRadius: 10,
  borderWidth: 1,
  borderColor: '#3a3a3a',
},
finalizarButtonText: {
  color: '#CCCCCC',
  fontSize: 16,
  fontWeight: '600',
  textAlign: 'center',
},
```

### Actualizar la firma del componente para recibir `navigation`:

```typescript
// ANTES:
const DashboardScreen: React.FC<Props> = () => {

// DESPUÉS:
const DashboardScreen: React.FC<Props> = ({navigation}) => {
```

---

## 5. Verificar que `App.tsx` tenga la ruta `LoginPolicial`

Asegúrate de que en `App.tsx` exista:

```typescript
<Stack.Screen name="LoginPolicial" component={LoginPolicialScreen} />
```

---

## Resumen de Cambios

1. ✅ **apiService.ts**: Actualizar `loginPolicial` y agregar `finalizarGuardia`
2. ✅ **HomeScreen.tsx**: Agregar botón "Acceso Policial" y estilos
3. ✅ **LoginPolicialScreen.tsx**: Obtener GPS y enviarlo al login
4. ✅ **DashboardScreen.tsx**: Actualizar 4 iconos, botón de pánico mejorado, y botón de finalizar guardia

---

## Notas Importantes

- El botón de pánico requiere mantener presionado 5 segundos para activarse
- Al iniciar sesión policial, se solicita GPS automáticamente
- Al finalizar guardia, se limpia el token y se redirige al Home
- Los oficiales cercanos (dentro de 5 km) recibirán notificación cuando se active el pánico

