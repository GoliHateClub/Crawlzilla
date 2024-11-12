package bot

import (
	"context"
	"fmt"
)

func StartBot(ctx context.Context) {
	fmt.Println("Hello from Bot")

	// Listen for the shutdown signal from the context
	select {
	case <-ctx.Done(): // triggered if the server's context is canceled
		fmt.Println("Bot received shutdown signal from context, stopping...")
		return
	}
}
