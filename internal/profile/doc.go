// Package profile manages named environment profiles.
//
// A profile is a saved collection of key/value pairs stored as a JSON file
// on disk. Profiles can be created from any map[string]string source,
// retrieved by name, enumerated, and deleted.
//
// Typical usage:
//
//	dir := filepath.Join(os.UserConfigDir(), "envoy-diff", "profiles")
//
//	// save current env as "production"
//	profile.Save(dir, "production", env)
//
//	// load it back later
//	p, err := profile.Load(dir, "production")
//
//	// list all saved profiles
//	names, err := profile.List(dir)
package profile
