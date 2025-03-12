# simple_router

**simple_router** es un paquete en Go diseñado para simplificar la configuración de rutas y la creación de servidores web. Utiliza el router **Chi** para gestionar rutas y proporciona una estructura modular y flexible para integrar servicios adicionales como **Swagger**, **Docsify** y funcionalidades básicas como un **ping endpoint**. Está pensado para ser fácil de configurar y extender, adaptándose a diferentes perfiles de la aplicación (desarrollo, prueba, producción).

## Características

- **Router con Chi**: Utiliza **Chi** para la gestión de rutas HTTP de forma eficiente y ligera.
- **Ping Endpoint**: Incluye un endpoint de **ping** (`/ping`) para verificar la disponibilidad del servicio.
- **Swagger**: Soporta la integración de **Swagger** para la documentación de la API, accesible en entornos de prueba o desarrollo.
- **Docsify**: Proporciona integración con **Docsify** para generar documentación estática de la API, accesible a través de rutas específicas.
- **Configuración de puerto flexible**: Permite establecer un puerto personalizado para el servidor o utilizar el puerto predeterminado.

## Instalación

Para utilizar el paquete `simple_router`, simplemente incluye el módulo en tu proyecto:

```bash
go get github.com/skolldire/web-simplify/pkg/utilities/simple_router
```

## Uso

```go

package main

import (
	"github.com/skolldire/web-simplify/pkg/utilities/simple_router"
	"log"
)

func main() {
	config := simple_router.Config{
		Port: "8080", // O dejar vacío para usar el puerto por defecto
	}
	service := simple_router.NewService(config)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

```
