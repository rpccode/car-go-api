package models

// Invoice es la estructura para una factura
type Invoice struct {
	ID        int
	Amount    float64
	PaymentID int
	// Agrega otros campos necesarios
}

// GenerateInvoice genera una factura a partir de un pago
func GenerateInvoice(payment Payment) Invoice {
	// Lógica para crear la factura
	return Invoice{
		ID:        1,              // Por ejemplo, asignar un ID adecuado
		Amount:    payment.Amount, // Utiliza el monto del pago
		PaymentID: payment.ID,     // Relacionar con el ID del pago
		// Rellena otros campos según sea necesario
	}
}
