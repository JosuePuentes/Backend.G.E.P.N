package routes

import (
	"gepn/handlers"
	"net/http"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Rutas públicas
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/api/health", handlers.HealthHandler)
	mux.HandleFunc("/ciudadano", handlers.CiudadanoHandler)
	mux.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	// Rutas de autenticación
	mux.HandleFunc("/api/policial/login", handlers.LoginPolicialHandler)
	mux.HandleFunc("/api/policial/finalizar-guardia", handlers.FinalizarGuardiaHandler)
	
	// Rutas de ciudadanos (públicas)
	mux.HandleFunc("/api/ciudadano/registro", handlers.RegistroCiudadanoHandler)
	mux.HandleFunc("/api/ciudadano/login", handlers.LoginCiudadanoHandler)
	
	// Rutas de denuncias (requieren autenticación de ciudadano)
	mux.HandleFunc("/api/denuncia/crear", handlers.CrearDenunciaHandler)
	mux.HandleFunc("/api/denuncia/mis-denuncias", handlers.MisDenunciasHandler)
	
	// Rutas de denuncias para usuarios del sistema (requieren autenticación de master con permiso "denuncias")
	mux.HandleFunc("/api/denuncia/listar", handlers.ListarTodasDenunciasHandler)
	mux.HandleFunc("/api/denuncia/obtener", handlers.ObtenerDenunciaHandler)
	mux.HandleFunc("/api/denuncia/actualizar-estado", handlers.ActualizarEstadoDenunciaHandler)

	// Rutas protegidas - Detenidos
	mux.HandleFunc("/api/detenidos", handlers.CrearDetenidoHandler)
	mux.HandleFunc("/api/detenidos/listar", handlers.ListarDetenidosHandler)
	mux.HandleFunc("/api/detenidos/obtener", handlers.ObtenerDetenidoHandler)

	// Rutas protegidas - Minutas
	mux.HandleFunc("/api/minutas", handlers.CrearMinutaHandler)
	mux.HandleFunc("/api/minutas/listar", handlers.ListarMinutasHandler)
	mux.HandleFunc("/api/minutas/obtener", handlers.ObtenerMinutaHandler)

	// Rutas protegidas - Búsqueda
	mux.HandleFunc("/api/buscar/cedula", handlers.BuscarCedulaHandler)
	mux.HandleFunc("/api/mas-buscados", handlers.ListarMasBuscadosHandler)

	// Rutas protegidas - Pánico
	mux.HandleFunc("/api/panico/activar", handlers.ActivarPanicoHandler)
	mux.HandleFunc("/api/panico/alertas", handlers.ListarAlertasPanicoHandler)

	// Rutas públicas de Master - Login y módulos
	mux.HandleFunc("/api/master/login", handlers.LoginMasterHandler)
	mux.HandleFunc("/api/master/modulos", handlers.ListarModulosHandler)
	mux.HandleFunc("/api/master/inicializar-admin", handlers.InicializarAdminHandler) // Temporal para crear admin
	mux.HandleFunc("/api/master/resetear-password-admin", handlers.ResetearPasswordAdminHandler) // Temporal para resetear contraseña

	// Rutas protegidas de Master (requieren autenticación)
	mux.HandleFunc("/api/master/crear-usuario", handlers.CrearUsuarioMasterHandler)
	mux.HandleFunc("/api/master/usuarios", handlers.ListarUsuariosMasterHandler)
	mux.HandleFunc("/api/master/usuarios/permisos/", handlers.ActualizarPermisosHandler)
	mux.HandleFunc("/api/master/usuarios/activar/", handlers.ActivarUsuarioMasterHandler)
	mux.HandleFunc("/api/master/verificar", handlers.VerificarMasterHandler)

	// Rutas de RRHH (requieren autenticación como master con permiso rrhh)
	mux.HandleFunc("/api/rrhh/registrar-oficial", handlers.RegistrarOficialHandler)
	mux.HandleFunc("/api/rrhh/generar-qr/", handlers.GenerarQRHandler)
	mux.HandleFunc("/api/rrhh/verificar-qr/", handlers.VerificarQRHandler)
	mux.HandleFunc("/api/rrhh/listar-oficiales", handlers.ListarOficialesHandler)
	mux.HandleFunc("/api/rrhh/ascensos-pendientes", handlers.AscensosPendientesHandler)
	mux.HandleFunc("/api/rrhh/aprobar-ascenso/", handlers.AprobarAscensoHandler)

	// Rutas de Centro de Coordinación (requieren autenticación como master)
	mux.HandleFunc("/api/centro-coordinacion/centros", handlers.CentrosHandler)
	mux.HandleFunc("/api/centro-coordinacion/estaciones", handlers.EstacionesHandler)
	mux.HandleFunc("/api/centro-coordinacion/estaciones/asignar", handlers.AsignarFuncionarioHandler)
	mux.HandleFunc("/api/centro-coordinacion/estaciones/funcionarios", handlers.ListarFuncionariosEstacionHandler)
	mux.HandleFunc("/api/centro-coordinacion/partes", handlers.PartesHandler)

	// Rutas de Patrullaje
	mux.HandleFunc("/api/patrullaje/login", handlers.LoginPatrullajeHandler)
	mux.HandleFunc("/api/patrullaje/crear-usuario-prueba", handlers.CrearUsuarioPruebaPatrullajeHandler)
	mux.HandleFunc("/api/patrullaje/iniciar", handlers.IniciarPatrullajeHandler)
	mux.HandleFunc("/api/patrullaje/actualizar-ubicacion", handlers.ActualizarUbicacionHandler)
	mux.HandleFunc("/api/patrullaje/activos", handlers.ObtenerPatrullajesActivosHandler)
	mux.HandleFunc("/api/patrullaje/finalizar", handlers.FinalizarPatrullajeHandler)
	mux.HandleFunc("/api/patrullaje/historial", handlers.HistorialPatrullajesHandler)

	return mux
}

