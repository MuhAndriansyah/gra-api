package repository

import (
	"backend-layout/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresUserRepository struct {
	conn *pgxpool.Pool
}

var querySelectUser = `SELECT id, name, email, password, photo, email_verify_code, email_verify_code_expired_at, verified_at FROM users WHERE 1=1`

// GetByEmail implements domain.UserRepository.
func (p *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	query := querySelectUser + ` AND email = $1;`

	err := p.conn.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Photo,
		&user.EmailVerifyCode,
		&user.EmailVerifyCodeExpiredAt,
		&user.VerifiedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil

}

// GetByEmailVerifyCode implements domain.UserRepository.
func (p *postgresUserRepository) GetByEmailVerifyCode(ctx context.Context, verifyCode string, id int64) (*domain.User, error) {
	user := &domain.User{}

	query := querySelectUser + ` AND email_verify_code = $1 AND id = $2;`

	err := p.conn.QueryRow(ctx, query, verifyCode, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Photo,
		&user.EmailVerifyCode,
		&user.EmailVerifyCodeExpiredAt,
		&user.VerifiedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

// GetByID implements domain.UserRepository.
func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}

	query := querySelectUser + ` AND id = $1;`

	err := p.conn.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Photo,
		&user.EmailVerifyCode,
		&user.EmailVerifyCodeExpiredAt,
		&user.VerifiedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

// Store implements domain.UserRepository.
func (p *postgresUserRepository) Store(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (name, email, password, email_verify_code, email_verify_code_expired_at) VALUES($1, $2, $3, $4, $5);`

	_, err := p.conn.Exec(ctx, query, user.Name, user.Email, user.Password, user.EmailVerifyCode, user.EmailVerifyCodeExpiredAt)

	if err != nil {
		return err
	}

	return nil
}

func NewPostgresUserRepository(conn *pgxpool.Pool) domain.UserRepository {
	return &postgresUserRepository{
		conn: conn,
	}
}
