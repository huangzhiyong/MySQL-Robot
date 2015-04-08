package common

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadLinesAll(t *testing.T) {
	res, err := ReadLinesAll("common_test.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res[0])
	if !strings.Contains(res[0], "package common") {
		t.Error("could not read correctly")
	}
}

func TestReadLinesOffset(t *testing.T) {
	res, err := ReadLinesOffset("common_test.go", 2, 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res[0])
	if !strings.Contains(res[0], "import (") {
		t.Error("could not read correctly")
	}
}
