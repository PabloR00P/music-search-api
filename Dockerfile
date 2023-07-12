# Imagen base de Go
FROM golang:1.16-alpine as builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Compilar la aplicación
RUN go build -o main .

# Imagen base de PostgreSQL
FROM postgres:13-alpine

# Copiar el archivo SQL de inicialización a la carpeta de scripts de PostgreSQL
COPY init.sql /docker-entrypoint-initdb.d/

# Exponer el puerto 8000 para la aplicación Go
EXPOSE 8000

# Copiar el archivo ejecutable de la aplicación Go desde el primer contenedor al segundo contenedor
COPY --from=builder /app/main /

# Establecer el comando de inicio para ejecutar la aplicación Go
CMD ["./main"]
