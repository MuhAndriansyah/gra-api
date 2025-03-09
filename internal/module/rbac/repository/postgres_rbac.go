package repository

import (
	"backend-layout/internal/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RBACRepository struct {
	conn *pgxpool.Pool
}

// CheckUserHasPermission implements domain.RBACRepository.
func (r *RBACRepository) CheckUserHasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_role ur
		JOIN role_has_permission rhp ON ur.role_id = rhp.role_id
		JOIN permissions p ON rhp.permission_id = p.id
		WHERE ur.user_id = $1 AND p.name = $2;
	`
	var count int
	err := r.conn.QueryRow(ctx, query, userID, permission).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheskUserHasRole implements domain.RBACRepository.
func (r *RBACRepository) CheskUserHasRole(ctx context.Context, userID int64, role string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM user_role ur
		JOIN roles on ur.role_id = ur.role_id
		where ur.user_id = $1 AND roles.name = $2;
	`

	var count int
	err := r.conn.QueryRow(ctx, query, userID, role).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewRBACRepository(conn *pgxpool.Pool) domain.RBACRepository {
	return &RBACRepository{conn}
}
