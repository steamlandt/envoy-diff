// Package export serialises diff results into shell-sourceable or
// dotenv-compatible text so that a user can pipe the output of envoy-diff
// directly into their shell or write it to a new .env file.
//
// Supported formats:
//
//	FormatShell  – emits "export KEY=VALUE" lines suitable for eval or source.
//	FormatDotenv – emits "KEY=VALUE" lines compatible with dotenv loaders.
//
// Usage:
//
//	err := export.Write(os.Stdout, entries, export.Options{
//	    Format:      export.FormatShell,
//	    OnlyChanged: true,
//	})
package export
