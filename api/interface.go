package api

import (
	"context"
	genlibapi "github.com/StephanHCB/go-generator-lib/api"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

type GitApiRepo interface {
	GetLocalPath() string
}

// Functionality that this library exposes.
type GitApi interface {
	// CreateTemporaryWorkdir creates a temporary working directory with a random name underneath basePath.
	//
	// You need to call this before any of the other methods can be called
	//
	// We use a random subdirectory so multiple goroutines can render in parallel
	CreateTemporaryWorkdir(ctx context.Context, basePath string) error

	// CloneSourceRepo clones the source repo into the working directory and switch to the given branch (or tag, or revision).
	CloneSourceRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// CloneTargetRepo clones the target repo into the working directory and sets up the given branch.
	//
	// If the branch does not yet exist, it will be created from the base branch (or tag, or revision),
	// otherwise we just check it out.
	CloneTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, baseBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// PrepareTargetRepo prepares the target repo inside the working directory.
	PrepareTargetRepo(ctx context.Context, gitRepoUrl string, gitBranch string, auth transport.AuthMethod) (GitApiRepo, error)

	// WriteRenderSpecFile writes the given parameters for the given generator to a render spec file in the target directory.
	//
	// unless some specific reason prevents you from this naming convention, renderSpecFile should be
	// 'generated-<generatorName>.yaml'
	//
	// the parameters will be validated against the generator spec file found in the source repo (called
	// 'generator-<generatorName>.yaml'). It is an error if any parameter does not conform to the specification,
	// or is missing and does not have a default, or if any parameters are unknown.
	//
	// Note that the render spec file in the target directory is silently overwritten. It is a git repo after all.
	// If the render spec file does not exist, that just means you are using the generator for the first time,
	// so it is silently created.
	//
	// Response is filled even in case of an error and will contain more details of what caused the error.
	WriteRenderSpecFile(ctx context.Context,
		generatorName string,
		renderSpecFile string,
		parameters map[string]interface{}) (*genlibapi.Response, error)

	// Generate generates files using the render spec file written by WriteRenderSpecFile.
	//
	// Response is filled even in case of an error and will contain more details of what caused the error
	// and what output files were affected. After a successful run, Response also contains the list of files
	// that were rendered.
	Generate(ctx context.Context) (*genlibapi.Response, error)

	// DeleteRenderSpecFile removes the render spec file written by WriteRenderSpecFile
	//
	// This is an optional step between Generate and CommitAndPush, in case you do not wish to commit the
	// render spec file.
	//
	// In most cases, deleting the render spec file is a BAD IDEA. The render spec file documents the parameters
	// used to perform the render run, so it is useful to check it in. Also, your user may wish to re-run
	// the render process with a new version of the template using the same parameters (or mostly the same parameters)
	// that were used during the original render run.
	//
	// There are however cases, where the generator incrementally adds a small number of files to a repository, and
	// it is undesirable to clutter the repository with lots of render spec files, and using the same filename
	// multiple times would create conflicts in git pull requests that otherwise do not touch the same files.
	DeleteRenderSpecFile(ctx context.Context) error

	// CommitAndPush commits the changes in the target and pushes them (if an auth method is supplied).
	CommitAndPush(ctx context.Context, name string, email string, message string, auth transport.AuthMethod) error

	// Cleanup deletes the temporary working directory, including the source and target clones underneath it.
	//
	// the base path given to CreateTemporaryWorkdir is left untouched, so it can be re-used for the next
	// (or concurrent) render operations.
	Cleanup(ctx context.Context) error
}
