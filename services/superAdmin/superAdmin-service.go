package superAdmin

import (
	"os"
	"strconv"
)

func IsSuperAdmin(userID int64) bool {
	superAdminId, err := strconv.ParseInt(os.Getenv("SUPER_ADMIN_ID"), 10, 64)
	if err != nil {
		return false
	}

	return superAdminId == userID
}
