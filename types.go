package main

type commitType int

func ParseCommitType(cc string) commitType {
	switch cc {
	case "ci":
		return ci
	case "fix":
		return fix
	case "docs":
		return docs
	case "build":
		return build
	case "chore":
		return chore
	case "feat":
		return feat
	case "perf":
		return perf
	case "refactor":
		return refactor
	case "revert":
		return revert
	case "style":
		return style
	case "test":
		return test
	default:
		return invalid
	}
}

const (
	invalid commitType = iota
	ci
	fix
	docs
	build
	chore
	feat
	perf
	refactor
	revert
	style
	test
)
