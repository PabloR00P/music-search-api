package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Estructura para representar una canción
type Song struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Duration string `json:"duration"`
	Album    string `json:"album"`
	Artwork  string `json:"artwork"`
	Price    string `json:"price"`
	Origin   string `json:"origin"`
}

// Estructura para representar la respuesta de la búsqueda de canciones
type SearchResponse struct {
	Results []Song `json:"results"`
	Message string `json:"message,omitempty"`
}

// Función para buscar canciones en iTunes
func searchIniTunes(name, artist, album string) ([]Song, error) {
	// Construir la URL de búsqueda en iTunes
	baseURL := "https://itunes.apple.com/search"
	params := url.Values{
		"term":  {name},
		"media": {"music"},
		"music": {"song"},
	}
	if artist != "" {
		params.Set("artist", artist)
	}
	if album != "" {
		params.Set("album", album)
	}
	searchURL := baseURL + "?" + params.Encode()

	// Realizar la solicitud HTTP a iTunes
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta JSON de iTunes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decodificar la respuesta JSON
	var iTunesResponse struct {
		Results []struct {
			TrackID    int     `json:"trackId"`
			TrackName  string  `json:"trackName"`
			ArtistName string  `json:"artistName"`
			Duration   int     `json:"trackTimeMillis"`
			AlbumName  string  `json:"collectionName"`
			ArtworkURL string  `json:"artworkUrl100"`
			Price      float64 `json:"trackPrice"`
		} `json:"results"`
	}
	err = json.Unmarshal(body, &iTunesResponse)
	if err != nil {
		return nil, err
	}

	// Convertir los resultados a la estructura Song
	var songs []Song
	for _, result := range iTunesResponse.Results {
		song := Song{
			ID:       strconv.Itoa(result.TrackID),
			Name:     result.TrackName,
			Artist:   result.ArtistName,
			Duration: formatDuration(result.Duration),
			Album:    result.AlbumName,
			Artwork:  result.ArtworkURL,
			Price:    fmt.Sprintf("GTQ %.2f", result.Price),
			Origin:   "iTunes",
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// Función para buscar canciones en ChartLyrics
func searchChartLyrics(name, artist, album string) ([]Song, error) {
	// Construir la URL de búsqueda en ChartLyrics
	baseURL := "http://api.chartlyrics.com/apiv1.asmx/SearchLyric"
	params := url.Values{
		"artist": {artist},
		"song":   {name},
	}
	searchURL := baseURL + "?" + params.Encode()

	// Realizar la solicitud HTTP a ChartLyrics
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Leer la respuesta XML de ChartLyrics
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Procesar la respuesta XML de ChartLyrics y convertirla en canciones
	songs, err := processChartLyricsResponse(body)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

// Función para procesar la respuesta XML de ChartLyrics y convertirla a canciones
func processChartLyricsResponse(body []byte) ([]Song, error) {
	// Decodificar la respuesta XML
	var chartLyricsResponse struct {
		TrackID []string `xml:"SearchLyricResult>TrackId"`
		Song    []string `xml:"SearchLyricResult>Song"`
		Artist  []string `xml:"SearchLyricResult>Artist"`
	}

	err := xml.Unmarshal(body, &chartLyricsResponse)
	if err != nil {
		return nil, err
	}

	// Convertir los resultados a la estructura Song
	var songs []Song
	for i := 0; i < len(chartLyricsResponse.TrackID); i++ {
		song := Song{
			ID:       chartLyricsResponse.TrackID[i],
			Name:     chartLyricsResponse.Song[i],
			Artist:   chartLyricsResponse.Artist[i],
			Duration: "",
			Album:    "",
			Artwork:  "",
			Price:    "",
			Origin:   "ChartLyrics",
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// Función para formatear la duración en milisegundos a formato minutos:segundos
func formatDuration(duration int) string {
	seconds := duration / 1000
	minutes := seconds / 60
	seconds %= 60

	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// Función para almacenar la respuesta en la base de datos
func storeResponseInDB(songs []Song) error {
	// Abrir la conexión a la base de datos
	db, err := sql.Open("postgres", "postgres://my-user:my-password@db:5432/my-database?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Crear la tabla si no existe
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS songs (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		artist TEXT NOT NULL,
		duration TEXT,
		album TEXT,
		artwork TEXT,
		price TEXT,
		origin TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	// Borrar los registros anteriores
	_, err = db.Exec("DELETE FROM songs")
	if err != nil {
		return err
	}

	// Insertar los registros de canciones
	for _, song := range songs {
		_, err = db.Exec("INSERT INTO songs (name, artist, duration, album, artwork, price, origin) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			song.Name, song.Artist, song.Duration, song.Album, song.Artwork, song.Price, song.Origin)
		if err != nil {
			return err
		}
	}

	return nil
}

// Función para obtener los registros almacenados en la base de datos
func getStoredSongsFromDB() ([]Song, error) {
	// Abrir la conexión a la base de datos
	db, err := sql.Open("postgres", "postgres://my-user:my-password@db:5432/my-database?sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Ejecutar la consulta SQL para obtener los registros almacenados
	rows, err := db.Query("SELECT id, name, artist, COALESCE(duration, ''), COALESCE(album, ''), COALESCE(artwork, ''), COALESCE(price, ''), origin FROM songs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Recorrer los resultados y crear la lista de canciones almacenadas
	var storedSongs []Song
	for rows.Next() {
		var id, name, artist, duration, album, artwork, price, origin string
		err := rows.Scan(&id, &name, &artist, &duration, &album, &artwork, &price, &origin)
		if err != nil {
			return nil, err
		}
		song := Song{
			ID:       id,
			Name:     name,
			Artist:   artist,
			Duration: duration,
			Album:    album,
			Artwork:  artwork,
			Price:    price,
			Origin:   origin,
		}
		storedSongs = append(storedSongs, song)
	}

	return storedSongs, nil
}

// Función de comparación para ordenar las canciones según la coincidencia con los criterios de búsqueda
func compareSongs(a, b *Song, name, artist, album string) bool {
	scoreA := 0
	scoreB := 0

	if strings.Contains(strings.ToLower(a.Name), strings.ToLower(name)) {
		scoreA++
	}
	if strings.Contains(strings.ToLower(a.Artist), strings.ToLower(artist)) {
		scoreA++
	}
	if strings.Contains(strings.ToLower(a.Album), strings.ToLower(album)) {
		scoreA++
	}

	if strings.Contains(strings.ToLower(b.Name), strings.ToLower(name)) {
		scoreB++
	}
	if strings.Contains(strings.ToLower(b.Artist), strings.ToLower(artist)) {
		scoreB++
	}
	if strings.Contains(strings.ToLower(b.Album), strings.ToLower(album)) {
		scoreB++
	}

	return scoreA > scoreB
}

func main() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Obtener los parámetros de canción, artista y álbum de la URL
		name := r.URL.Query().Get("name")
		artist := r.URL.Query().Get("artist")
		album := r.URL.Query().Get("album")

		// Realizar la búsqueda en ChartLyrics y en iTunes en goroutines
		chartLyricsChan := make(chan []Song)
		iTunesChan := make(chan []Song)

		go func() {
			chartLyricsSongs, err := searchChartLyrics(name, artist, album)
			if err != nil {
				log.Println("Error en la búsqueda de canciones en ChartLyrics:", err)
			}
			chartLyricsChan <- chartLyricsSongs
		}()

		go func() {
			iTunesSongs, err := searchIniTunes(name, artist, album)
			if err != nil {
				log.Println("Error en la búsqueda de canciones en iTunes:", err)
			}
			iTunesChan <- iTunesSongs
		}()

		// Recibir los resultados de las búsquedas
		chartLyricsSongs := <-chartLyricsChan
		iTunesSongs := <-iTunesChan

		// Combinar los resultados de ambas búsquedas
		var selectedSongs []Song

		for _, song := range chartLyricsSongs {
			if strings.Contains(strings.ToLower(song.Name), strings.ToLower(name)) ||
				strings.Contains(strings.ToLower(song.Artist), strings.ToLower(artist)) {
				selectedSongs = append(selectedSongs, song)
			}
		}

		for _, song := range iTunesSongs {
			if strings.Contains(strings.ToLower(song.Name), strings.ToLower(name)) ||
				strings.Contains(strings.ToLower(song.Artist), strings.ToLower(artist)) ||
				strings.Contains(strings.ToLower(song.Album), strings.ToLower(album)) {
				selectedSongs = append(selectedSongs, song)
			}
		}

		// Ordenar las canciones según la coincidencia con los criterios de búsqueda
		sort.Slice(selectedSongs, func(i, j int) bool {
			return compareSongs(&selectedSongs[i], &selectedSongs[j], name, artist, album)
		})

		// Almacenar la respuesta en la base de datos
		err := storeResponseInDB(selectedSongs)
		if err != nil {
			log.Println("Error al almacenar la respuesta en la base de datos:", err)
			http.Error(w, "Error al almacenar la respuesta", http.StatusInternalServerError)
			return
		}

		// Obtener los registros almacenados desde la base de datos
		storedSongs, err := getStoredSongsFromDB()
		if err != nil {
			log.Println("Error al obtener los registros almacenados:", err)
			http.Error(w, "Error al obtener los registros almacenados", http.StatusInternalServerError)
			return
		}

		// Crear la respuesta de búsqueda
		response := SearchResponse{
			Results: storedSongs,
		}

		// Agregar mensaje si no se encontraron canciones
		if len(storedSongs) == 0 {
			response.Message = "No se encontraron canciones que coincidan con los criterios de búsqueda."
		}

		// Convertir la respuesta a formato JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Println("Error al convertir la respuesta a JSON:", err)
			http.Error(w, "Error al convertir la respuesta a JSON", http.StatusInternalServerError)
			return
		}

		// Establecer la cabecera de tipo de contenido JSON
		w.Header().Set("Content-Type", "application/json")

		// Escribir la respuesta JSON
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println("Error al escribir la respuesta JSON:", err)
			http.Error(w, "Error al escribir la respuesta JSON", http.StatusInternalServerError)
			return
		}
	})

	log.Println("Iniciando el servidor en http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
