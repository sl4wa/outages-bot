package domain

// FindOutageForNotification finds the first matching outage for a user that they haven't been notified about.
func FindOutageForNotification(user *User, allOutages []*Outage) *Outage {
	for _, outage := range allOutages {
		if !outage.AffectsUserAddress(user.Address) {
			continue
		}

		outageInfo := NewOutageInfo(outage.Period, outage.Description)

		if user.IsAlreadyNotifiedAbout(outageInfo) {
			return nil
		}

		return outage
	}

	return nil
}
