package main

var (
	// Version package version, specified when built with ldflags
	Version = `Not specified, use --ldflags "-X main.Version "1.0.0""`

	// Sha git commit sha, specified when built with ldflags
	Sha = "Not specified"
)
