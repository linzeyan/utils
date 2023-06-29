package utils

import (
	"context"

	"googlemaps.github.io/maps"
)

type GoogleMaps struct {
	ctx    context.Context
	Client *maps.Client
}

func GoogleMapsNewClient(apiKey string) *GoogleMaps {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	return &GoogleMaps{ctx: ctx, Client: client}
}

func (g *GoogleMaps) Directions(req *maps.DirectionsRequest) ([]maps.Route, []maps.GeocodedWaypoint, error) {
	return g.Client.Directions(g.ctx, req)
}

func (g *GoogleMaps) DistanceMatrix(req *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
	return g.Client.DistanceMatrix(g.ctx, req)
}

func (g *GoogleMaps) Elevation(req *maps.ElevationRequest) ([]maps.ElevationResult, error) {
	return g.Client.Elevation(g.ctx, req)
}

func (g *GoogleMaps) Geocode(req *maps.GeocodingRequest) ([]maps.GeocodingResult, error) {
	return g.Client.Geocode(g.ctx, req)
}

func (g *GoogleMaps) Geolocate(req *maps.GeolocationRequest) (*maps.GeolocationResult, error) {
	return g.Client.Geolocate(g.ctx, req)
}

func (g *GoogleMaps) PlaceDetails(req *maps.PlaceDetailsRequest) (maps.PlaceDetailsResult, error) {
	return g.Client.PlaceDetails(g.ctx, req)
}

func (g *GoogleMaps) PlaceNearbySearch(req *maps.NearbySearchRequest) (maps.PlacesSearchResponse, error) {
	return g.Client.NearbySearch(g.ctx, req)
}

func (g *GoogleMaps) RoadsNearest(req *maps.NearestRoadsRequest) (*maps.NearestRoadsResponse, error) {
	return g.Client.NearestRoads(g.ctx, req)
}

func (g *GoogleMaps) RoadsSnapTo(req *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error) {
	return g.Client.SnapToRoad(g.ctx, req)
}

func (g *GoogleMaps) RoadsSpeedLimits(req *maps.SpeedLimitsRequest) (*maps.SpeedLimitsResponse, error) {
	return g.Client.SpeedLimits(g.ctx, req)
}
