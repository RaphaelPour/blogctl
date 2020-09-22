package main

import (
	"os"
	"strings"
	"testing"
)

func Test_Main(t *testing.T) {

	/* Overwrite go test arguments with the one we actually
	 * want to test. This is necessary while cobra doesn't
	 * work with go test args.
	 */
	os.Args = append([]string{os.Args[0]}, strings.Split(os.Getenv("TEST_ARGS"), " ")...)

	main()
}
