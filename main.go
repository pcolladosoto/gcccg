package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
	"github.com/spf13/cobra"

	_ "embed"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&fromTag, "from-tag", "v1.0.0", "git tag to begin tracking commits from")
	rootCmd.PersistentFlags().StringVar(&toTag, "to-tag", "auto", "git tag to end tracking commits on (auto selects the previous tag)")
	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", ".", "path to the local git repository")
	rootCmd.PersistentFlags().StringVar(&logLevelFlag, "log-level", "info", "log level: one of debug, info, warn, error")
	rootCmd.PersistentFlags().BoolVar(&addEmail, "email", false, "add committer emails to the changelog?")
	rootCmd.PersistentFlags().BoolVar(&ghMarkdown, "gh", false, "use GitHub flavoured Markdown?")
}

type ChangelogData struct {
	ReleaseName string
	FromTag     string
	ToTag       string
	Fixes       map[string]string
	Ci          map[string]string
	Docs        map[string]string
	Authors     map[string]string
	AddEmail    bool
	GhMarkdown  bool
}

type commitType int

func ParseCommitType(cc string) commitType {
	switch cc {
	case "ci":
		return ci
	case "fix":
		return fix
	case "docs":
		return docs
	default:
		return invalid
	}
}

const (
	invalid commitType = iota
	ci
	fix
	docs
)

var (
	rootCmd = &cobra.Command{
		Use:   "gcccg",
		Short: "Go Conventional Commits Changelog Generator.",
		Long:  "Go nuts!",

		// Note we need to configure logging here as flags are not parsed before
		// calling rootCmd.Execute on main... As this PreRun function is persistent,
		// it'll be inherited by subcommands
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logLevel, ok := logLevelMap[logLevelFlag]
			if !ok {
				logLevel = slog.LevelInfo
			}

			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     logLevel,
			}))
			slog.SetDefault(logger)
		},
		Run: func(cmd *cobra.Command, args []string) {
			repo, err := git.PlainOpen(repoPath)
			if err != nil {
				fmt.Printf("error opening the repo: %v\n", err)
				return
			}

			fromRef, err := repo.Tag(fromTag)
			if err != nil {
				fmt.Printf("error getting the from tag: %v\n", err)
				return
			}

			if strings.ToLower(toTag) == "auto" {
				tagsIter, err := repo.Tags()
				if err != nil {
					fmt.Printf("error getting the tags: %v\n", err)
					return
				}

				tags := []string{}
				tagsIter.ForEach(func(ref *plumbing.Reference) error {
					tags = append(tags, ref.Name().Short())
					return nil
				})

				sort.Strings(tags)
				for i, tag := range tags {
					if tag == fromTag {
						if i == 0 {
							fmt.Printf("the specified from tag is the first one!\n")
							return
						}
						toTag = tags[i-1]
					}
				}

				if strings.ToLower(toTag) == "auto" {
					fmt.Printf("couldn't find a to tag\n")
					return
				}
			}

			toRef, err := repo.Tag(toTag)
			if err != nil {
				fmt.Printf("error getting the to tag: %v\n", err)
				return
			}

			commits, err := repo.Log(&git.LogOptions{From: fromRef.Hash(), To: toRef.Hash()})
			if err != nil {
				fmt.Printf("error getting commits: %v\n", err)
			}

			opts := []conventionalcommits.MachineOption{
				parser.WithTypes(conventionalcommits.TypesConventional),
			}
			pMachine := parser.NewMachine(opts...)

			authors := map[string]string{}
			pCommits := map[commitType]map[string]string{}
			if err := commits.ForEach(func(c *object.Commit) error {
				if c.Hash == toRef.Hash() {
					return nil
				}

				pMessage, err := pMachine.Parse([]byte(strings.TrimSpace(c.Message)))
				if err != nil {
					fmt.Printf("error parsing commit %s (%s): %v\n", c.Hash, strings.TrimSpace(c.Message), err)
					return nil
				}
				ccMessage, ok := pMessage.(*conventionalcommits.ConventionalCommit)
				if !ok {
					fmt.Errorf("not a conventional commit")
				}
				// fmt.Printf("hash: %s, type: %s, description: %s\n", c.Hash, ccMessage.Type, ccMessage.Description)

				commitMap, ok := pCommits[ParseCommitType(ccMessage.Type)]
				if !ok {
					pCommits[ParseCommitType(ccMessage.Type)] = map[string]string{}
					commitMap = pCommits[ParseCommitType(ccMessage.Type)]
				}
				commitMap[c.Hash.String()] = ccMessage.Description

				authors[c.Author.Email] = c.Author.Name

				return nil
			}); err != nil {
				fmt.Printf("error processing commits: %v\n", err)
			}

			// fmt.Printf("%v\n", pCommits)

			parsedTemplate, err := template.New("changelog").Parse(changelogTemplateRaw)
			parsedTemplate.Execute(os.Stdout, ChangelogData{
				ReleaseName: fromTag,
				FromTag:     fromTag,
				ToTag:       toTag,
				Ci:          pCommits[ci],
				Fixes:       pCommits[fix],
				Docs:        pCommits[docs],
				Authors:     authors,
				AddEmail:    addEmail,
				GhMarkdown:  ghMarkdown,
			})
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the built version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("built commit: %s\nbase version: %s\n", builtCommit, baseVersion)
		},
	}

	tmpCmd = &cobra.Command{
		Use:   "tmp",
		Short: "A temporary command...",
	}

	fromTag    string
	toTag      string
	repoPath   string
	addEmail   bool
	ghMarkdown bool

	logLevelFlag string

	builtCommit string
	baseVersion string

	//go:embed changelog.tmpl
	changelogTemplateRaw string

	logLevelMap = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

func init() {
	// Disable completion please!
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add the different sub-commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(tmpCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
