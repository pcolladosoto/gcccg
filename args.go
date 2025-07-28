package main

var (
	fromTag    string
	toTag      string
	repoPath   string
	outputPath string
	addEmail   bool
	ghMarkdown bool
	toStdout   bool

	logLevelFlag string

	builtCommit string
	baseVersion string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&fromTag, "from-tag", "auto", "git tag to begin tracking commits from (auto selects the latest tag)")
	rootCmd.PersistentFlags().StringVar(&toTag, "to-tag", "auto", "git tag to end tracking commits on (auto selects the previous tag)")
	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", ".", "path to the local git repository")
	rootCmd.PersistentFlags().StringVar(&outputPath, "out", "", "path to the changelog output file; if empty it won't be written to a file.")
	rootCmd.PersistentFlags().StringVar(&logLevelFlag, "log-level", "info", "log level: one of debug, info, warn, error")
	rootCmd.PersistentFlags().BoolVar(&addEmail, "email", false, "add committer emails to the changelog?")
	rootCmd.PersistentFlags().BoolVar(&ghMarkdown, "gh", false, "use GitHub flavoured Markdown?")
	rootCmd.PersistentFlags().BoolVar(&toStdout, "stdout", false, "print the changelog to stdout?")
}
