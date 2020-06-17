package domain

import "context"

// Roles ...
type Roles struct {
	ID        int64  `json:"id"`
	Data      string `json:"data"`
}

type RolesData struct {
	RoleName string `json:"role_name"`
	Description string `json:"description"`
}

// RolesRepository represent the author's repository contract
type RolesRepository interface {
	GetByID(ctx context.Context, id int64) (Roles, error)
}
