# Simple Toolkit

`simple-toolkit` es una librería de utilidades para acelerar el desarrollo en Go, proporcionando módulos reutilizables para **logging, clientes HTTP/SQS, Base de datos Dynamo y Redis, manejo de errores, configuración**, y más.


[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=uala-challenge_simple-toolkit&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=uala-challenge_simple-toolkit)
![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![Viper](https://img.shields.io/badge/configuration-viper-green.svg)

## **Características Principales**
* Modularidad y reutilización  
* Integraciones con servicios en la nube (AWS, Redis, etc.)  
* Simplificación de tareas comunes en microservicios  
* Código optimizado y fácil de extender

---

## **Módulos Disponibles**

### **Logging (`utilities/log`)**
Manejo de logs centralizado con soporte para diferentes niveles de severidad.

 **Funcionalidades:**
- Formato estructurado (`JSON`)
- Soporte para logs en consola y servicios externos
- Diferentes niveles: `INFO`, `ERROR`, `WARN`, `DEBUG`

 **Ejemplo de Uso:**
```go
log := log.NewService()
log.Info(ctx, "Iniciando aplicación", nil)
log.Error(ctx, err, "Ocurrió un error crítico", nil)
```

### **Cliente HTTP (client/rest)**

Cliente HTTP con timeouts, retries y manejo de errores mejorado.

**Ejemplo de Uso:**
```go
client := rest.NewClient(rest.Config{
    Timeout: 5 * time.Second,
})
resp, err := client.Get(ctx, "https://api.example.com/data", nil)
```
### **Cliente SQS (client/sqs)**
Interfaz para interactuar con AWS SQS de forma eficiente.
```go
sqsClient := sqs.NewClient()
msg := sqs.SendMessageInput{
    QueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue",
    MessageBody: "Hello SQS",
}
err := sqsClient.SendMessage(ctx, msg)
```
### **Cliente SNS (client/sns)**
Publicación de mensajes en AWS SNS con retries.
```go
snsClient := sns.NewClient()
msg := sns.PublishInput{
    Message: "Nuevo evento",
    TopicArn: "arn:aws:sns:us-east-1:123456789012:MyTopic",
}
err := snsClient.Publish(ctx, msg)
```
### **Manejo de Errores (utilities/error_handler)**
Proporciona una forma estándar de manejar y estructurar errores.
```go

func (s service) Init(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("bad request", err.Error(), err, http.StatusBadRequest), w)
		return
	}....
}
```
### **Configuración (utilities/config)**
Carga y parseo de variables de entorno y archivos de configuración.
```go
config := config.NewConfig()
config.LoadEnv()
config.LoadFile("config.yaml")
``` 

### **Instalación**
```bash
go get github.com/uala-challenge/simple-toolkit
```
