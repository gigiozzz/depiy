package bundles

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestExtractImageTo(t *testing.T) {
	repository := "docker.io/gigiozzz/bundle-test-op"
	concat := "@sha256:"
	digest := "70ba938d4e11f219fc9dc0424e3e55173419a1da51598b341bb2162ea088a8a4"
	dir, err := ioutil.TempDir("/tmp", "crane-"+digest)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("dir: " + dir)

	err = ExtractImageTo(repository+concat+digest, dir)
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := os.Stat(dir + "/descriptor.yaml"); err != nil {
		t.Fatal(err.Error())
	}

	defer os.RemoveAll(dir)

}
