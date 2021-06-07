package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
	"keboola-as-code/src/log"
	"keboola-as-code/src/options"
	"keboola-as-code/src/utils"
	"keboola-as-code/src/version"
	"os"
	"path"
	"regexp"
	"time"
)

const description = `
Keboola Connection pull/push client
for components configurations.

Configurations can be synchronized in both
directions [KBC project] <-> [a local directory].
`

const usageTemplate = `Usage:{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{else if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}

Aliases:`

type rootCommand struct {
	cmd          *cobra.Command
	options      *options.Options   // parsed flags and env variables
	prompt       *Prompt            // user interaction
	initialized  bool               // init method was called
	logFile      *os.File           // log file instance
	logFileClear bool               // is log file temporary? if yes, it will be removed at the end, if no error occurs
	logger       *zap.SugaredLogger // log to console and logFile
}

// NewRootCommand creates parent of all sub-commands
func NewRootCommand(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, prompt *Prompt) *rootCommand {
	root := &rootCommand{options: &options.Options{}, prompt: prompt}

	// Command definition
	root.cmd = &cobra.Command{
		Use:          path.Base(os.Args[0]), // name of the binary
		Version:      version.Version(),
		Short:        description,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Print help if no command specified
			return root.cmd.Help()
		},
	}

	// Setup in/out
	root.cmd.SetIn(stdin)
	root.cmd.SetOut(stdout)
	root.cmd.SetErr(stderr)

	// Setup templates
	root.cmd.SetVersionTemplate("{{.Version}}")
	root.cmd.SetUsageTemplate(
		regexp.MustCompile(`Usage:(.|\n)*Aliases:`).ReplaceAllString(root.cmd.UsageTemplate(), usageTemplate),
	)

	// Persistent flags for all sub-commands
	root.options.BindPersistentFlags(root.cmd.PersistentFlags())

	// Root command flags
	root.cmd.Flags().SortFlags = true
	root.cmd.Flags().Bool("version", false, "print version")

	// Init when flags are parsed
	root.cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return root.init(cmd)
	}

	// Sub-commands
	root.cmd.AddCommand(
		initCommand(root),
	)

	return root
}

// Execute command or sub-command
func (root *rootCommand) Execute() {
	defer root.tearDown()
	if err := root.cmd.Execute(); err != nil {
		// Init, it can be uninitialized, if error occurred before PersistentPreRun call
		_ = root.init(root.cmd)
		// Error is already logged
		os.Exit(1)
	}
}

func (root *rootCommand) GetCommandByName(name string) *cobra.Command {
	for _, cmd := range root.cmd.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}

	return nil
}

// tearDown makes clean-up after command execution
func (root *rootCommand) tearDown() {
	if err := recover(); err == nil {
		// No error -> remove log file if temporary
		if root.logFile != nil && root.logFileClear {
			if err = root.logFile.Close(); err != nil {
				panic(fmt.Errorf("cannot close log file \"%s\": %s", root.options.LogFilePath, err))
			}
			if err = os.Remove(root.options.LogFilePath); err != nil {
				panic(fmt.Errorf("cannot remove temp log file \"%s\": %s", root.options.LogFilePath, err))
			}
		}
	} else {
		// Error -> process and close log file
		exitCode := utils.ProcessPanic(err, root.logger, root.options.LogFilePath)
		if root.logFile != nil {
			if err = root.logFile.Close(); err != nil {
				panic(fmt.Errorf("cannot close log file \"%s\": %s", root.options.LogFilePath, err))
			}
		}
		os.Exit(exitCode)
	}
}

// init sets logger and options after flags are parsed
func (root *rootCommand) init(cmd *cobra.Command) (err error) {
	if root.initialized {
		return
	}
	root.initialized = true

	// Load values from flags and envs
	warnings, err := root.options.Load(cmd.Flags())

	// Setup logger and log options load warnings
	root.setupLogger()
	root.logDebugInfo()
	for _, msg := range warnings {
		root.logger.Debug(msg)
	}

	// Return load error
	return
}

// setupLogger according to the options
func (root *rootCommand) setupLogger() {
	logFile, logFileErr := root.getLogFile()
	root.logger = log.NewLogger(root.cmd.OutOrStdout(), root.cmd.ErrOrStderr(), logFile, root.options.Verbose)
	root.logFile = logFile
	root.cmd.SetOut(log.ToInfoWriter(root.logger))
	root.cmd.SetErr(log.ToWarnWriter(root.logger))

	// Warn if user specified log file and it cannot be opened
	if logFileErr != nil && !root.logFileClear {
		root.logger.Warnf("Cannot open log file: %s", logFileErr)
	}
}

func (root *rootCommand) logDebugInfo() {
	// Version
	_, err := log.ToDebugWriter(root.logger).WriteString(root.cmd.Version)
	if err != nil {
		panic(err)
	}

	// Command
	root.logger.Debugf("Running command %v", os.Args)

	// Options
	root.logger.Debug(root.options.Dump())
}

// Get log file defined in the flags or create a temp file
func (root *rootCommand) getLogFile() (logFile *os.File, logFileErr error) {
	if len(root.options.LogFilePath) > 0 {
		root.logFileClear = false // log file defined by user will be preserved
	} else {
		root.options.LogFilePath = path.Join(os.TempDir(), fmt.Sprintf("keboola-as-code-%d.txt", time.Now().Unix()))
		root.logFileClear = true // temp log file will be removed. It will be preserved only in case of error
	}

	logFile, logFileErr = os.OpenFile(root.options.LogFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if logFileErr != nil {
		root.options.LogFilePath = ""
	}
	return
}
