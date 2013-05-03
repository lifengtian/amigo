package utilspub

import (
	"testing"
)

func TestIsFile(t *testing.T) {
	ok := IsFile("/Users/tianl/workspace/amigo/src/CAG/utils/utils.go")
	if !ok {
		t.Errorf("error isFile ")
	}

	cmd := []Cmd{{"pwd", ""}}
	ok = RunIt(false, cmd)
	if !ok {
		t.Errorf("cmd error")
	}

	cmd = []Cmd{{"whoami", ""}}
	ok = RunIt(false, cmd)
	if !ok {
		t.Errorf("cmd error")
	}
}
