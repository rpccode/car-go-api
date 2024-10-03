package controllers

import (
	"go-auth-api/src/config"
	"go-auth-api/src/models"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func ProcessReservationPayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Definimos una WaitGroup para esperar a que ambos procesos terminen
	var wg sync.WaitGroup
	wg.Add(2) // Vamos a realizar dos tareas en paralelo

	// Canal para capturar errores
	errChan := make(chan error, 2)
	// Canal para capturar la factura generada
	invoiceChan := make(chan models.Invoice)

	// Procesar el pago en paralelo
	go func() {
		defer wg.Done() // Marca la goroutine como terminada
		if err := payment.ProcessPayment(config.DB); err != nil {
			errChan <- err
		}
	}()

	// Generar la factura en paralelo
	go func() {
		defer wg.Done()                            // Marca la goroutine como terminada
		invoice := models.GenerateInvoice(payment) // Esto debería devolver un models.Invoice
		invoiceChan <- invoice                     // Enviar el Invoice al canal
	}()

	// Esperamos a que ambos procesos terminen
	wg.Wait()
	close(errChan)
	close(invoiceChan)

	// Comprobar si hubo algún error
	for err := range errChan {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el pago o generar la factura"})
			return
		}
	}

	// Recuperar la factura generada
	invoice := <-invoiceChan

	// Respuesta de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Pago procesado", "invoice": invoice})
}
