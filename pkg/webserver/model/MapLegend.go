package model

type ResourceIdentity struct {
	Name string
	Id   string
}

type MapLegend struct {
	Harbors    []ResourceIdentity
	Landscapes []ResourceIdentity
}
