package controllers

import (
	"encoding/json"
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type RickAndMortyAPIResponse struct {
	Info struct {
		Count int    `json:"count"`
		Pages int    `json:"pages"`
		Next  string `json:"next"`
		Prev  string `json:"prev"`
	} `json:"info"`
	Results []models.Character `json:"results"`
}

type RickAndMortyCharacter struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Gender  string `json:"gender"`
	Species string `json:"species"`
	Image   string `json:"image"`
}

// Obtener todos los personajes desde la API externa y guardarlos en la base de datos
func FetchAndSaveAllCharacters(c *gin.Context) {
	apiURL := "https://rickandmortyapi.com/api/character"
	var wg sync.WaitGroup // WaitGroup para esperar a que terminen las goroutines

	for {
		// Llamar a la API
		resp, err := http.Get(apiURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los datos de la API externa"})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer los datos de la API"})
			return
		}

		// Parsear la respuesta JSON
		var apiResponse RickAndMortyAPIResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear los datos de la API"})
			return
		}

		// Canal para manejar errores
		errCh := make(chan error, len(apiResponse.Results))

		// Guardar los personajes en la base de datos en paralelo
		for _, apiCharacter := range apiResponse.Results {
			wg.Add(1) // Incrementar el contador del WaitGroup

			go func(apiCharacter models.Character) {
				defer wg.Done() // Decrementar el contador cuando termine
				character := models.Character{
					ID:      apiCharacter.ID,
					Name:    apiCharacter.Name,
					Image:   apiCharacter.Image,
					Status:  apiCharacter.Status,
					Gender:  apiCharacter.Gender,
					Species: apiCharacter.Species,
				}

				if err := character.SaveCharacter(config.DB); err != nil {
					errCh <- err // Enviar el error al canal
				}
			}(apiCharacter)
		}

		// Esperar a que terminen todas las goroutines
		wg.Wait()
		close(errCh)

		// Revisar errores
		for err := range errCh {
			if err != nil {
				log.Printf("Error al guardar el personaje en la base de datos: %v", err)
			}
		}

		// Si no hay más páginas, salir del bucle
		if apiResponse.Info.Next == "" {
			break
		}

		// Configurar la URL para la siguiente página
		apiURL = apiResponse.Info.Next
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todos los personajes han sido guardados"})
}

// Obtener personajes con paginación y búsqueda
func GetPaginatedCharacters(c *gin.Context) {
	// Parámetros de paginación y búsqueda
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	search := c.Query("search")

	// Conversión de strings a enteros con valores por defecto
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Valor por defecto
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Valor por defecto
	}

	// Obtener personajes de la base de datos
	characters, err := models.GetCharacters(config.DB, limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los personajes"})
		return
	}

	c.JSON(http.StatusOK, characters)
}
