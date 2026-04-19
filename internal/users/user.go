package users

import "outages-bot/internal/outage"

// User represents a subscribed user.
type User struct {
	ID         int64
	Address    Address
	OutageInfo *OutageInfo
}

// WithNotifiedOutage returns a new User with the outage info set from the given outage.
func (u *User) WithNotifiedOutage(current *outage.Outage) *User {
	info := NewOutageInfo(current.Period, current.Description)
	return &User{
		ID:         u.ID,
		Address:    u.Address,
		OutageInfo: &info,
	}
}

// FindOutageForNotification finds the first matching outage for a user that they haven't been notified about.
func (u *User) FindOutageForNotification(allOutages []*outage.Outage) *outage.Outage {
	for _, current := range allOutages {
		if current.Address.StreetID != u.Address.StreetID {
			continue
		}
		matchedBuilding := false
		for _, building := range current.Address.Buildings {
			if building == u.Address.Building {
				matchedBuilding = true
				break
			}
		}
		if !matchedBuilding {
			continue
		}

		outageInfo := NewOutageInfo(current.Period, current.Description)

		if u.OutageInfo != nil && u.OutageInfo.Equals(outageInfo) {
			return nil
		}

		return current
	}

	return nil
}
