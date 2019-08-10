package model

// ResourceIdentity Mapping resource names for both Harbors and Landscapes to their respective ID
type ResourceIdentity struct {
	Name string
	Id   string
}

// MapLegend Legend for API uses, which allows use of codes (which can than be mapped via the Legend
type MapLegend struct {
	Harbors    []ResourceIdentity
	Landscapes []ResourceIdentity
}
