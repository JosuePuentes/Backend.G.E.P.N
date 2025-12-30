# Instrucciones Frontend - Sistema de Denuncias

Este documento contiene todas las instrucciones para implementar la funcionalidad de visualización de denuncias en el frontend.

## 📋 Resumen

Cuando un ciudadano envía una denuncia, esa denuncia queda guardada en la base de datos y puede ser vista por:
1. **El ciudadano que la creó** - A través de "Mis Denuncias"
2. **Usuarios del sistema con permiso "denuncias"** - Pueden ver todas las denuncias y gestionar su estado

## 🔌 Endpoints Disponibles

### Para Ciudadanos

#### 1. Crear Denuncia
- **Endpoint:** `POST /api/denuncia/crear`
- **Autenticación:** Requiere token de ciudadano (Bearer token)
- **Request Body:**
```json
{
  "denunciante": {
    "nombre": "Juan Pérez",
    "cedula": "12345678",
    "telefono": "04121234567",
    "fechaNacimiento": "15/05/1990",
    "parroquia": "Parroquia X"
  },
  "denuncia": {
    "motivo": "Robo",
    "hechos": "Descripción detallada de los hechos"
  },
  "denunciado": {
    "nombre": "Persona Denunciada",
    "direccion": "Dirección del denunciado",
    "estado": "Estado",
    "municipio": "Municipio",
    "parroquia": "Parroquia"
  }
}
```

- **Response:**
```json
{
  "success": true,
  "message": "Denuncia registrada correctamente",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "numero_denuncia": "DEN-2025-0001",
    "fecha": "2025-01-15T10:30:00Z"
  }
}
```

#### 2. Mis Denuncias (Ciudadano)
- **Endpoint:** `GET /api/denuncia/mis-denuncias`
- **Autenticación:** Requiere token de ciudadano (Bearer token)
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "numero_denuncia": "DEN-2025-0001",
      "motivo": "Robo",
      "fecha_denuncia": "2025-01-15T10:30:00Z",
      "estado": "Pendiente"
    }
  ]
}
```

### Para Usuarios del Sistema (Master con permiso "denuncias")

#### 3. Listar Todas las Denuncias
- **Endpoint:** `GET /api/denuncia/listar`
- **Autenticación:** Requiere token de master con permiso "denuncias" (Bearer token)
- **Query Parameters:**
  - `page` (opcional): Número de página (default: 1)
  - `limit` (opcional): Cantidad de resultados por página (default: 20)
  - `estado` (opcional): Filtrar por estado ("Pendiente", "En Proceso", "Resuelta", "Archivada")

- **Ejemplo de uso:**
```
GET /api/denuncia/listar?page=1&limit=20&estado=Pendiente
```

- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "ciudadano_id": "507f1f77bcf86cd799439012",
      "numero_denuncia": "DEN-2025-0001",
      "nombre_denunciante": "Juan Pérez",
      "cedula_denunciante": "12345678",
      "telefono_denunciante": "04121234567",
      "fecha_nacimiento_denunciante": "1990-05-15",
      "parroquia_denunciante": "Parroquia X",
      "motivo": "Robo",
      "hechos": "Descripción detallada de los hechos",
      "nombre_denunciado": "Persona Denunciada",
      "direccion_denunciado": "Dirección del denunciado",
      "estado_denunciado": "Estado",
      "municipio_denunciado": "Municipio",
      "parroquia_denunciado": "Parroquia",
      "fecha_denuncia": "2025-01-15T10:30:00Z",
      "estado": "Pendiente"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50
  }
}
```

#### 4. Obtener Denuncia Específica
- **Endpoint:** `GET /api/denuncia/obtener`
- **Autenticación:** Requiere token de master con permiso "denuncias" (Bearer token)
- **Query Parameters:**
  - `id` (requerido): ID de la denuncia

- **Ejemplo de uso:**
```
GET /api/denuncia/obtener?id=507f1f77bcf86cd799439011
```

- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "ciudadano_id": "507f1f77bcf86cd799439012",
    "numero_denuncia": "DEN-2025-0001",
    "nombre_denunciante": "Juan Pérez",
    "cedula_denunciante": "12345678",
    "telefono_denunciante": "04121234567",
    "fecha_nacimiento_denunciante": "1990-05-15",
    "parroquia_denunciante": "Parroquia X",
    "motivo": "Robo",
    "hechos": "Descripción detallada de los hechos",
    "nombre_denunciado": "Persona Denunciada",
    "direccion_denunciado": "Dirección del denunciado",
    "estado_denunciado": "Estado",
    "municipio_denunciado": "Municipio",
    "parroquia_denunciado": "Parroquia",
    "fecha_denuncia": "2025-01-15T10:30:00Z",
    "estado": "Pendiente"
  }
}
```

#### 5. Actualizar Estado de Denuncia
- **Endpoint:** `PUT /api/denuncia/actualizar-estado`
- **Autenticación:** Requiere token de master con permiso "denuncias" (Bearer token)
- **Request Body:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "estado": "En Proceso"
}
```

- **Estados válidos:**
  - `"Pendiente"`
  - `"En Proceso"`
  - `"Resuelta"`
  - `"Archivada"`

- **Response:**
```json
{
  "success": true,
  "message": "Estado de denuncia actualizado correctamente"
}
```

## 🎨 Implementación en el Frontend

### 1. Pantalla de Listado de Denuncias (Para Usuarios del Sistema)

```javascript
import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';

const ListadoDenunciasScreen = ({ navigation }) => {
  const [denuncias, setDenuncias] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [estadoFiltro, setEstadoFiltro] = useState('');
  const token = 'TU_TOKEN_MASTER_AQUI'; // Obtener del contexto/AsyncStorage

  useEffect(() => {
    cargarDenuncias();
  }, [page, estadoFiltro]);

  const cargarDenuncias = async () => {
    try {
      setLoading(true);
      let url = `https://tu-api.com/api/denuncia/listar?page=${page}&limit=20`;
      if (estadoFiltro) {
        url += `&estado=${estadoFiltro}`;
      }

      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      const data = await response.json();
      
      if (data.success) {
        setDenuncias(data.data);
        setTotal(data.pagination.total);
      }
    } catch (error) {
      console.error('Error al cargar denuncias:', error);
    } finally {
      setLoading(false);
    }
  };

  const renderDenuncia = ({ item }) => (
    <TouchableOpacity
      style={styles.denunciaCard}
      onPress={() => navigation.navigate('DetalleDenuncia', { id: item.id })}
    >
      <View style={styles.denunciaHeader}>
        <Text style={styles.numeroDenuncia}>{item.numero_denuncia}</Text>
        <View style={[styles.badge, getBadgeColor(item.estado)]}>
          <Text style={styles.badgeText}>{item.estado}</Text>
        </View>
      </View>
      <Text style={styles.motivo}>{item.motivo}</Text>
      <Text style={styles.denunciante}>Denunciante: {item.nombre_denunciante}</Text>
      <Text style={styles.fecha}>
        {new Date(item.fecha_denuncia).toLocaleDateString('es-VE')}
      </Text>
    </TouchableOpacity>
  );

  const getBadgeColor = (estado) => {
    const colors = {
      'Pendiente': { backgroundColor: '#FFA500' },
      'En Proceso': { backgroundColor: '#2196F3' },
      'Resuelta': { backgroundColor: '#4CAF50' },
      'Archivada': { backgroundColor: '#9E9E9E' },
    };
    return colors[estado] || { backgroundColor: '#9E9E9E' };
  };

  if (loading && denuncias.length === 0) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#2196F3" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      {/* Filtros */}
      <View style={styles.filtros}>
        <TouchableOpacity
          style={[styles.filtroBtn, estadoFiltro === '' && styles.filtroBtnActive]}
          onPress={() => setEstadoFiltro('')}
        >
          <Text>Todas</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.filtroBtn, estadoFiltro === 'Pendiente' && styles.filtroBtnActive]}
          onPress={() => setEstadoFiltro('Pendiente')}
        >
          <Text>Pendientes</Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[styles.filtroBtn, estadoFiltro === 'En Proceso' && styles.filtroBtnActive]}
          onPress={() => setEstadoFiltro('En Proceso')}
        >
          <Text>En Proceso</Text>
        </TouchableOpacity>
      </View>

      {/* Lista */}
      <FlatList
        data={denuncias}
        renderItem={renderDenuncia}
        keyExtractor={(item) => item.id}
        onEndReached={() => {
          if (denuncias.length < total) {
            setPage(page + 1);
          }
        }}
        onEndReachedThreshold={0.5}
        ListEmptyComponent={
          <View style={styles.center}>
            <Text>No hay denuncias disponibles</Text>
          </View>
        }
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  filtros: {
    flexDirection: 'row',
    padding: 10,
    backgroundColor: '#fff',
    borderBottomWidth: 1,
    borderBottomColor: '#e0e0e0',
  },
  filtroBtn: {
    paddingHorizontal: 15,
    paddingVertical: 8,
    marginRight: 10,
    borderRadius: 20,
    backgroundColor: '#e0e0e0',
  },
  filtroBtnActive: {
    backgroundColor: '#2196F3',
  },
  denunciaCard: {
    backgroundColor: '#fff',
    padding: 15,
    margin: 10,
    borderRadius: 8,
    elevation: 2,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
  },
  denunciaHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 10,
  },
  numeroDenuncia: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#333',
  },
  badge: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
  },
  badgeText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
  motivo: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    marginBottom: 5,
  },
  denunciante: {
    fontSize: 12,
    color: '#666',
    marginBottom: 5,
  },
  fecha: {
    fontSize: 12,
    color: '#999',
  },
});

