package run

const (
	SourceAPI       Source = "tfe-api"
	SourceUI        Source = "tfe-ui"
	SourceTerraform Source = "terraform+cloud"
	SourceGithub    Source = "github"
	SourceGitlab    Source = "gitlab"
)

// Source represents a source type of a run.
type Source string
