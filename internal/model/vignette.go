package model

import "time"

type VignetteResult struct {
	Plate   string
	Charges []VignetteCharge
}

type VignetteCharge struct {
	ValidSince time.Time
	ValidUntil time.Time
}

func (charge *VignetteCharge) IsValidFor(now time.Time) bool {
	return charge.ValidSince.Before(now) && charge.ValidUntil.After(now)
}

type VignetteAuth struct {
	Token      string
	Expiration time.Time
}

func (auth *VignetteAuth) ExpiresSoon(soon time.Duration) bool {
	now := time.Now()
	if auth.Expiration.Before(now) {
		return true
	}

	diff := now.Sub(auth.Expiration)
	return diff <= soon
}
