//nolint:forbidigo
package testhelper

import (
	"fmt"
	"strings"

	"github.com/keboola/go-utils/pkg/wildcards"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

// fileNode is one file/dir in expected or actual directory.
type fileNode struct {
	isDir   bool
	absPath string
}

// fileNodeState in expected and actual directory.
type fileNodeState struct {
	relPath  string
	expected *fileNode
	actual   *fileNode
}

// DirectoryContentsSame compares two directories, in expected file content can be used wildcards.
func DirectoryContentsSame(expectedFs filesystem.Fs, expectedDir string, actualFs filesystem.Fs, actualDir string) error {
	nodesState := compareDirectories(expectedFs, expectedDir, actualFs, actualDir)
	var errs []string
	for _, node := range nodesState {
		// Check if present if both dirs (actual/expected) and if has same type (file/dir)
		switch {
		case node.actual == nil:
			errs = append(errs, fmt.Sprintf("only in expected \"%s\"", node.expected.absPath))
		case node.expected == nil:
			errs = append(errs, fmt.Sprintf("only in actual \"%s\"", node.actual.absPath))
		case node.actual.isDir != node.expected.isDir:
			if node.actual.isDir {
				errs = append(errs, fmt.Sprintf("\"%s\" is dir in actual, but file in expected", node.relPath))
			} else {
				errs = append(errs, fmt.Sprintf("\"%s\" is file in actual, but dir in expected", node.relPath))
			}
		default:
			// Compare content
			if !node.actual.isDir {
				expectedFile, err := expectedFs.ReadFile(filesystem.NewFileDef(node.expected.absPath))
				if err != nil {
					return err
				}
				actualFile, err := actualFs.ReadFile(filesystem.NewFileDef(node.actual.absPath))
				if err != nil {
					return err
				}
				err = wildcards.Compare(
					expectedFile.Content,
					actualFile.Content,
				)
				if err != nil {
					return errors.PrefixErrorf(err, `different content of the file "%s"`, node.relPath)
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.New("Directories are not same:\n" + strings.Join(errs, "\n"))
	}
	return nil
}

// AssertDirectoryContentsSame compares two directories, in expected file content can be used wildcards.
func AssertDirectoryContentsSame(t assert.TestingT, expectedFs filesystem.Fs, expectedDir string, actualFs filesystem.Fs, actualDir string) {
	err := DirectoryContentsSame(expectedFs, expectedDir, actualFs, actualDir)
	if err != nil {
		assert.Fail(t, err.Error())
	}
}

func compareDirectories(expectedFs filesystem.Fs, expectedDir string, actualFs filesystem.Fs, actualDir string) map[string]*fileNodeState {
	// relative path -> state
	hashMap := map[string]*fileNodeState{}
	var err error

	// Process actual dir
	err = actualFs.Walk(actualDir, func(path string, info filesystem.FileInfo, err error) error {
		// Stop on error
		if err != nil {
			return err
		}

		// Ignore root
		if path == actualDir {
			return nil
		}

		// Ignore hidden files, except .env*, .gitignore
		if IsIgnoredFile(path, info) {
			return nil
		}

		// Get relative path
		relPath, err := filesystem.Rel(actualDir, path)
		if err != nil {
			return err
		}

		// Create node
		hashMap[relPath] = &fileNodeState{
			relPath: relPath,
			actual:  &fileNode{info.IsDir(), path},
		}

		return nil
	})

	if err != nil {
		panic(errors.Errorf(`cannot iterate over directory "%s" in "%s": %w`, actualDir, actualFs.BasePath(), err))
	}

	// Process expected dir
	err = expectedFs.Walk(expectedDir, func(path string, info filesystem.FileInfo, err error) error {
		// Stop on error
		if err != nil {
			return err
		}

		// Ignore root
		if path == expectedDir {
			return nil
		}

		// Ignore hidden files, except .env*, .gitignore
		if IsIgnoredFile(path, info) {
			return nil
		}

		// Get relative path
		relPath, err := filesystem.Rel(expectedDir, path)
		if err != nil {
			return err
		}

		// Create node if not exists
		if _, ok := hashMap[relPath]; !ok {
			hashMap[relPath] = &fileNodeState{}
		}
		hashMap[relPath].relPath = relPath
		hashMap[relPath].expected = &fileNode{info.IsDir(), path}

		return nil
	})

	if err != nil {
		panic(errors.Errorf(`cannot iterate over directory "%s" in "%s": %w`, expectedDir, expectedFs.BasePath(), err))
	}

	return hashMap
}
