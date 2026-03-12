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

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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

	cfg := tablewriter.NewConfigBuilder().
		WithHeaderAutoFormat(tw.Off).
		WithRowAutoWrap(tw.WrapNormal).
		ForColumn(3).WithMaxWidth(30).Build().
		ForColumn(5).WithMaxWidth(30).Build().
		ForColumn(6).WithMaxWidth(20).Build().
		Build()

	table := tablewriter.NewTable(w,
		tablewriter.WithConfig(cfg),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{})),
	)
	table.Header([]string{"Chat ID", "Username", "Name", "Street", "Building", "Outage", "Comment"})

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

		firstName := sanitizeDisplayText(info.FirstName)
		lastName := sanitizeDisplayText(info.LastName)
		var nameParts []string
		if firstName != "-" {
			nameParts = append(nameParts, firstName)
		}
		if lastName != "-" {
			nameParts = append(nameParts, lastName)
		}
		name := "-"
		if len(nameParts) > 0 {
			name = strings.Join(nameParts, " ")
		}

		table.Append([]string{
			fmt.Sprintf("%d", info.ChatID),
			username,
			name,
			user.Address.StreetName,
			user.Address.Building,
			outageStr,
			commentStr,
		})
		successCount++
	}

	table.Render()
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
