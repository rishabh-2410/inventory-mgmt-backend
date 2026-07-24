package requestcontext

import (
	"context"
	enums "inventory-backend/src/models/enums"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const BusinessIDKey contextKey = "business_id"
const RoleIDKey contextKey = "role"

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func GetBusinessID(ctx context.Context) (uuid.UUID, bool) {
	businessID, ok := ctx.Value(BusinessIDKey).(uuid.UUID)
	return businessID, ok
}

func GetUserRole(ctx context.Context) (enums.UserRole, bool) {
	role, ok := ctx.Value(RoleIDKey).(enums.UserRole)
	return role, ok
}
