package domain

import "math"

type QualityBand string

const (
	QualityUnknown QualityBand = "unknown"
	QualityPoor    QualityBand = "poor"
	QualityUsable  QualityBand = "usable"
	QualityStrong  QualityBand = "strong"
)

type QualityProfile struct {
	Minimum float64     `json:"minimum"`
	Maximum float64     `json:"maximum"`
	Band    QualityBand `json:"band"`
}

func BandForQuality(value float64) QualityBand {
	switch {
	case value < 0 || value > 1 || math.IsNaN(value):
		return QualityUnknown
	case value < 0.35:
		return QualityPoor
	case value < 0.75:
		return QualityUsable
	default:
		return QualityStrong
	}
}

func NewQualityProfile(values []float64) QualityProfile {
	if len(values) == 0 {
		return QualityProfile{Band: QualityUnknown}
	}
	min, max, sum := values[0], values[0], 0.0
	for _, value := range values {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
		sum += value
	}
	return QualityProfile{Minimum: min, Maximum: max, Band: BandForQuality(sum / float64(len(values)))}
}

func (p QualityProfile) Accepts(value float64) bool {
	if p.Band == QualityUnknown {
		return false
	}
	return value >= p.Minimum && value <= p.Maximum
}

func (p QualityProfile) Spread() float64 { return p.Maximum - p.Minimum }

func (p QualityProfile) IsConsistent() bool {
	return p.Band != QualityUnknown && p.Spread() <= 0.55
}

type QualityDistribution struct {
	Unknown int `json:"unknown"`
	Poor    int `json:"poor"`
	Usable  int `json:"usable"`
	Strong  int `json:"strong"`
}

func (d *QualityDistribution) Add(value float64) {
	switch BandForQuality(value) {
	case QualityPoor:
		d.Poor++
	case QualityUsable:
		d.Usable++
	case QualityStrong:
		d.Strong++
	default:
		d.Unknown++
	}
}

func (d QualityDistribution) Total() int    { return d.Unknown + d.Poor + d.Usable + d.Strong }
func (d QualityDistribution) Reliable() int { return d.Usable + d.Strong }
func (d QualityDistribution) ReliabilityRatio() float64 {
	if d.Total() == 0 {
		return 0
	}
	return float64(d.Reliable()) / float64(d.Total())
}
