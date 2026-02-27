package cli

import (
	"fmt"
	"io"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/application/admin"
	"outages-bot/internal/domain"
	"text/tabwriter"
)

// RunUsersCommand lists all users with their Telegram info and addresses.
func RunUsersCommand(
	userRepo domain.UserRepository,
	userInfoProvider application.TelegramUserInfoProvider,
	w io.Writer,
	logger *log.Logger,
) error {
	users := admin.ListUsers(userRepo)

	if len(users) == 0 {
		fmt.Fprintln(w, "No users found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Chat ID\tUsername\tFirst Name\tLast Name\tStreet\tBuilding\tOutage\tComment")

	successCount := 0
	for _, user := range users {
		info, err := userInfoProvider.GetUserInfo(user.ID)
		if err != nil {
			logger.Printf("Failed to get info for chat %d: %v", user.ID, err)
			continue
		}

		username := "-"
		if info.Username != "" {
			username = "@" + info.Username
		}

		outageStr := "-"
		commentStr := "-"
		if user.OutageInfo != nil {
			outageStr = admin.PeriodFormatter(
				user.OutageInfo.Period.StartDate,
				user.OutageInfo.Period.EndDate,
			)
			commentStr = user.OutageInfo.Description.Value
			if commentStr == "" {
				commentStr = "-"
			}
		}

		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			info.ChatID,
			username,
			Sanitize(info.FirstName),
			Sanitize(info.LastName),
			user.Address.StreetName,
			user.Address.Building,
			outageStr,
			commentStr,
		)
		successCount++
	}

	tw.Flush()
	fmt.Fprintf(w, "\nTotal Users: %d\n", successCount)
	return nil
}
