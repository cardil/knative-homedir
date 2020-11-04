// +build mage

package main

import (
	"log"
	"os"
	"path"

	"github.com/cardil/homedir/internal"
	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/tasks"
)

// Default target is set to binary
var Default = magetasks.Binary

// Publish will publish built images to a remote registry
func Publish() {
	mg.Deps(Images)
	if len(config.Binaries) > 0 {
		t := tasks.StartMultiline("📤", "Publishing OCI images")
		errs := make([]error, 0)
		for _, binary := range config.Binaries {
			args := []string{
				"push", imagename(binary),
			}
			args = append(args)
			err := sh.RunV("podman", args...)
			errs = append(errs, err)
		}
		t.End(errs...)
	}
}

// Images builds a OCI images of binaries
func Images() {
	mg.Deps(magetasks.Binary)

	if len(config.Binaries) > 0 {
		t := tasks.StartMultiline("📦", "Packaging OCI images")
		errs := make([]error, 0)
		for _, binary := range config.Binaries {
			args := []string{
				"build",
				"-f", containerfile(binary),
				"-t", imagename(binary),
				".",
			}
			args = append(args)
			err := sh.RunV("podman", args...)
			errs = append(errs, err)
		}
		t.End(errs...)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.Binaries = append(config.Binaries, config.Binary{
		Name: internal.BinaryName,
	})
	config.VersionVariablePath = "github.com/cardil/homedir/internal.Version"
}

func containerfile(bin config.Binary) string {
	return path.Join("cmd", bin.Name, "Containerfile")
}

func imagename(bin config.Binary) string {
	return os.Getenv("CONTAINER_BASENAME") + bin.Name
}
