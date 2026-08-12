package app

import "github.com/gofiber/fiber/v3"

type Authorization struct {
	Id        string `json:"id"`
	AccountID string `json:"account_id"`
}

const (
	HeaderYgkAccountID = "X-Ygk-Meta-Account"
)

func ExtractorUserInfo(ctx fiber.Ctx) (auth Authorization, err error) {
	accountID := ctx.Get(HeaderYgkAccountID)
	if accountID == "" {
		// Fallback for local development or simple tests
		accountID = ctx.Get("X-Owned-By")
		if accountID == "" {
			accountID = "local"
		}
	}

	// This is a simplified version of the reference.
	// In production, this would extract from token locals.
	return Authorization{
		Id:        "system", // Placeholder
		AccountID: accountID,
	}, nil
}
