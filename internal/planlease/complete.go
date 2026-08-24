package planlease

import "context"

func (l *Lease) Complete(ctx context.Context, campaignID, stationID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.campaignID != campaignID || l.stationID != stationID {
		return ErrActive
	}
	return nil
}
