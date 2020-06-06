package app

import (
	"context"
	"time"

	"github.com/jkuri/abstruse/internal/pkg/scm"
	"github.com/jkuri/abstruse/internal/pkg/shared"
	"github.com/jkuri/abstruse/internal/server/db/model"
	"github.com/jkuri/abstruse/internal/server/parser"
	jsoniter "github.com/json-iterator/go"
)

// TriggerBuild temp func.
func (app *App) TriggerBuild(repoID, userID uint) error {
	repo, err := app.repoRepository.Find(repoID, userID)
	if err != nil {
		return err
	}
	scm, err := scm.NewSCM(context.Background(), repo.Provider.Name, repo.Provider.URL, repo.Provider.AccessToken)
	if err != nil {
		return err
	}
	commit, err := scm.LastCommit(repo.FullName, repo.DefaultBranch)
	if err != nil {
		return err
	}
	content, err := scm.FindContent(repo.FullName, commit.Sha, ".abstruse.yml")
	if err != nil {
		return err
	}
	config := parser.ConfigParser{Raw: string(content.Data)}
	if err := config.Parse(); err != nil {
		return err
	}
	commandsJSON, err := jsoniter.Marshal(config.Commands)
	if err != nil {
		return err
	}

	buildModel := model.Build{
		Branch:          repo.DefaultBranch,
		Commit:          commit.Sha,
		CommitMessage:   commit.Message,
		Config:          string(content.Data),
		AuthorLogin:     commit.Author.Login,
		AuthorName:      commit.Author.Name,
		AuthorEmail:     commit.Author.Email,
		AuthorAvatar:    commit.Author.Avatar,
		CommitterLogin:  commit.Committer.Name,
		CommitterName:   commit.Committer.Name,
		CommitterEmail:  commit.Committer.Email,
		CommitterAvatar: commit.Committer.Avatar,
		RepositoryID:    repo.ID,
		StartTime:       func(t time.Time) *time.Time { return &t }(time.Now()),
	}
	build, err := app.buildRepository.Create(buildModel)
	if err != nil {
		return err
	}

	for _, env := range config.Env {
		jobModel := model.Job{
			Image:    config.Parsed.Image,
			Commands: string(commandsJSON),
			Env:      env,
			BuildID:  build.ID,
		}
		job, err := app.jobRepository.Create(jobModel)
		if err != nil {
			return err
		}
		if err := app.scheduleJob(job, repo.Provider, commit.Sha, repo.FullName); err != nil {
			return err
		}
	}

	if err := app.broadcastNewBuild(build.ID); err != nil {
		return err
	}

	return nil
}

// StopBuild stops the build and related jobs.
func (app *App) StopBuild(buildID uint) error {
	// build, err := app.buildRepository.FindAll(buildID)
	// if err != nil {
	// 	return err
	// }
	// for _, job := range build.Jobs {
	// 	go app.Scheduler.StopJob(job.ID)
	// }
	return nil
}

// RestartBuild stops the current build related jobs if any, then start them again.
func (app *App) RestartBuild(buildID uint) error {
	build, err := app.buildRepository.FindAll(buildID)
	if err != nil {
		return err
	}
	build.StartTime = func(t time.Time) *time.Time { return &t }(time.Now())
	build.EndTime = nil
	if build, err = app.buildRepository.Update(build); err != nil {
		return err
	}
	for _, job := range build.Jobs {
		// app.Scheduler.StopJob(job.ID)
		if err := app.scheduleJob(*job, build.Repository.Provider, build.Commit, build.Repository.FullName); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) scheduleJob(job model.Job, provider model.Provider, commitSHA, repo string) error {
	errch := make(chan error)

	go func() {
		j := shared.Job{
			ID:            job.ID,
			BuildID:       job.BuildID,
			Commands:      job.Commands,
			Image:         job.Image,
			Env:           job.Env,
			ProviderName:  provider.Name,
			ProviderURL:   provider.URL,
			ProviderToken: provider.AccessToken,
			CommitSHA:     commitSHA,
			RepoName:      repo,
			Priority:      uint16(1000),
			Status:        shared.StatusUnknown,
		}
		errch <- app.scheduler.scheduleJob(j)
	}()

	return <-errch
}

// func (app *App) saveJob(job *Job) error {
// 	j := model.Job{
// 		ID:      job.ID,
// 		BuildID: job.BuildID,
// 		Log:     strings.Join(job.Log, ""),
// 		Status:  job.Status,
// 	}
// 	if start, err := ptypes.Timestamp(job.Task.StartTime); err == nil {
// 		j.StartTime = &start
// 	}
// 	if end, err := ptypes.Timestamp(job.Task.EndTime); err == nil {
// 		j.EndTime = &end
// 	}
// 	j, err := app.jobRepository.Update(j)
// 	if err != nil {
// 		return err
// 	}
// 	build, err := app.buildRepository.FindAll(j.BuildID)
// 	if err != nil {
// 		return err
// 	}
// 	if build.EndTime == nil {
// 		alldone := true
// 		for _, j := range build.Jobs {
// 			if j.EndTime == nil {
// 				alldone = false
// 				break
// 			}
// 		}
// 		if alldone {
// 			build.EndTime = func(t time.Time) *time.Time { return &t }(time.Now())
// 			if _, err := app.buildRepository.Update(build); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return err
// }

func (app *App) updateBuildTime(buildID uint) error {
	build, err := app.buildRepository.FindAll(buildID)
	if err != nil {
		return err
	}
	if build.EndTime == nil {
		alldone := true
		for _, j := range build.Jobs {
			if j.EndTime == nil {
				alldone = false
				break
			}
		}
		if alldone {
			build.EndTime = func(t time.Time) *time.Time { return &t }(time.Now())
			if _, err := app.buildRepository.Update(build); err != nil {
				return err
			}
		}
	}
	return nil
}
