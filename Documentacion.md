# Documentación del servicio de búsqueda de canciones

El servicio de búsqueda de canciones es una aplicación desarrollada en Go que permite buscar canciones en iTunes y ChartLyrics. La aplicación se ejecuta en un entorno Docker.

## Requisitos

Para utilizar el servicio de búsqueda de canciones, necesitarás tener instalados los siguientes componentes en tu sistema:

- Docker

## Configuración

Antes de ejecutar el servicio, es necesario realizar la siguiente configuración:

1. Abre el archivo `docker-compose.yml` y asegúrate de que los puertos `8000` y `5432` no estén siendo utilizados por otros servicios en tu sistema. Si están en uso, puedes modificar los puertos según tus preferencias.

2. Opcionalmente, si deseas cambiar el nombre de la base de datos, el usuario o la contraseña, puedes modificar las variables de entorno en el archivo `docker-compose.yml` en la sección `db` bajo `environment`.

## Ejecución del servicio

Una vez que hayas completado la configuración, puedes ejecutar el servicio siguiendo estos pasos:

1. Abre una terminal y navega hasta la ubicación del proyecto.

2. Ejecuta el siguiente comando para construir las imágenes Docker y levantar los contenedores:

   ```shell
   docker-compose up
   ```

   Este comando descargará las imágenes necesarias, creará los contenedores y ejecutará la aplicación en el puerto `8000`.

3. Una vez que los contenedores se hayan levantado correctamente, verás un mensaje que indica que el servidor se está ejecutando en `http://localhost:8000`.

## Uso del servicio

El servicio de búsqueda de canciones proporciona un punto final `/search` al que puedes realizar una solicitud HTTP GET para buscar canciones.

### Parámetros de la solicitud

La solicitud GET debe incluir los siguientes parámetros en la URL:

- `name` (obligatorio): El nombre de la canción que deseas buscar.
- `artist` (obligatorio): El nombre del artista de la canción.
- `album` (opcional): El nombre del álbum al que pertenece la canción.

Ejemplo de URL de solicitud de búsqueda:

```
http://localhost:8000/search?name=antologia&artist=shakira&album=pies
```

### Respuesta

El servicio responderá con una lista de canciones que coincidan con los criterios de búsqueda. La respuesta estará en formato JSON y tendrá la siguiente estructura:

```json
{
  "results": [
    {
        "id": "1411628233",
        "name": "Believer",
        "artist": "Imagine Dragons",
        "duration": "3:24",
        "album": "Evolve",
        "artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music126/v4/11/7a/b8/117ab805-6811-8929-18b9-0fad7baf0c25/17UMGIM98210.rgb.jpg/100x100bb.jpg",
        "price": "GTQ 1.29",
        "origin": "iTunes"
    }
  ]
}
```

Cada objeto de canción en la lista de resultados contendrá los siguientes campos:

- `id`: El identificador único de la canción.
- `name`: El nombre de la canción.
- `artist`: El nombre del artista de la canción.
- `duration`: La duración de la canción en formato minutos:segundos. Este campo puede estar vacío en algunas canciones.
- `album`: El nombre del álbum al que pertenece la canción. Este campo puede estar vacío en algunas canciones.
- `artwork`: La URL de la imagen de portada de la canción. Este campo puede estar vacío en algunas canciones.
- `price`: El precio de la canción en formato monetario. Este campo puede estar vacío en algunas canciones.
- `origin`: El origen de la canción, que puede ser "iTunes" o "ChartLyrics".

## Almacenamiento de registros

El servicio almacena los registros de las canciones en una base de datos PostgreSQL. Puedes acceder a la base de datos utilizando las siguientes credenciales:

- Nombre de usuario: `my-user`
- Contraseña: `my-password`
- Nombre de la base de datos: `my-database`
- Host: `localhost`
- Puerto: `5432`

La tabla `songs` contiene los siguientes campos:

- `id`: El identificador único de la canción (clave primaria).
- `name`: El nombre de la canción.
- `artist`: El nombre del artista de la canción.
- `duration`: La duración de la canción en formato minutos:segundos. Este campo puede estar vacío en algunas canciones.
- `album`: El nombre del álbum al que pertenece la canción. Este campo puede estar vacío en algunas canciones.
- `artwork`: La URL de la imagen de portada de la canción. Este campo puede estar vacío en algunas canciones.
- `price`: El precio de la canción en formato monetario. Este campo puede estar vacío en algunas canciones.
- `origin`: El origen de la canción, que puede ser "iTunes" o "ChartLyrics".

## Detener el servicio

Para detener el servicio, puedes presionar `Ctrl + C` en la terminal donde se esté ejecutando el comando `docker-compose up`. Esto detendrá los contenedores y liberará los puertos utilizados.

## Notas adicionales

- El servicio de búsqueda de canciones utiliza el servicio de iTunes y ChartLyrics para obtener información sobre las canciones. Asegúrate de tener una conexión a Internet activa para que el servicio funcione correctamente.
