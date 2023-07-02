package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"googlemaps.github.io/maps"
)

const api = "APIKEY"

func TestGoogleMapsDirections(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.DirectionsRequest{
		Origin:      "Google Sydney",
		Destination: "Glebe Pt Rd, Glebe",
		Mode:        maps.TravelModeTransit,
	}
	ctx := context.Background()
	result, _, err := GoogleMapsAPINewClient(api).Directions(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result[0])
	assertion.Nil(err)
	// assertion.JSONEq(response, s)
}

func TestGoogleMapsDistanceMatrix(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.DistanceMatrixRequest{
		Origins:      []string{},
		Destinations: []string{},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).DistanceMatrix(ctx, r)
	requirement.Error(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsElevation(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	response := `{"location":{"lat":24.72969,"lng":120.87882},"elevation":21.09626007080078,"resolution":9.543951988220215}`
	r := &maps.ElevationRequest{
		Locations: []maps.LatLng{
			{Lng: 120.8788299, Lat: 24.72969236},
		},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).Elevation(ctx, r)
	requirement.Nil(err)
	s, err := JSONMarshalString(result[0])
	assertion.Nil(err)
	assertion.Equal(response, s)
}

func TestGoogleMapsGeocoding(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lng: 120.8788299,
			Lat: 24.72969236,
		},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).Geocode(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result[0])
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsGeolocation(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.GeolocationRequest{}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).Geolocate(ctx, r)
	requirement.Error(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsNearbySearch(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lng: 120.8788299,
			Lat: 24.72969236,
		},
		Radius: 100,
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).NearbySearch(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsNearestRoads(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.NearestRoadsRequest{
		Points: []maps.LatLng{{Lng: 120.8788299, Lat: 24.72969236}},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).NearestRoads(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsPlaceDetails(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.PlaceDetailsRequest{
		PlaceID: "ChIJeSlIfru0aTQR0LK8Cz0vfDA",
		Fields:  []maps.PlaceDetailsFieldMask{},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).PlaceDetails(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsSnapToRoads(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.SnapToRoadRequest{
		Path: []maps.LatLng{{Lng: 120.8788299, Lat: 24.72969236}},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).SnapToRoad(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}

func TestGoogleMapsSpeedLimits(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	r := &maps.SpeedLimitsRequest{
		Path: []maps.LatLng{{Lng: 120.8788299, Lat: 24.72969236}},
	}
	ctx := context.Background()
	result, err := GoogleMapsAPINewClient(api).SpeedLimits(ctx, r)
	requirement.Nil(err)
	_, err = JSONMarshalString(result)
	assertion.Nil(err)
	// assertion.Equal("", s)
}
