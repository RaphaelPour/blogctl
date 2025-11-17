package common

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kballard/go-shellquote"
)

func Unwrap[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

var SlugRegex = regexp.MustCompile(`[^A-Za-z0-9-]`)

func Slug(s string) string {
	return SlugRegex.ReplaceAllString(
		strings.ReplaceAll(strings.ToLower(s), " ", "-"),
		"",
	)
}

func CopyFile(sourceFile, destinationFile string) error {
	src, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("error open source file %s: %w", sourceFile, err)
	}
	defer src.Close()

	dst, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("error open destination file %s: %w", destinationFile, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}
	return nil
}

/* https://github.com/MarkusFreitag/changelogger/blob/master/pkg/editor/editor.go */
func Editor(initialValue *string) error {
	var editor string
	if val, ok := os.LookupEnv("EDITOR"); ok {
		editor = val
	}
	if val, ok := os.LookupEnv("VISUAL"); ok {
		editor = val
	}
	if editor == "" {
		return fmt.Errorf("set EDITOR or VISUAL variable")
	}

	file, err := os.CreateTemp("", "new-post*.md")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(*initialValue); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	args, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	args = append(args, file.Name())

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	raw, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	*initialValue = string(raw)
	return nil
}
