package notifier

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

type Snapshot struct {
	Count       int
	Active      string
	LastError   error
	CampaignIDs []string
}

func (r *Registry) Snapshot() Snapshot {
	ids := make([]string, 0, len(r.signals))
	for campaignID := range r.signals {
		ids = append(ids, campaignID)
	}
	return Snapshot{Count: len(r.signals), Active: r.activeCampaign, LastError: r.lastError, CampaignIDs: ids}
}

func (r *Registry) HasSignal(campaignID string) bool {
	_, ok := r.signals[campaignID]
	return ok
}

func (r *Registry) SignalFor(campaignID string) (ReleaseSignal, bool) {
	signal, ok := r.signals[campaignID]
	return signal, ok
}

func (r *Registry) Match(release domain.Release) bool {
	signal, ok := r.SignalFor(release.CampaignID)
	return ok && signal.ReleaseID == release.ID && signal.Reviewer == release.ReleasedBy
}