export default ListadoDenunciasScreen;
```

### 2. Pantalla de Detalle de Denuncia

```javascript
import React, { useState, useEffect } from 'react';
import { View, Text, ScrollView, StyleSheet, ActivityIndicator, TouchableOpacity } from 'react-native';

const DetalleDenunciaScreen = ({ route, navigation }) => {
  const { id } = route.params;
  const [denuncia, setDenuncia] = useState(null);
  const [loading, setLoading] = useState(true);
  const token = 'TU_TOKEN_MASTER_AQUI';

  useEffect(() => {
    cargarDenuncia();
  }, []);

  const cargarDenuncia = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `https://tu-api.com/api/denuncia/obtener?id=${id}`,
        {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      const data = await response.json();
      if (data.success) {
        setDenuncia(data.data);
      }
    } catch (error) {
      console.error('Error al cargar denuncia:', error);
    } finally {
      setLoading(false);
    }
  };

  const actualizarEstado = async (nuevoEstado) => {
    try {
      const response = await fetch(
        'https://tu-api.com/api/denuncia/actualizar-estado',
        {
          method: 'PUT',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            id: id,
            estado: nuevoEstado,
          }),
        }
      );

      const data = await response.json();
      if (data.success) {
        // Recargar la denuncia
        cargarDenuncia();
        alert('Estado actualizado correctamente');
      }
    } catch (error) {
      console.error('Error al actualizar estado:', error);
      alert('Error al actualizar el estado');
    }
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#2196F3" />
      </View>
    );
  }

  if (!denuncia) {
    return (
      <View style={styles.center}>
        <Text>Denuncia no encontrada</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Número de Denuncia</Text>
        <Text style={styles.value}>{denuncia.numero_denuncia}</Text>
      </View>

      <View style={styles.card}>
        <Text style={styles.title}>Estado</Text>
        <View style={[styles.badge, getBadgeColor(denuncia.estado)]}>
          <Text style={styles.badgeText}>{denuncia.estado}</Text>
        </View>
      </View>

      <View style={styles.card}>
        <Text style={styles.sectionTitle}>Datos del Denunciante</Text>
        <Text style={styles.label}>Nombre:</Text>
        <Text style={styles.value}>{denuncia.nombre_denunciante}</Text>
        <Text style={styles.label}>Cédula:</Text>
        <Text style={styles.value}>{denuncia.cedula_denunciante}</Text>
        <Text style={styles.label}>Teléfono:</Text>
        <Text style={styles.value}>{denuncia.telefono_denunciante}</Text>
        {denuncia.parroquia_denunciante && (
          <>
            <Text style={styles.label}>Parroquia:</Text>
            <Text style={styles.value}>{denuncia.parroquia_denunciante}</Text>
          </>
        )}
      </View>

      <View style={styles.card}>
        <Text style={styles.sectionTitle}>Datos de la Denuncia</Text>
        <Text style={styles.label}>Motivo:</Text>
        <Text style={styles.value}>{denuncia.motivo}</Text>
        <Text style={styles.label}>Hechos:</Text>
        <Text style={styles.value}>{denuncia.hechos}</Text>
        <Text style={styles.label}>Fecha:</Text>
        <Text style={styles.value}>
          {new Date(denuncia.fecha_denuncia).toLocaleString('es-VE')}
        </Text>
      </View>

      {denuncia.nombre_denunciado && (
        <View style={styles.card}>
          <Text style={styles.sectionTitle}>Datos del Denunciado</Text>
          <Text style={styles.label}>Nombre:</Text>
          <Text style={styles.value}>{denuncia.nombre_denunciado}</Text>
          {denuncia.direccion_denunciado && (
            <>
              <Text style={styles.label}>Dirección:</Text>
              <Text style={styles.value}>{denuncia.direccion_denunciado}</Text>
            </>
          )}
        </View>
      )}

      {/* Botones para cambiar estado */}
      <View style={styles.acciones}>
        <Text style={styles.accionesTitle}>Cambiar Estado</Text>
        <View style={styles.botonesEstado}>
          <TouchableOpacity
            style={[styles.btnEstado, denuncia.estado === 'Pendiente' && styles.btnEstadoActive]}
            onPress={() => actualizarEstado('Pendiente')}
          >
            <Text style={styles.btnEstadoText}>Pendiente</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.btnEstado, denuncia.estado === 'En Proceso' && styles.btnEstadoActive]}
            onPress={() => actualizarEstado('En Proceso')}
          >
            <Text style={styles.btnEstadoText}>En Proceso</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.btnEstado, denuncia.estado === 'Resuelta' && styles.btnEstadoActive]}
            onPress={() => actualizarEstado('Resuelta')}
          >
            <Text style={styles.btnEstadoText}>Resuelta</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.btnEstado, denuncia.estado === 'Archivada' && styles.btnEstadoActive]}
            onPress={() => actualizarEstado('Archivada')}
          >
            <Text style={styles.btnEstadoText}>Archivada</Text>
          </TouchableOpacity>
        </View>
      </View>
    </ScrollView>
  );
};

