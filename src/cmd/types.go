// encdec
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/types.go
// Original timestamp : 2026.11.08

package cmd

// buildVersion and buildDate are injected at link time via
// `-ldflags "-X encdec/cmd.buildVersion=... -X encdec/cmd.buildDate=..."`
// by every build path (src/build.sh and each __*/{Makefile,APKBUILD,PKGBUILD,spec});
// a plain `go build` leaves them at these fallbacks.
var buildVersion = "dev"
var buildDate = "unknown"
