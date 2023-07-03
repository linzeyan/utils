package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"googlemaps.github.io/maps"
)

func TestGoogleMapsAPINewClient(t *testing.T) {
	assertion := assert.New(t)
	requirement := require.New(t)
	var api string

	api = ""
	client, err := maps.NewClient(maps.WithAPIKey(api), maps.WithRateLimit(0))
	requirement.Error(err)
	requirement.Nil(client)
	requirement.Panics(func() {
		GoogleMapsAPINewClient(api)
	})

	api = "APIKEY"
	client, err = maps.NewClient(maps.WithAPIKey(api), maps.WithRateLimit(0))
	requirement.Nil(err)
	requirement.NotNil(client)
	client1 := GoogleMapsAPINewClient(api)
	assertion.NotNil(client1)
	assertion.EqualValues(client, client1)
}
