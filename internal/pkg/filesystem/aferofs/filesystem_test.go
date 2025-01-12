package aferofs_test

import (
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/keboola/go-utils/pkg/orderedmap"
	"github.com/stretchr/testify/assert"

	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	. "github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs/mountfs"
	"github.com/keboola/keboola-as-code/internal/pkg/log"
)

func TestLocalFilesystem(t *testing.T) {
	t.Parallel()
	createFs := func() (filesystem.Fs, log.DebugLogger) {
		logger := log.NewDebugLogger()
		fs, err := NewLocalFs(t.TempDir(), filesystem.WithLogger(logger), filesystem.WithWorkingDir(filesystem.Join("my", "dir")))
		assert.NoError(t, err)
		return fs, logger
	}
	cases := &testCases{createFs}
	cases.runTests(t)
}

func TestMemoryFilesystem(t *testing.T) {
	t.Parallel()
	createFs := func() (filesystem.Fs, log.DebugLogger) {
		logger := log.NewDebugLogger()
		fs := NewMemoryFs(filesystem.WithLogger(logger), filesystem.WithWorkingDir(filesystem.Join("my", "dir")))
		return fs, logger
	}
	cases := &testCases{createFs}
	cases.runTests(t)
}

func TestMountFilesystem_WithoutMountPoint(t *testing.T) {
	t.Parallel()
	createFs := func() (filesystem.Fs, log.DebugLogger) {
		logger := log.NewDebugLogger()
		rootFs := NewMemoryFs(filesystem.WithLogger(logger), filesystem.WithWorkingDir(filesystem.Join("my", "dir")))
		fs, err := NewMountFs(rootFs, nil)
		assert.NoError(t, err)
		return fs, logger
	}
	cases := &testCases{createFs}
	cases.runTests(t)
}

func TestMountFilesystem_WithMountPoint(t *testing.T) {
	t.Parallel()
	createFs := func() (filesystem.Fs, log.DebugLogger) {
		logger := log.NewDebugLogger()
		rootFs := NewMemoryFs(filesystem.WithLogger(logger), filesystem.WithWorkingDir(filesystem.Join("my", "dir")))
		mountPointFs := NewMemoryFs(filesystem.WithLogger(logger))
		fs, err := NewMountFs(rootFs, []mountfs.MountPoint{mountfs.NewMountPoint(filesystem.Join("sub", "dir1"), mountPointFs)})
		assert.NoError(t, err)
		return fs, logger
	}
	cases := &testCases{createFs}
	cases.runTests(t)
}

func TestMountFilesystem_WithNestedMountPoint(t *testing.T) {
	t.Parallel()
	createFs := func() (filesystem.Fs, log.DebugLogger) {
		logger := log.NewDebugLogger()
		rootFs := NewMemoryFs(filesystem.WithLogger(logger), filesystem.WithWorkingDir(filesystem.Join("my", "dir")))
		mountPoint1Fs := NewMemoryFs(filesystem.WithLogger(logger))
		mountPoint2Fs := NewMemoryFs(filesystem.WithLogger(logger))
		fs, err := NewMountFs(
			rootFs,
			[]mountfs.MountPoint{
				mountfs.NewMountPoint("sub/dir1", mountPoint1Fs),
				mountfs.NewMountPoint("sub/dir1/dir2", mountPoint2Fs),
			},
		)
		assert.NoError(t, err)
		return fs, logger
	}
	cases := &testCases{createFs}
	cases.runTests(t)
}

type testCases struct {
	createFs func() (filesystem.Fs, log.DebugLogger)
}

func (tc *testCases) runTests(t *testing.T) {
	t.Helper()
	// Call all Test* methods
	tp := reflect.TypeOf(tc)
	prefix := `Test`
	for i := 0; i < tp.NumMethod(); i++ {
		method := tp.Method(i)
		if strings.HasPrefix(method.Name, prefix) {
			fs, logger := tc.createFs()
			testName := strings.TrimPrefix(method.Name, prefix)
			t.Run(testName, func(t *testing.T) {
				t.Parallel()
				reflect.ValueOf(tc).MethodByName(method.Name).Call([]reflect.Value{
					reflect.ValueOf(t),
					reflect.ValueOf(fs),
					reflect.ValueOf(logger),
				})
			})
		}
	}
}

func (*testCases) TestApiName(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NotEmpty(t, fs.APIName())
}

func (*testCases) TestBasePath(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NotEmpty(t, fs.BasePath())
}

func (*testCases) TestWorkingDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.Equal(t, "my/dir", fs.WorkingDir())
	assert.NoError(t, fs.Mkdir("some/dir"))
	fs.SetWorkingDir("some/dir")
	assert.Equal(t, "some/dir", fs.WorkingDir())
	fs.SetWorkingDir("some/missing")
	assert.Equal(t, "", fs.WorkingDir())
}

