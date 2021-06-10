package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"keboola-as-code/src/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkingDirFromOs(t *testing.T) {
	options := &Options{}
	flags := &pflag.FlagSet{}

	// Load
	_, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, wd, options.WorkingDirectory())
}

func TestWorkingDirFromFlag(t *testing.T) {
	tempDir := t.TempDir()
	flags := &pflag.FlagSet{}
	options := &Options{}
	options.BindPersistentFlags(flags)
	assert.NoError(t, flags.Set("working-dir", tempDir))

	// Load
	_, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, tempDir, options.WorkingDirectory())
}

func TestProjectDirNotFound(t *testing.T) {
	// Load
	options := &Options{}
	flags := &pflag.FlagSet{}
	_, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	assert.Empty(t, options.projectDirectory)
	assert.False(t, options.HasProjectDirectory())
}

func TestProjectDirExpectedDirButFoundFile(t *testing.T) {
	projectDir := t.TempDir()
	metadataDir := filepath.Join(projectDir, ".keboola")
	fakeMetadataFile := filepath.Join(projectDir, "foo", ".keboola")
	workingDir := filepath.Join(projectDir, "foo", "bar", "baz")

	// Create ".keboola" dir in project directory
	assert.NoError(t, os.Mkdir(metadataDir, 0600))

	// Working dir = project sub-dir
	assert.NoError(t, os.MkdirAll(workingDir, 0600))
	assert.NoError(t, os.Chdir(workingDir))

	// Create ".keboola" file in "foo" dir -> invalid ".keboola" should be dir
	file, err := os.Create(fakeMetadataFile)
	assert.NoError(t, err)
	assert.NoError(t, file.Close())

	// Load
	options := &Options{}
	flags := &pflag.FlagSet{}
	warnings, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Equal(t, []string{fmt.Sprintf("Expected dir, but found file at \"%s\"", fakeMetadataFile)}, warnings)
}

func TestProjectDirSameAsWorkingDir(t *testing.T) {
	projectDir := t.TempDir()
	metadataDir := filepath.Join(projectDir, ".keboola")

	// Create ".keboola" dir in project directory
	assert.NoError(t, os.Mkdir(metadataDir, 0600))

	// Working dir = project dir
	assert.NoError(t, os.Chdir(projectDir))

	// Load
	options := &Options{}
	flags := &pflag.FlagSet{}
	warnings, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Empty(t, warnings)
}

func TestProjectDirIsParentOfWorkingDir(t *testing.T) {
	projectDir := t.TempDir()
	metadataDir := filepath.Join(projectDir, ".keboola")
	workingDir := filepath.Join(projectDir, "foo", "bar", "baz")

	// Create ".keboola" dir in project directory
	assert.NoError(t, os.Mkdir(metadataDir, 0600))

	// Working dir = project dir sub-dir
	assert.NoError(t, os.MkdirAll(workingDir, 0600))
	assert.NoError(t, os.Chdir(workingDir))

	// Load
	options := &Options{}
	flags := &pflag.FlagSet{}
	warnings, err := options.Load(flags)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Empty(t, warnings)
}

func TestValuesPriority(t *testing.T) {
	defer utils.ResetEnv(t, os.Environ())

	// Create working and project dir
	projectDir := t.TempDir()
	metadataDir := filepath.Join(projectDir, ".keboola")
	workingDir := filepath.Join(projectDir, "foo", "bar")
	assert.NoError(t, os.MkdirAll(workingDir, 0600))
	assert.NoError(t, os.Chdir(workingDir))

	// Create structs
	flags := &pflag.FlagSet{}
	options := &Options{}
	options.BindPersistentFlags(flags)

	// No values defined
	warnings, err := options.Load(flags)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Equal(t, workingDir, options.WorkingDirectory())
	assert.False(t, options.HasProjectDirectory())
	assert.Equal(t, "", options.ApiHost)

	// 1. Lowest priority, ".env" file from project dir
	os.Clearenv()
	assert.NoError(t, os.Mkdir(metadataDir, 0600))
	file, err := os.Create(filepath.Join(projectDir, ".env"))
	assert.NoError(t, err)
	_, err = file.WriteString("KBC_STORAGE_API_HOST=connection.keboola.com")
	assert.NoError(t, file.Close())
	assert.NoError(t, err)
	warnings, err = options.Load(flags)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Equal(t, workingDir, options.WorkingDirectory())
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Equal(t, "connection.keboola.com", options.ApiHost)

	// 2. Higher priority, ".env" file from working dir
	os.Clearenv()
	file, err = os.Create(filepath.Join(workingDir, ".env"))
	assert.NoError(t, err)
	_, err = file.WriteString("KBC_STORAGE_API_HOST=connection.north-europe.azure.keboola.com/")
	assert.NoError(t, file.Close())
	assert.NoError(t, err)
	warnings, err = options.Load(flags)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Equal(t, workingDir, options.WorkingDirectory())
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Equal(t, "connection.north-europe.azure.keboola.com", options.ApiHost)

	// 3. Higher priority , ENV defined in OS
	os.Clearenv()
	assert.NoError(t, os.Setenv("KBC_STORAGE_API_HOST", "https://connection.eu-central-1.keboola.com/"))
	warnings, err = options.Load(flags)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Equal(t, workingDir, options.WorkingDirectory())
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Equal(t, "connection.eu-central-1.keboola.com", options.ApiHost)

	// 4. The highest priority, flag
	assert.NoError(t, flags.Set("storage-api-host", "connection.keboola.cloud"))
	warnings, err = options.Load(flags)
	assert.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Equal(t, workingDir, options.WorkingDirectory())
	assert.Equal(t, projectDir, options.ProjectDirectory())
	assert.Equal(t, "connection.keboola.cloud", options.ApiHost)
}

func TestValidateNoRequired(t *testing.T) {
	options := &Options{}
	assert.Empty(t, options.Validate([]string{}))
}

func TestValidateAllRequired(t *testing.T) {
	options := &Options{}
	errors := options.Validate([]string{"projectDirectory", "ApiHost", "ApiToken"})

	// Assert
	expected := []string{
		`- This or any parent directory is not a Keboola project dir.`,
		`  Project directory must contain the ".keboola" metadata directory.`,
		`  Please change working directory to a project directory or use the "init" command.`,
		`- Missing api host. Please use "--storage-api-host" flag or ENV variable "KBC_STORAGE_API_HOST".`,
		`- Missing api token. Please use "--storage-api-token" flag or ENV variable "KBC_STORAGE_API_TOKEN".`,
	}
	assert.Equal(t, strings.Join(expected, "\n"), errors)
}

func TestDump(t *testing.T) {
	options := &Options{}
	options.ApiHost = "connection.keboola.com"
	options.ApiToken = "12345-67890123abcd"
	expected := `Parsed options: &options.Options{` +
		`Verbose:false, ` +
		`VerboseApi:false, ` +
		`LogFilePath:"", ` +
		`ApiHost:"connection.keboola.com", ` +
		`ApiToken:"12345-6*****", ` +
		`workingDirectory:"", ` +
		`projectDirectory:"", ` +
		`metadataDirectory:""` +
		`}`
	assert.Equal(t, expected, options.Dump())
}
