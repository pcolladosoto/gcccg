package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gcccg",
		Short: "Go Conventional Commits Changelog Generator.",
		Long:  "Go Conventional Commits Changelog Generator.",

		// Note we need to configure logging here as flags are not parsed before
		// calling rootCmd.Execute on main... As this PreRun function is persistent,
		// it'll be inherited by subcommands
		PersistentPreRun: setupLogging,

		Run: func(cmd *cobra.Command, args []string) {
			repo, err := git.PlainOpen(repoPath)
			if err != nil {
				slog.Error("error opening the repo", "err", err)
				os.Exit(1)
			}

			sortedTags, err := getTags(repo)
			if err != nil {
				slog.Error("error getting repository tags", "err", err)
				os.Exit(1)
			}

			if len(sortedTags) < 2 {
				slog.Error("not enough tags defined", "nTags", len(sortedTags))
				os.Exit(1)
			}

			if strings.ToLower(fromTag) == "auto" {
				fromTag = sortedTags[len(sortedTags)-1]
			}

			if strings.ToLower(toTag) == "auto" {
				for i, tag := range sortedTags {
					if tag == fromTag {
						if i == 0 {
							slog.Error("the specified from tag is the first one")
							os.Exit(1)
						}
						toTag = sortedTags[i-1]
					}
				}

				if strings.ToLower(toTag) == "auto" {
					slog.Error("couldn't find a to tag")
					os.Exit(1)
				}
			}

			fromRef, err := repo.Tag(fromTag)
			if err != nil {
				slog.Error("error getting from tag", "err", err)
				os.Exit(1)
			}

			toRef, err := repo.Tag(toTag)
			if err != nil {
				slog.Error("error getting to tag", "err", err)
				os.Exit(1)
			}

			commits, err := repo.Log(&git.LogOptions{
				From:  fromRef.Hash(),
				To:    toRef.Hash(),
				Order: git.LogOrderCommitterTime,
			})
			if err != nil {
				slog.Error("error getting commits", "err", err)
				os.Exit(1)
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
					slog.Error("error parsing commit", "err", err, "hash", c.Hash, "message", strings.TrimSpace(c.Message))
					return nil
				}
				ccMessage, ok := pMessage.(*conventionalcommits.ConventionalCommit)
				if !ok {
					slog.Warn("not a conventional commit", "hash", c.Hash.String())
					return nil
				}

				commitMap, ok := pCommits[ParseCommitType(ccMessage.Type)]
				if !ok {
					pCommits[ParseCommitType(ccMessage.Type)] = map[string]string{}
					commitMap = pCommits[ParseCommitType(ccMessage.Type)]
				}
				commitMap[c.Hash.String()] = ccMessage.Description

				authors[c.Author.Email] = c.Author.Name

				return nil
			}); err != nil {
				slog.Error("error processing commits", "err", err)
				os.Exit(1)
			}

			if err := executeTemplate(ChangelogData{
				ReleaseName: fromTag,
				FromTag:     fromTag,
				ToTag:       toTag,

				Ci:    pCommits[ci],
				Fixes: pCommits[fix],
				Docs:  pCommits[docs],
				Build: pCommits[build],
				Feat:  pCommits[feat],
				Perf:  pCommits[perf],
				Test:  pCommits[test],

				Authors:    authors,
				AddEmail:   addEmail,
				GhMarkdown: ghMarkdown,
			}, toStdout, outputPath); err != nil {
				slog.Error("error executing the template", "err", err)
			}
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the built version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("built commit: %s\nbase version: %s\n", builtCommit, baseVersion)
		},
	}
)

func init() {
	// Disable completion please!
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add the different sub-commands
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("couldn't execute the root command", "err", err)
		os.Exit(1)
	}
}
