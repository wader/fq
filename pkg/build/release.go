//go:build release

package build

func initReleaseOptions() {
	CheckChildAlreadyExists = false
}
