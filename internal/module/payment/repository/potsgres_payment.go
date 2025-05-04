package repository

import (
	"backend-layout/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPaymentRepository struct {
	conn *pgxpool.Pool
}

// UpdatePaymentStatus implements domain.PaymentRepository.
func (p *PostgresPaymentRepository) SetPaymentFailed(ctx context.Context, tx pgx.Tx, orderId int64) error {
	query := `UPDATE orders 
	          SET payment_status = 'failure' 
			  WHERE id = $1 AND payment_status = 'Pending';`

	cmdTag, err := tx.Exec(ctx, query, orderId)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("payment not processed: order not found or already paid")
	}

	return nil
}

const paymentStatusPaid = "Paid"

// GetTx implements domain.OrderRepository.
func (p *PostgresPaymentRepository) GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.conn.Begin(ctx)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

// ProcessPayment implements domain.PaymentRepository.
func (p *PostgresPaymentRepository) ProcessPayment(ctx context.Context, tx pgx.Tx, userId int64, paymentMethod, orderNumber string) error {
	query := `UPDATE orders 
	          SET payment_status = $1, payment_date = timezone('UTC', now()), payment_method = $2 
			  WHERE order_number = $3 AND user_id = $4 AND payment_status = 'Pending';`

	cmdTag, err := tx.Exec(ctx, query, paymentStatusPaid, paymentMethod, orderNumber, userId)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("order already paid")
	}

	return nil
}

func NewPostgresPaymentRepository(conn *pgxpool.Pool) domain.PaymentRepository {
	return &PostgresPaymentRepository{conn: conn}
}
