package util

import (
	"io/ioutil"
	"os"
	"testing"
)

const perm = 0755

func TestCopyDir(t *testing.T) {
	srcDir, err := ioutil.TempDir("", "srcDir")
	if err != nil {
		t.Fatalf("error creating temp directory (%s)", err)
	}
	destDir, err := ioutil.TempDir("", "destDir")
	if err != nil {
		t.Fatalf("error creating temp directory (%s)", err)
	}
	data := []byte("web: example-go")
	if err = ioutil.WriteFile(srcDir+"/chart", data, perm); err != nil {
		t.Fatalf("error creating %s/chart (%s)", srcDir, err)
	}
	defer func() {
		if err = os.RemoveAll(srcDir); err != nil {
			t.Fatalf("failed to remove %s directory (%s)", srcDir, err)
		}
		if err = os.RemoveAll(destDir); err != nil {
			t.Fatalf("failed to remove %s directory (%s)", destDir, err)
		}
	}()
	if err = CopyDir(srcDir, destDir); err != nil {
		t.Errorf("[Expected] error not have occured\n[Got] error (%s)\n", err)
	}
	fileInfo, err := os.Stat(destDir + "/chart")
	if err != nil {
		t.Fatalf("error getting file info in destination directory (%s)", err)
	}

	if perm != fileInfo.Mode() {
		t.Errorf("\n[Expected] %d\n[Got] %d\n", perm, fileInfo.Mode())
	}
}
