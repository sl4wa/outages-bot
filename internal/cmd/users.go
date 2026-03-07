package cli

import (
	"fmt"
	"io"
	"log"
	"outages-bot/internal/application"
	"outages-bot/internal/application/admin"
	"outages-bot/internal/domain"
	"regexp"
	"strings"
	"text/tabwriter"
)

// RunUsersCommand lists all users with their Telegram info and addresses.
func RunUsersCommand(
	userRepo domain.UserRepository,
	userInfoProvider application.UserInfoProvider,
	w io.Writer,
	logger *log.Logger,
) {
	users := admin.ListUsers(userRepo)

	if len(users) == 0 {
		fmt.Fprintln(w, "No users found.")
		return
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
			sanitizeDisplayText(info.FirstName),
			sanitizeDisplayText(info.LastName),
			user.Address.StreetName,
			user.Address.Building,
			outageStr,
			commentStr,
		)
		successCount++
	}

	tw.Flush()
	fmt.Fprintf(w, "\nTotal Users: %d\n", successCount)
}

// sanitizeRegex matches Unicode format characters (\p{Cf}) and Hangul filler (U+3164).
var sanitizeRegex = regexp.MustCompile(`[\p{Cf}\x{3164}]`)

// sanitizeDisplayText removes invisible/control Unicode characters from a string.
// Returns "-" for empty or whitespace-only results.
func sanitizeDisplayText(value string) string {
	if value == "" {
		return "-"
	}

	cleaned := sanitizeRegex.ReplaceAllString(value, "")
	trimmed := strings.TrimSpace(cleaned)

	if trimmed == "" {
		return "-"
	}

	return trimmed
}
