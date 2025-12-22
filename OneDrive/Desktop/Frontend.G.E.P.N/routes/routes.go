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

	// Rutas públicas de Master (RRHH) - Registro y Login
	mux.HandleFunc("/api/rrhh/master/registro", handlers.RegistroMasterHandler)
	mux.HandleFunc("/api/rrhh/master/login", handlers.LoginMasterHandler)
	mux.HandleFunc("/api/rrhh/master/verificar", handlers.VerificarMasterHandler)

	// Rutas de RRHH (requieren autenticación como master)
	mux.HandleFunc("/api/rrhh/registrar-oficial", handlers.RegistrarOficialHandler)
	mux.HandleFunc("/api/rrhh/generar-qr/", handlers.GenerarQRHandler)
	mux.HandleFunc("/api/rrhh/verificar-qr/", handlers.VerificarQRHandler)
	mux.HandleFunc("/api/rrhh/listar-oficiales", handlers.ListarOficialesHandler)
	mux.HandleFunc("/api/rrhh/ascensos-pendientes", handlers.AscensosPendientesHandler)
	mux.HandleFunc("/api/rrhh/aprobar-ascenso/", handlers.AprobarAscensoHandler)

	return mux
}