const getBadgeColor = (estado) => {
  const colors = {
    'Pendiente': { backgroundColor: '#FFA500' },
    'En Proceso': { backgroundColor: '#2196F3' },
    'Resuelta': { backgroundColor: '#4CAF50' },
    'Archivada': { backgroundColor: '#9E9E9E' },
  };
  return colors[estado] || { backgroundColor: '#9E9E9E' };
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  card: {
    backgroundColor: '#fff',
    padding: 15,
    margin: 10,
    borderRadius: 8,
    elevation: 2,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#333',
    marginBottom: 10,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#2196F3',
    marginBottom: 15,
  },
  label: {
    fontSize: 12,
    color: '#666',
    marginTop: 10,
    marginBottom: 5,
  },
  value: {
    fontSize: 14,
    color: '#333',
  },
  badge: {
    alignSelf: 'flex-start',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 15,
  },
  badgeText: {
    color: '#fff',
    fontSize: 14,
    fontWeight: 'bold',
  },
  acciones: {
    backgroundColor: '#fff',
    padding: 15,
    margin: 10,
    borderRadius: 8,
    elevation: 2,
  },
  accionesTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  botonesEstado: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 10,
  },
  btnEstado: {
    paddingHorizontal: 15,
    paddingVertical: 10,
    borderRadius: 8,
    backgroundColor: '#e0e0e0',
    marginRight: 10,
    marginBottom: 10,
  },
  btnEstadoActive: {
    backgroundColor: '#2196F3',
  },
  btnEstadoText: {
    color: '#333',
    fontWeight: '600',
  },
});

export default DetalleDenunciaScreen;
```

### 3. Pantalla "Mis Denuncias" (Para Ciudadanos)

```javascript
import React, { useState, useEffect } from 'react';
import { View, Text, FlatList, TouchableOpacity, StyleSheet, ActivityIndicator } from 'react-native';

