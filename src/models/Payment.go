package models

import (
	"database/sql"
	"time"
)

type Payment struct {
	ID            int       `json:"id"`
	ReservationID int       `json:"reservation_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"` // pagado, pendiente
	PaymentDate   time.Time `json:"payment_date"`
}

// Procesar el pago
func (p *Payment) ProcessPayment(db *sql.DB) error {
	query := `INSERT INTO payments (reservation_id, amount, status, payment_date) 
              VALUES ($1, $2, 'pagado', $3) RETURNING id`
	return db.QueryRow(query, p.ReservationID, p.Amount, time.Now()).Scan(&p.ID)
}

// // Generar factura
// func GenerateInvoice(payment Payment) string {
// 	return fmt.Sprintf("Factura ID: %d, Monto: %.2f, Fecha: %s", payment.ID, payment.Amount, payment.PaymentDate.String())
// }
