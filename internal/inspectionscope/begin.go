package inspectionscope

import "context"

func (s *Scope) Begin(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.owner != "" && s.owner != campaignID {
		return ErrBusy
	}
	s.owner = campaignID
	s.depth++
	return nil
}
