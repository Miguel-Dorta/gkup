package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Args			[]string
	BackupName      string
	BackupDate      string
	BufferSize      int
	Cmd             string
	NumberOfThreads int
	OmitHidden      bool
	OmitErrors      bool
	ReadSymLinks    bool
	RepoPath        string
	Sum             string
	VerboseLevel    int

	ArgsErrors []error
)

func parseCmd(cmd *cobra.Command, args []string) {
	Cmd = cmd.Name()
	Args = args
}

func addFlagBackupName(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&BackupName, "name", "n", "", "backup name")
}

func addFlagBackupDate(cmd *cobra.Command) {
	d := *cmd.Flags().StringP("date", "d", "", `date of the backup to restore.
	It format must be YYYY-MM-DD_hh-mm-ss
	and must match with the backup date.`)
	if len(d) != 19 {
		ArgsErrors = append(ArgsErrors, errors.New("invalid backup date"))
	}
	BackupDate = d
}

func addFlagBufferSize(cmd *cobra.Command) {
	bSize := *cmd.Flags().IntP("buffer-size", "b", 4 * 1024 * 1024, "buffer size, in bytes, per thread")
	if bSize < 512 {
		bSize = 512
	}
	BufferSize = bSize
}

func addFlagNumberOfThreads(cmd *cobra.Command) {
	threads := *cmd.Flags().IntP("threads", "t", runtime.NumCPU(), "number of threads in parallel operations")
	if threads < 1 {
		ArgsErrors = append(ArgsErrors, errors.New("invalid number of threads"))
	}
	NumberOfThreads = threads
}

func addFlagOmitHidden(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&OmitHidden, "omit-hidden", false, "omit hidden files")
}

func addFlagReadSymLinks(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&ReadSymLinks, "read-symlinks", false, "read symlinks (will not avoid infinite loops)")
}

func addFlagSum(cmd *cobra.Command) {
	sum := *cmd.Flags().StringP("sum", "s", "", `hash algorithm used in this repository. It cannot be changed later.
Supported algorithms:
    - MD5
    - SHA1
    - SHA256 (default)
    - SHA512
    - SHA3-256
    - SHA3-512`)
	if sum == "" {
		sum = "sha256"
	}
	Sum = sum
}

func addPersistentFlagOmitErrors(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&OmitErrors, "omit-errors", false, "omit non-critical errors")
}

func addPersistentFlagRepoPath(cmd *cobra.Command) {
	repoPath := *cmd.PersistentFlags().StringP("repo", "r", "", `path of your repository
    If not provided, working directory will be used`)
	if repoPath == "" {
		repoPath = "."
	}
	RepoPath = repoPath
}

func addPersistentFlagVerboseLevel(cmd *cobra.Command) {
	verbose := *cmd.PersistentFlags().IntP("verbose", "v", 2, `verbose level
    0: no output
    1: critical messages only
    2: critical and error messages
    3: detailed info and errors
    4: debug
   `)
	if verbose < 0 || verbose > 4 {
		ArgsErrors = append(ArgsErrors, errors.New("invalid verbose level"))
	}
	VerboseLevel = verbose
}
