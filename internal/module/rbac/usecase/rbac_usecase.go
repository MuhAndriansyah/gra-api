package usecase

import (
	"backend-layout/internal/domain"
	"context"
)

type RBACUsecase struct {
	repo domain.RBACRepository
}

// CheckUserHasPermission implements domain.RBACUsecase.
func (r *RBACUsecase) CheckUserHasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	hasPermission, err := r.repo.CheckUserHasPermission(ctx, userID, permission)

	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

// CheskUserHasRole implements domain.RBACUsecase.
func (r *RBACUsecase) CheskUserHasRole(ctx context.Context, userID int64, role string) (bool, error) {
	hasRole, err := r.repo.CheskUserHasRole(ctx, userID, role)

	if err != nil {
		return false, err
	}

	return hasRole, nil
}

func NewRBACUsecase(repo domain.RBACRepository) domain.RBACUsecase {
	return &RBACUsecase{repo}
}
