package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/spf13/cobra"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func getTags(repo *git.Repository) ([]string, error) {
	tagsIter, err := repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("error getting the tags: %w", err)
	}

	tags := []string{}
	tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref.Name().Short())
		return nil
	})

	sort.Strings(tags)
	return tags, nil
}

func setupLogging(cmd *cobra.Command, args []string) {
	logLevel, ok := logLevelMap[logLevelFlag]
	if !ok {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time.
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}

			// Remove the directory from the source's filename.
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}

			return a
		},
	}))
	slog.SetDefault(logger)
}
