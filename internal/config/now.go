package config

import "time"

var (
	// Now is a function that returns the current time.
	// This is useful for testing purposes, as it allows us to
	// mock the current time in tests.
	Now = time.Now
)
