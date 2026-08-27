package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

type Repository interface {
	CreateCampaign(domain.Campaign) error
	GetCampaign(string) (domain.Campaign, error)
	UpdateCampaign(domain.Campaign) error
	ListCampaigns() []domain.Campaign
	AddStation(domain.Station) error
	GetStation(string) (domain.Station, error)
	ListStations(string) []domain.Station
	CreateShift(domain.Shift) error
	GetShift(string) (domain.Shift, error)
	UpdateShift(domain.Shift) error
	ListShifts(string) []domain.Shift
	CreateObservation(domain.Observation) error
	GetObservation(string) (domain.Observation, error)
	UpdateObservation(domain.Observation) error
	ListObservations(string) []domain.Observation
	SaveRelease(domain.Release) error
	GetRelease(string) (domain.Release, error)
}
