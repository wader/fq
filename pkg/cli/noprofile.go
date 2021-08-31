//go:build !profile

package cli

func maybeProfile() func() { return func() {} }
