// Build and run the hello-world Echo web service.

package main

import (
	"context"
	"dagger/dagger-github-cache-test/internal/dagger"
)

type DaggerGithubCacheTest struct{}

// Build a static Linux binary of the web service from source.
func (m *DaggerGithubCacheTest) Build(source *dagger.Directory) *dagger.File {
	return dag.Container().
		From("golang:1.25").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", "linux").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build")).
		WithExec([]string{"go", "build", "-o", "/out/server", "."}).
		File("/out/server")
}

// Package the web service into a minimal runtime container exposing port 8080.
func (m *DaggerGithubCacheTest) Container(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("alpine:3.20").
		WithFile("/usr/local/bin/server", m.Build(source)).
		WithExposedPort(8080).
		WithEntrypoint([]string{"/usr/local/bin/server"})
}

// Publish the web service container image to the given registry address.
func (m *DaggerGithubCacheTest) Publish(ctx context.Context, source *dagger.Directory, address string) (string, error) {
	return m.Container(source).Publish(ctx, address)
}

// Start the web service as an ephemeral Dagger service on port 8080.
func (m *DaggerGithubCacheTest) Serve(source *dagger.Directory) *dagger.Service {
	return m.Container(source).AsService()
}
