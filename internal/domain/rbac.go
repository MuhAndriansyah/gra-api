package domain

import "context"

type RBACRepository interface {
	CheckUserHasPermission(ctx context.Context, userID int64, permission string) (bool, error)
	CheskUserHasRole(ctx context.Context, userID int64, role string) (bool, error)
}

type RBACUsecase interface {
	CheckUserHasPermission(ctx context.Context, userID int64, permission string) (bool, error)
	CheskUserHasRole(ctx context.Context, userID int64, role string) (bool, error)
}
