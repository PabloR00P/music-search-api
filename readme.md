
# Servicio de Búsqueda de Canciones

El servicio de búsqueda de canciones es una aplicación desarrollada en Go que permite buscar canciones en iTunes y ChartLyrics. Utiliza Docker para facilitar su ejecución y configuración.

## Instalación

Antes de comenzar, asegúrate de tener instalados los siguientes componentes en tu sistema:

- Docker

## Configuración

Antes de ejecutar el servicio, hay algunas configuraciones que debes realizar:

1. Abre el archivo `docker-compose.yml` y revisa los puertos y variables de entorno utilizados por el servicio. Puedes personalizarlos según tus preferencias.

2. Opcionalmente, puedes modificar el archivo `nginx.conf` para ajustar la configuración del servidor Nginx utilizado en el servicio.

## Inicio Rápido

Sigue los pasos a continuación para ejecutar el servicio de búsqueda de canciones:

1. Clona el repositorio en tu máquina local:

   ```bash
   git clone https://github.com/PabloR00P/music-search-api.git
   ```

2. Navega hasta la carpeta del proyecto:

   ```bash
   cd music-search-api
   ```

3. Ir hacia la rama del proyecto music-search-api

   ```bash
   git checkout project/music-search-api
   ```

4. Ejecuta el siguiente comando para construir las imágenes Docker y levantar los contenedores:

   ```bash
   docker-compose up
   ```

   Esto descargará las imágenes necesarias, creará los contenedores y ejecutará la aplicación en el puerto 8000.
   

5. Una vez que los contenedores se hayan levantado correctamente, podrás acceder al servicio de búsqueda de canciones en tu navegador web en la siguiente URL:

   ```
   http://localhost:8000
   ```

## Uso del Servicio

El servicio de búsqueda de canciones proporciona un end point `/search` al que puedes realizar una solicitud GET para buscar canciones. Puedes proporcionar los siguientes parámetros en la URL:

- `name` (obligatorio): El nombre de la canción que deseas buscar.
- `artist` (obligatorio): El nombre del artista de la canción.
- `album` (opcional): El nombre del álbum al que pertenece la canción.

Aquí tienes un ejemplo de solicitud de búsqueda utilizando `curl`:

```bash
curl "http://localhost:8000/search?name=antologia&artist=shakira&album=pies"
```

La respuesta será una lista de canciones que coincidan con los criterios de búsqueda en formato JSON.

## Detener el Servicio

Para detener el servicio, puedes ejecutar el siguiente comando en la carpeta del proyecto:

```bash
docker-compose down
```

Esto detendrá los contenedores y liberará los puertos utilizados.



## Licencia

[Licencia MIT](LICENSE).
