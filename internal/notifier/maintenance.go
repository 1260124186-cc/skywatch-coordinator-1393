package notifier

import "context"

func (r *Registry) Reset(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.activeCampaign = ""
	r.lastError = nil
	return nil
}

func (r *Registry) Clear(ctx context.Context) error {
	if err := r.Reset(ctx); err != nil {
		return err
	}
	r.signals = make(map[string]ReleaseSignal)
	return nil
}

func (r *Registry) Healthy() bool {
	return r.lastError == nil && r.activeCampaign == ""
}
