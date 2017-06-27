package kons

import "math"

type UpsertSettings struct {
	MakePath bool
}

type FindSettings struct {
	Limit int32
}

func DefaultUpsertSettings() *UpsertSettings {
	return &UpsertSettings{
		MakePath: true,
	}
}

func DefaultFindSettings() *FindSettings {
	return &FindSettings{
		Limit: math.MaxInt32,
	}
}