const MisDenunciasScreen = ({ navigation }) => {
  const [denuncias, setDenuncias] = useState([]);
  const [loading, setLoading] = useState(true);
  const token = 'TU_TOKEN_CIUDADANO_AQUI'; // Obtener del contexto/AsyncStorage

  useEffect(() => {
    cargarMisDenuncias();
  }, []);

  const cargarMisDenuncias = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        'https://tu-api.com/api/denuncia/mis-denuncias',
        {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      const data = await response.json();
      if (data.success) {
        setDenuncias(data.data);
      }
    } catch (error) {
      console.error('Error al cargar denuncias:', error);
    } finally {
      setLoading(false);
    }
  };

  const renderDenuncia = ({ item }) => (
    <TouchableOpacity
      style={styles.denunciaCard}
      onPress={() => navigation.navigate('DetalleMiDenuncia', { id: item.id })}
    >
      <View style={styles.denunciaHeader}>
        <Text style={styles.numeroDenuncia}>{item.numero_denuncia}</Text>
        <View style={[styles.badge, getBadgeColor(item.estado)]}>
          <Text style={styles.badgeText}>{item.estado}</Text>
        </View>
      </View>
      <Text style={styles.motivo}>{item.motivo}</Text>
      <Text style={styles.fecha}>
        {new Date(item.fecha_denuncia).toLocaleDateString('es-VE')}
      </Text>
    </TouchableOpacity>
  );

  const getBadgeColor = (estado) => {
    const colors = {
      'Pendiente': { backgroundColor: '#FFA500' },
      'En Proceso': { backgroundColor: '#2196F3' },
      'Resuelta': { backgroundColor: '#4CAF50' },
      'Archivada': { backgroundColor: '#9E9E9E' },
    };
    return colors[estado] || { backgroundColor: '#9E9E9E' };
  };

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" color="#2196F3" />
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <FlatList
        data={denuncias}
        renderItem={renderDenuncia}
        keyExtractor={(item) => item.id}
        ListEmptyComponent={
          <View style={styles.center}>
            <Text>No has realizado ninguna denuncia</Text>
          </View>
        }
        refreshing={loading}
        onRefresh={cargarMisDenuncias}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  center: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  denunciaCard: {
    backgroundColor: '#fff',
    padding: 15,
    margin: 10,
    borderRadius: 8,
    elevation: 2,
  },
  denunciaHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 10,
  },
  numeroDenuncia: {
    fontSize: 16,
    fontWeight: 'bold',
    color: '#333',
  },
  badge: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
  },
  badgeText: {
    color: '#fff',
    fontSize: 12,
    fontWeight: 'bold',
  },
  motivo: {
    fontSize: 14,
    fontWeight: '600',
    color: '#333',
    marginBottom: 5,
  },
  fecha: {
    fontSize: 12,
    color: '#999',
  },
});

export default MisDenunciasScreen;
```

## 🔐 Autenticación

### Para Ciudadanos
- Obtener el token después del login en `/api/ciudadano/login`
- Guardar el token en AsyncStorage o Context
- Incluir el token en el header `Authorization: Bearer {token}`

### Para Usuarios del Sistema
- Obtener el token después del login en `/api/master/login`
- Verificar que el usuario tenga el permiso `"denuncias"` en su array de permisos
- Guardar el token en AsyncStorage o Context
- Incluir el token en el header `Authorization: Bearer {token}`

## 📱 Navegación

Agregar las pantallas a tu navegador:

```javascript
// Ejemplo con React Navigation
import { createStackNavigator } from '@react-navigation/stack';

const Stack = createStackNavigator();

function DenunciasStack() {
  return (
    <Stack.Navigator>
      <Stack.Screen 
        name="ListadoDenuncias" 
        component={ListadoDenunciasScreen}
        options={{ title: 'Denuncias' }}
      />
      <Stack.Screen 
        name="DetalleDenuncia" 
        component={DetalleDenunciaScreen}
        options={{ title: 'Detalle de Denuncia' }}
      />
    </Stack.Navigator>
  );
}
```

## ✅ Checklist de Implementación

- [ ] Crear servicio/API helper para las llamadas a los endpoints
- [ ] Implementar pantalla de listado de denuncias (para usuarios del sistema)
- [ ] Implementar pantalla de detalle de denuncia
- [ ] Implementar funcionalidad de cambio de estado
- [ ] Implementar pantalla "Mis Denuncias" (para ciudadanos)
- [ ] Agregar filtros por estado
- [ ] Implementar paginación
- [ ] Manejar errores y estados de carga
- [ ] Agregar navegación entre pantallas
- [ ] Probar con diferentes tokens (ciudadano y master)

## 🎯 Notas Importantes

1. **Permisos**: Solo los usuarios master con el permiso `"denuncias"` pueden ver todas las denuncias y cambiar su estado.

2. **Estados**: Los estados válidos son:
   - `"Pendiente"` - Denuncia recién creada
   - `"En Proceso"` - Denuncia siendo procesada
   - `"Resuelta"` - Denuncia resuelta
   - `"Archivada"` - Denuncia archivada

3. **Paginación**: El endpoint de listado soporta paginación. Implementa scroll infinito o paginación manual según tu preferencia.

4. **Filtros**: Puedes filtrar las denuncias por estado usando el parámetro `estado` en la query string.

5. **Formato de Fechas**: Las fechas vienen en formato ISO 8601 (RFC3339). Usa `new Date()` para parsearlas.

---

**Estado**: ✅ Backend listo, Frontend pendiente de implementación
**Fecha**: 2025-01-15

