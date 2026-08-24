package inspectionscope

import "context"

func (s *Scope) Begin(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.campaignID != "" && s.campaignID != campaignID {
		return ErrBusy
	}
	s.campaignID = campaignID
	return nil
}
