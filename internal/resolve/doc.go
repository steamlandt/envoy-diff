// Package resolve implements variable interpolation for environment maps.
//
// It expands references of the form $VAR and ${VAR} found in environment
// values, using other keys within the same map as the source of truth.
// Optionally, unresolved references can fall back to the host process
// environment via os.Getenv.
//
// Typical usage:
//
//	env, _ := parser.ParseFile(".env")
//	resolved := resolve.Apply(env, resolve.DefaultOptions())
//
// Circular references are broken by a configurable MaxDepth limit (default 10)
// rather than by cycle detection, keeping the implementation simple.
package resolve
