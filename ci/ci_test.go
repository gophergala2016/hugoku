package ci_test

import (
	"github.com/gophergala2016/hugoku/ci"
	"testing"
)

func TestBuild(t *testing.T) {
	err := ci.Build("josgilmo", "example-site", "josgilmo/example-site")
	if err != nil {
		t.Error("It shuild be nil")
	}
}
