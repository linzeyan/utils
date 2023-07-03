package utils

import (
	"googlemaps.github.io/maps"
)

/* GoogleMapsAPINewClient returns pointer of maps.Client with disables rate limiting. */
func GoogleMapsAPINewClient(apiKey string) *maps.Client {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(0))
	if err != nil {
		logPanic(err)
	}
	return client
}
