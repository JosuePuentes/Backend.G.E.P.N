# Subir cambios del backend para que Render despliegue

## Archivos que debes incluir en el commit

- `database/repositories.go` – Usuario PATRULLA-TEST se crea al iniciar el backend
- `handlers/patrullaje.go` – Endpoint crear-usuario-prueba
- `middleware/middleware.go` – CORS: OPTIONS responde 200
- `routes/routes.go` – Ruta `/api/patrullaje/crear-usuario-prueba`

## Comandos (ejecutar en la raíz del repo de backend)

Si tu repo de backend es **Frontend.G.E.P.N** (monorepo):

```bash
cd c:\Users\puent\OneDrive\Desktop\Frontend.G.E.P.N
git add database/repositories.go handlers/patrullaje.go middleware/middleware.go routes/routes.go
git commit -m "Backend: CORS OPTIONS 200, usuario PATRULLA-TEST en BD, endpoint crear-usuario-prueba"
git push origin main
```

Si el backend está en **otro repo** (por ejemplo Backend.G.E.P.N):

1. Copia los 4 archivos modificados al repo del backend.
2. En ese repo:
   ```bash
   git add database/repositories.go handlers/patrullaje.go middleware/middleware.go routes/routes.go
   git commit -m "Backend: CORS OPTIONS 200, usuario PATRULLA-TEST en BD, endpoint crear-usuario-prueba"
   git push origin main
   ```

## Después del push

- Si Render está conectado a ese repo y a la rama `main`, el despliegue se iniciará solo.
- Tras el deploy, el usuario **PATRULLA-TEST** / **PIN 123456** se creará en la base de datos al arrancar el servicio.
- Entra en https://frontendgepn.vercel.app/patrullaje/login con esas credenciales.
