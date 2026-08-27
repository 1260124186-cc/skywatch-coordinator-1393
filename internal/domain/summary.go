package domain

type CampaignSummary struct {
	Campaign         Campaign `json:"campaign"`
	StationCount     int      `json:"stationCount"`
	ShiftCount       int      `json:"shiftCount"`
	OpenShiftCount   int      `json:"openShiftCount"`
	ObservationCount int      `json:"observationCount"`
	CandidateCount   int      `json:"candidateCount"`
	AcceptedCount    int      `json:"acceptedCount"`
	RejectedCount    int      `json:"rejectedCount"`
	Release          *Release `json:"release,omitempty"`
}

func (s CampaignSummary) IsReadyForRelease() bool {
	return s.OpenShiftCount == 0 && s.CandidateCount == 0 && s.AcceptedCount > 0
}
