package ci_test

import (
	"testing"

	"github.com/gophergala2016/hugoku/ci"
)

func TestBuild(t *testing.T) {
	_, err := ci.Build("josgilmo", "example-site", "josgilmo/example-site")
	if err != nil {
		t.Error("It shuild be nil")
	}
}
