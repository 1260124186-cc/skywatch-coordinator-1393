package planlease

func (l *Lease) Active() (string, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.campaignID, l.stationID
}
