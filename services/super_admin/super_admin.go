package super_admin

import (
	"Crawlzilla/database/repositories"
	"Crawlzilla/models/ads"
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

func ValidateCrawlResultData(result *ads.CrawlResult) error {
	// some codes.
	return nil
}

func AddAdForSuperAdmin(result *ads.CrawlResult) error {
	if err := ValidateCrawlResultData(result); err != nil {
		return err
	}
	return repositories.AddCrawlResult(result)
}
