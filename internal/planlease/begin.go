package planlease

import "context"

func (l *Lease) Begin(ctx context.Context, campaignID, stationID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.campaignID != "" && (l.campaignID != campaignID || l.stationID != stationID) {
		return ErrActive
	}
	l.campaignID, l.stationID = campaignID, stationID
	return nil
}
