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
	// Release the lease so independent campaigns/stations can plan afterward.
	// Without clearing, the first Begin would leave the slot occupied forever
	// and every subsequent Begin for a different campaign/station would fail.
	l.campaignID, l.stationID = "", ""
	return nil
}