func (*testCases) TestSubDirFs(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("sub/dir1/dir2/file.txt", "foo\n")))
	assert.True(t, fs.IsFile(`sub/dir1/dir2/file.txt`))

	// Empty path is not allowed
	_, err := fs.SubDirFs("  ")
	assert.Error(t, err)
	assert.Equal(t, `cannot get sub directory "  ": path cannot be empty`, err.Error())

	// /sub/dir1
	subDirFs1, err := fs.SubDirFs(`/sub/dir1`)
	assert.NoError(t, err)
	assert.Equal(t, `/`, subDirFs1.WorkingDir())
	assert.Equal(t, filepath.Join(fs.BasePath(), `sub`, `dir1`), subDirFs1.BasePath()) // nolint: forbidigo
	assert.False(t, subDirFs1.IsFile(`sub/dir1/dir2/file.txt`))
	assert.True(t, subDirFs1.IsFile(`dir2/file.txt`))
	file1, err := subDirFs1.ReadFile(filesystem.NewFileDef(`dir2/file.txt`))
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", file1.Content)

	// /sub/dir1/dir2
	subDirFs2, err := subDirFs1.SubDirFs(`/dir2`)
	assert.NoError(t, err)
	assert.Equal(t, `/`, subDirFs2.WorkingDir())
	assert.False(t, subDirFs2.IsFile(`sub/dir1/dir2/file.txt`))
	assert.False(t, subDirFs2.IsFile(`dir2/file.txt`))
	assert.True(t, subDirFs2.IsFile(`file.txt`))
	file2, err := subDirFs2.ReadFile(filesystem.NewFileDef(`file.txt`))
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", file2.Content)

	// file
	subDirFs3, err := subDirFs2.SubDirFs(`/file.txt`)
	assert.Error(t, err)
	assert.Equal(t, `cannot get sub directory "file.txt": path "file.txt" is not directory`, err.Error())
	assert.Nil(t, subDirFs3)

	// not found
	subDirFs4, err := subDirFs2.SubDirFs(`/abc`)
	assert.Error(t, err) // msg differs between backends
	assert.Nil(t, subDirFs4)

	// check working dir inheritance
	// original FS has working dir "my/dir"
	assert.Equal(t, filesystem.Join(`my`, `dir`), fs.WorkingDir())
	// create directory "my/dir/foo/bar"
	assert.NoError(t, fs.Mkdir(filesystem.Join(`my`, `dir`, `foo`, `bar`)))
	// get sub FS for "my" dir -> working dir is inherited "dir"
	subDirFs5, err := fs.SubDirFs(`my`)
	assert.NoError(t, err)
	assert.Equal(t, `dir`, subDirFs5.WorkingDir())
	// get sub FS for "my/dir" dir -> working dir is "/"
	subDirFs6, err := fs.SubDirFs(filesystem.Join(`my`, `dir`))
	assert.NoError(t, err)
	assert.Equal(t, ``, subDirFs6.WorkingDir())
	// get sub FS for "my/dir/foo/bar" dir -> working dir is "/"
	subDirFs7, err := fs.SubDirFs(filesystem.Join(`my`, `dir`, `foo`, `bar`))
	assert.NoError(t, err)
	assert.Equal(t, `/`, subDirFs7.WorkingDir())
}

func (*testCases) TestSetLogger(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	logger := log.NewNopLogger()
	assert.NotPanics(t, func() {
		fs.SetLogger(logger)
	})
}

func (*testCases) TestWalk(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("sub/dir1/dir2/dir3"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("sub/dir2/file.txt", "foo\n")))

	paths := make([]string, 0)
	root := "."
	err := fs.Walk(root, func(path string, info iofs.FileInfo, err error) error {
		// Skip root
		if root == path {
			return nil
		}

		assert.NoError(t, err)
		paths = append(paths, path)
		return err
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{
		`sub`,
		`sub/dir1`,
		`sub/dir1/dir2`,
		`sub/dir1/dir2/dir3`,
		`sub/dir2`,
		`sub/dir2/file.txt`,
	}, paths)
}

func (*testCases) TestGlob(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir2/file.txt", "foo\n")))

	matches, err := fs.Glob(`my/abc/*`)
	assert.NoError(t, err)
	assert.Empty(t, matches)

	matches, err = fs.Glob(`my/*`)
	assert.NoError(t, err)
	assert.Equal(t, []string{`my/dir1`, `my/dir2`}, matches)

	matches, err = fs.Glob(`my/*/*`)
	assert.NoError(t, err)
	assert.Equal(t, []string{`my/dir2/file.txt`}, matches)

	matches, err = fs.Glob(`my/*/*.txt`)
	assert.NoError(t, err)
	assert.Equal(t, []string{`my/dir2/file.txt`}, matches)
}

