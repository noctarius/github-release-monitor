package main

import (
	"github.com/jawher/mow.cli"
	"github.com/google/go-github/github"
	"log"
	"context"
	"fmt"
	"strconv"
	"regexp"
	"time"
	"github.com/araddon/dateparse"
	"sync"
	"github.com/vbauerster/mpb"
	"github.com/wushilin/threads"
	"github.com/vbauerster/mpb/decor"
	"strings"
	"net/http"
	"grm/config"
)

func cmdReport(cmd *cli.Cmd) {
	cmd.Spec = "NAME  [ -p=<private_repos> ] [ --repository-pattern=<repository-pattern> ] [ --since=<since> ]"

	var (
		name              = cmd.StringArg("NAME", "", "The name of the remote definition")
		private           = cmd.BoolOpt("p private", false, "Analyze private repositories, default: false")
		repositoryPattern = cmd.StringOpt("repository-pattern", "", "A pattern to match repository names")
		since             = cmd.StringOpt("since", "", "Date of search begin in ISO format YYYY-MM-DD")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No remote name specified")
		}

		username, ok := configuration.NamedSectionGet(*name, config.Remote, config.Username, "")
		if !ok {
			log.Fatal("Could not retrieve username from config, please run 'grm init'")
		}
		pass, ok := configuration.NamedSectionGet(*name, config.Remote, config.Password, "")
		if !ok {
			log.Fatal("Could not retrieve password from config, please run 'grm init'")
		}

		salt, ok := configuration.NamedSectionGet(*name, config.Remote, config.Salt, "")
		if !ok {
			log.Fatal("Could not retrieve salt from config, please run 'grm init'")
		}

		basicAuth := github.BasicAuthTransport{
			Username: username,
			Password: decrypt(pass, salt, machineKey),
		}

		client := github.NewClient(basicAuth.Client())

		remoteAccount := username
		showPrivate := *private
		repositoryPattern := *repositoryPattern
		if r, ok := configuration.NamedSectionGet(*name, config.Remote, config.RepositoryPattern, ""); ok {
			repositoryPattern = r
		}

		date := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		if *since != "" {
			d, err := dateparse.ParseIn(*since, time.UTC)
			if err != nil {
				log.Fatal("Could not parse since data", err)
			}
			date = d
		}

		if u, ok := configuration.NamedSectionGet(*name, config.Remote, config.RemoteUser, ""); ok {
			remoteAccount = u
		}
		if p, ok := configuration.NamedSectionGet(*name, config.Remote, config.ShowPrivate, ""); ok {
			sp, err := strconv.ParseBool(p)
			if err != nil {
				showPrivate = false
			} else {
				showPrivate = sp
			}
		}

		visibility := "public"
		if showPrivate {
			visibility = "all"
		}

		print("Reading repositories... ")
		repos := readRepositories(remoteAccount, visibility, repositoryPattern, client)
		println("done.")

		repositories := selectRepositories(repos, remoteAccount, date, client)
		println(fmt.Sprintf("Found %d repositories", len(repositories)))
		for _, rep := range repositories {
			for _, rel := range rep.releases {
				if rel.milestone != nil {
					println(fmt.Sprintf("New %s release: %s (%s)", rep.name, rel.name, rel.created.Format("2006-01-02")))
					println("Release Notes: " + rel.milestoneUrl)
					if rel.downloadUrl != "" {
						println("Download: " + rel.downloadUrl)
					}
					println("")
				}
			}
		}
	}
}

func readMilestones(account, repository string, client *github.Client) []*github.Milestone {
	ctx := context.Background()

	milestones := make([]*github.Milestone, 0)

	page := 1
	for {
		s, response, err := client.Issues.ListMilestones(ctx, account, repository, &github.MilestoneListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		})

		if rateLimit(response) {
			continue
		}

		if err != nil {
			log.Fatal(fmt.Sprintf("Could not retrieve commit for repository %s", repository), err)
		}

		for _, milestone := range s {
			milestones = append(milestones, milestone)
		}

		if hasMorePages(response) {
			page++
			continue
		}

		return milestones
	}
}

func selectRepositories(repositories []*github.Repository, account string, since time.Time, client *github.Client) []*repository {
	tasks := new(sync.WaitGroup)
	tasks.Add(len(repositories))

	p := mpb.New(mpb.WithWaitGroup(tasks))
	bar := p.AddBar(int64(len(repositories)-1),
		mpb.PrependDecorators(
			decor.Name("Filtering repositories", decor.WCSyncSpaceR),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
		),
	)

	jobs := make([]func() interface{}, len(repositories))
	for i, repo := range repositories {
		repoName := repo.GetName()
		jobs[i] = func() interface{} {
			milestones := readMilestones(account, repoName, client)
			tags := readTags(account, repoName, client)
			releases := filterTags(tags, account, repoName, since, client)

			var pattern *regexp.Regexp = nil
			milestonePattern, ok := configuration.NamedSectionGet(account, config.Remote, config.MilestonePattern, repoName)
			if !ok {
				log.Fatal("No milestone pattern defined to extract milestone naming scheme")
			}
			pattern, err := regexp.Compile(milestonePattern)
			if err != nil {
				log.Fatal(fmt.Sprintf("Cannot compile regex: %s", milestonePattern))
			}

			downloadUrl, _ := configuration.NamedSectionGet(account, config.Remote, config.DownloadUrl, repoName)

			for _, release := range releases {
				milestone := findMatchingMilestone(release, milestones, pattern)
				if milestone != nil {
					release.milestone = milestone
					release.milestoneUrl = fmt.Sprintf("%s?closed=1", milestone.GetHTMLURL())
					release.milestoneState = milestone.GetState()
					release.downloadUrl = buildDownloadUrl(account, repoName, downloadUrl, milestone)
				}
			}

			if len(releases) > 0 {
				rep := &repository{
					name:     repoName,
					releases: releases,
				}

				bar.Increment()
				tasks.Done()
				return rep
			}
			bar.Increment()
			tasks.Done()
			return nil
		}
	}

	futureGroup := threads.ParallelDoWithLimit(jobs, 8)
	ret := futureGroup.WaitAll()

	reps := make([]*repository, 0)
	for _, r := range ret {
		if r != nil {
			reps = append(reps, r.(*repository))
		}
	}

	return reps
}

func buildDownloadUrl(account, repository, downloadUrl string, milestone *github.Milestone) string {
	downloadUrl = strings.Replace(downloadUrl, "{name}", account, -1)
	downloadUrl = strings.Replace(downloadUrl, "{repository}", repository, -1)
	downloadUrl = strings.Replace(downloadUrl, "{version}", milestone.GetTitle(), -1)
	response, err := http.Get(downloadUrl)
	if err != nil {
		log.Fatal("Cannot test download url")
	}
	if response.StatusCode == http.StatusOK {
		return downloadUrl
	}
	return ""
}

func findMatchingMilestone(release *release, milestones []*github.Milestone, pattern *regexp.Regexp) *github.Milestone {
	substrings := pattern.FindAllStringSubmatch(release.name, 1)
	if len(substrings) > 0 && len(substrings[0]) > 1 {
		milestoneName := substrings[0][1]
		for _, milestone := range milestones {
			if milestone.GetTitle() == milestoneName {
				return milestone
			}
		}
	}

	return nil
}

func filterTags(tags []*github.RepositoryTag, account, repository string, since time.Time, client *github.Client) []*release {
	filteredTags := make([]*release, 0)
	for _, tag := range tags {
		commit := readCommit(account, repository, tag.GetCommit().GetSHA(), client)
		if since.Before(commit.GetCommit().GetCommitter().GetDate()) {
			filteredTags = append(filteredTags, &release{
				created: commit.GetCommit().GetCommitter().GetDate(),
				name:    tag.GetName(),
			})
		}
	}

	return filteredTags
}

func readCommit(account, repository, sha string, client *github.Client) *github.RepositoryCommit {
	ctx := context.Background()

	commit, response, err := client.Repositories.GetCommit(ctx, account, repository, sha)
	for {
		if rateLimit(response) {
			continue
		}

		if err != nil {
			log.Fatal(fmt.Sprintf("Could not retrieve commit for commitId %s", sha), err)
		}

		return commit
	}
}

func readTags(account, repository string, client *github.Client) []*github.RepositoryTag {
	ctx := context.Background()

	releases := make([]*github.RepositoryTag, 0)

	var pattern *regexp.Regexp = nil
	if r, ok := configuration.NamedSectionGet(account, config.Remote, config.ReleasePattern, repository); ok {
		p, err := regexp.Compile(r)
		if err != nil {
			log.Fatal(fmt.Sprintf("Cannot compile regex: %s", r))
		}
		pattern = p
	}

	page := 1
	for {
		r, response, err := client.Repositories.ListTags(ctx, account, repository, &github.ListOptions{
			PerPage: 100,
			Page:    page,
		})

		if rateLimit(response) {
			continue
		}

		if err != nil {
			log.Fatal(fmt.Sprintf("Could not retrieve tags for repository %s", repository), err)
		}

		for _, release := range r {
			if pattern != nil && !pattern.MatchString(release.GetName()) {
				continue
			}
			releases = append(releases, release)
		}

		if hasMorePages(response) {
			page++
			continue
		}

		return releases
	}
}

func readRepositories(account, visibility, repositoryPattern string, client *github.Client) []*github.Repository {
	ctx := context.Background()

	repositories := make([]*github.Repository, 0)

	var pattern *regexp.Regexp = nil
	if repositoryPattern != "" {
		p, err := regexp.Compile(repositoryPattern)
		if err != nil {
			log.Fatal(fmt.Sprintf("Cannot compile regex: %s", repositoryPattern))
		}
		pattern = p
	}

	page := 1
	for {
		r, response, err := client.Repositories.List(ctx, account, &github.RepositoryListOptions{
			Visibility:  visibility,
			Type:        "owner",
			Affiliation: "owner",
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		})

		if rateLimit(response) {
			continue
		}

		if err != nil {
			log.Fatal("Could not retrieve repositories", err)
		}

		for _, repository := range r {
			if pattern == nil || pattern.MatchString(repository.GetName()) {
				if !isBlacklisted(account, repository.GetName()) {
					repositories = append(repositories, repository)
				}
			}
		}

		if hasMorePages(response) {
			page++
			continue
		}

		return repositories
	}
}

func isBlacklisted(account, repository string) bool {
	if r, ok := configuration.NamedSectionGet(account, config.Remote, config.RepositoryBlacklisted, repository); ok {
		b, err := strconv.ParseBool(r)
		if err != nil {
			log.Fatal("Could not parse boolean", err)
		}
		return b
	}
	return false
}

type repository struct {
	name     string
	releases []*release
	url      string
}

type release struct {
	name           string
	created        time.Time
	milestoneUrl   string
	milestoneState string
	downloadUrl    string
	milestone      *github.Milestone
}