func (*testCases) TestStat(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir/file.txt", "foo\n")))
	s, err := fs.Stat(`my/dir/file.txt`)
	assert.NoError(t, err)
	assert.False(t, s.IsDir())
}

func (*testCases) TestReadDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir/subdir"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir/file.txt", "foo\n")))

	items, err := fs.ReadDir(`my/dir`)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}

func (*testCases) TestExists(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	// Create file
	filePath := "file.txt"
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(filePath, "foo\n")))

	// Assert
	logger.Truncate()
	assert.True(t, fs.Exists(filePath))
	assert.False(t, fs.Exists("file-non-exists.txt"))
	assert.Equal(t, "", logger.AllMessages())
}

func (*testCases) TestIsFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir/file.txt", "foo\n")))

	// Assert
	assert.True(t, fs.IsFile("my/dir/file.txt"))
	assert.False(t, fs.IsFile("my/dir"))
	assert.False(t, fs.IsFile("file-non-exists.txt"))
}

func (*testCases) TestIsDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir/file.txt", "foo\n")))

	// Assert
	assert.True(t, fs.IsDir("my/dir"))
	assert.False(t, fs.IsDir("my/dir/file.txt"))
	assert.False(t, fs.IsDir("file-non-exists.txt"))
}

func (*testCases) TestCreate(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	fd, err := fs.Create(`file.txt`)
	assert.NoError(t, err)

	n, err := fd.WriteString("foo\n")
	assert.Equal(t, 4, n)
	assert.NoError(t, err)
	assert.NoError(t, fd.Close())

	file, err := fs.ReadFile(filesystem.NewFileDef("file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", file.Content)
}

func (*testCases) TestOpen(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(`file.txt`, "foo\n")))

	fd, err := fs.Open(`file.txt`)
	assert.NoError(t, err)

	content, err := io.ReadAll(fd)
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", string(content))
	assert.NoError(t, fd.Close())
}

func (*testCases) TestOpenFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(`file.txt`, "foo\n")))

	fd, err := fs.OpenFile(`file.txt`, os.O_RDONLY, 0o600)
	assert.NoError(t, err)

	content, err := io.ReadAll(fd)
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", string(content))
	assert.NoError(t, fd.Close())
}

func (*testCases) TestMkdir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.False(t, fs.Exists("my/dir"))
	assert.NoError(t, fs.Mkdir("my/dir"))
	assert.True(t, fs.Exists("my/dir"))
	assert.NoError(t, fs.Mkdir("my/dir"))
	assert.True(t, fs.Exists("my/dir"))
}

func (*testCases) TestCopyFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.False(t, fs.Exists("my/dir2/file.txt"))

	assert.NoError(t, fs.Copy("my/dir1/file.txt", "my/dir2/file.txt"))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestCopyFileExists(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir2/file.txt", "bar\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
	err := fs.Copy("my/dir1/file.txt", "my/dir2/file.txt")
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), `cannot copy "my/dir1/file.txt" -> "my/dir2/file.txt": destination exists`))
}

func (*testCases) TestCopyForce(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir2/file.txt", "bar\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
	assert.NoError(t, fs.CopyForce("my/dir1/file.txt", "my/dir2/file.txt"))

	file, err := fs.ReadFile(filesystem.NewFileDef("my/dir2/file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", file.Content)
	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestCopyDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.False(t, fs.Exists("my/dir2/file.txt"))

	assert.NoError(t, fs.Copy("my/dir1", "my/dir2"))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestCopyDirExists(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.Mkdir("my/dir2"))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2"))
	err := fs.Copy("my/dir1", "my/dir2")
	assert.True(t, strings.HasPrefix(err.Error(), `cannot copy "my/dir1" -> "my/dir2": destination exists`))
}

func (*testCases) TestMoveFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.False(t, fs.Exists("my/dir2/file.txt"))

	assert.NoError(t, fs.Move("my/dir1/file.txt", "my/dir2/file.txt"))

	assert.False(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestMoveFileExists(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir2/file.txt", "bar\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
	err := fs.Move("my/dir1/file.txt", "my/dir2/file.txt")
	assert.Error(t, err)
	assert.Equal(t, `cannot move "my/{dir1/file.txt -> dir2/file.txt}": destination exists`, err.Error())
}

func (*testCases) TestMoveForce(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.Mkdir("my/dir2"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir2/file.txt", "bar\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
	assert.NoError(t, fs.MoveForce("my/dir1/file.txt", "my/dir2/file.txt"))

	file, err := fs.ReadFile(filesystem.NewFileDef("my/dir2/file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "foo\n", file.Content)
	assert.False(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestMoveDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.False(t, fs.Exists("my/dir2/file.txt"))

	assert.NoError(t, fs.Move("my/dir1", "my/dir2"))

	assert.False(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2/file.txt"))
}

func (*testCases) TestMoveDirExists(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))
	assert.NoError(t, fs.Mkdir("my/dir2"))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.True(t, fs.Exists("my/dir2"))
	err := fs.Move("my/dir1", "my/dir2")
	assert.Equal(t, `cannot move "my/{dir1 -> dir2}": destination exists`, err.Error())
}

func (*testCases) TestRemoveFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1/file.txt"))
	assert.NoError(t, fs.Remove("my/dir1/file.txt"))
	assert.False(t, fs.Exists("my/dir1/file.txt"))
}

func (*testCases) TestRemoveDir(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Mkdir("my/dir1"))
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile("my/dir1/file.txt", "foo\n")))

	assert.True(t, fs.Exists("my/dir1"))
	assert.NoError(t, fs.Remove("my/dir1"))
	assert.False(t, fs.Exists("my/dir1"))
}

func (*testCases) TestRemoveNotExist(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	assert.NoError(t, fs.Remove("my/dir1/file.txt"))
}

func (*testCases) TestReadFile(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	// Create file
	filePath := "file.txt"
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(filePath, "foo\n")))

	// Read
	logger.Truncate()
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "foo\n", file.Content)
	assert.Equal(t, `DEBUG  Loaded "file.txt"`, strings.TrimSpace(logger.AllMessages()))
}

func (*testCases) TestReadFileNotFound(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	filePath := "file.txt"
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.Error(t, err)
	assert.Nil(t, file)
	assert.True(t, strings.HasPrefix(err.Error(), `missing file "file.txt"`))
	assert.Equal(t, "", logger.AllMessages())
}

func (*testCases) TestWriteFile(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	filePath := "file.txt"

	// Write
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(filePath, "foo\n")))
	assert.Equal(t, `DEBUG  Saved "file.txt"`, strings.TrimSpace(logger.AllMessages()))

	// Read
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "foo\n", file.Content)
}

func (*testCases) TestWriteFileDirNotFound(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	filePath := "my/dir/file.txt"

	// Write
	assert.NoError(t, fs.WriteFile(filesystem.NewRawFile(filePath, "foo\n")))
	expectedLogs := `
DEBUG  Created directory "my/dir"
DEBUG  Saved "my/dir/file.txt"
`
	assert.Equal(t, strings.TrimSpace(expectedLogs), strings.TrimSpace(logger.AllMessages()))

	// Read - dir is auto created
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "foo\n", file.Content)
}

func (*testCases) TestWriteFile_JsonFile(t *testing.T, fs filesystem.Fs, logger log.DebugLogger) {
	filePath := "file.json"

	// Write
	data := orderedmap.New()
	data.Set(`foo`, `bar`)
	assert.NoError(t, fs.WriteFile(filesystem.NewJSONFile(filePath, data)))
	assert.Equal(t, `DEBUG  Saved "file.json"`, strings.TrimSpace(logger.AllMessages()))

	// Read
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "{\n  \"foo\": \"bar\"\n}\n", file.Content)
}

func (*testCases) TestCreateOrUpdateFile(t *testing.T, fs filesystem.Fs, _ log.DebugLogger) {
	filePath := "file.txt"

	// Create empty file
	updated, err := fs.CreateOrUpdateFile(filesystem.NewFileDef(filePath), []filesystem.FileLine{})
	assert.False(t, updated)
	assert.NoError(t, err)
	assert.True(t, fs.Exists(filePath))
	file, err := fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.Equal(t, "", file.Content)

	// Add some lines
	updated, err = fs.CreateOrUpdateFile(filesystem.NewFileDef(filePath), []filesystem.FileLine{
		{Line: "foo"},
		{Line: "bar\n"},
		{Line: "BAZ1=123\n", Regexp: "^BAZ1="},
		{Line: "BAZ2=456\n", Regexp: "^BAZ2=.*$"},
	})
	assert.NoError(t, err)
	assert.True(t, updated)
	assert.True(t, fs.Exists(filePath))
	file, err = fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.Equal(t, "foo\nbar\nBAZ1=123\nBAZ2=456\n", file.Content)

	// Update some lines
	updated, err = fs.CreateOrUpdateFile(filesystem.NewFileDef(filePath), []filesystem.FileLine{
		{Line: "BAZ1=new123\n", Regexp: "^BAZ1="},
		{Line: "BAZ2=new456\n", Regexp: "^BAZ2=.*$"},
	})
	assert.True(t, updated)
	assert.NoError(t, err)
	assert.True(t, fs.Exists(filePath))
	file, err = fs.ReadFile(filesystem.NewFileDef(filePath))
	assert.NoError(t, err)
	assert.Equal(t, "foo\nbar\nBAZ1=new123\nBAZ2=new456\n", file.Content)
}
