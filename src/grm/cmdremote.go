package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"strconv"
	"grm/config"
	"fmt"
)

func cmdRemote(cmd *cli.Cmd) {
	cmd.Command("add", "Adds a remote Github user", cmdRemoteAdd)
	cmd.Command("remove", "Removes a remote Github user", cmdRemoteRemove)
}

func cmdRemoteAdd(cmd *cli.Cmd) {
	cmd.Spec = "NAME USER [ -p=<private> ] [ --release-pattern=<release-pattern> ] [ --repository-pattern=<repository-pattern> ] [ --milestone-pattern=<milestone-pattern> ] [ --download-url=<download-url> ]"

	var (
		name              = cmd.StringArg("NAME", "", "The name of the remote definition")
		user              = cmd.StringArg("USER", "", "The remote user to be registered")
		private           = cmd.BoolOpt("p private", false, "Will analyze private repositories, default: false")
		releasePattern    = cmd.StringOpt("release-pattern", "", "The default pattern to match tag names")
		repositoryPattern = cmd.StringOpt("repository-pattern", "", "The default pattern to match repository names")
		milestonePattern  = cmd.StringOpt("milestone-pattern", "", "The default pattern to match milestone names")
		downloadUrl       = cmd.StringOpt("download-url", "", "The default download url pattern")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		if *name == "" {
			log.Fatal("No remote user specified")
		}

		showPrivate := *private

		realRepositoryPattern := *repositoryPattern
		if realRepositoryPattern == "" {
			realRepositoryPattern = readLine("Pattern to match repository names: [.*]",
				false, ".*")
		}

		realMilestonePattern := *milestonePattern
		if realMilestonePattern == "" {
			realMilestonePattern = readLine("Pattern to match milestone names: [^[a-zA-Z-_]-(.*)]",
				false, "^[a-zA-Z-_]-(.*)")
		}

		realReleasePattern := *releasePattern
		if realReleasePattern == "" {
			realReleasePattern = readLine("Pattern to match release names: []",
				false, "")
		}

		realDownloadUrl := *downloadUrl
		if realDownloadUrl == "" {
			fmt.Println("The basic download url is the template to generate download urls for releases.")
			fmt.Println("The template can contain template placeholders to be filled automatically.")
			fmt.Println("Current template placeholders are:")
			fmt.Println(" * {account}: The Github remote username (from the remote definition)")
			fmt.Println(" * {repository}: The name of the current repository")
			fmt.Println(" * {version}: The current release version")
			realDownloadUrl = readLine("Basic download url: [http://download.example.com/{account}/{repository}/{version}]",
				false, "http://download.example.com/{account}/{repository}/{version}")
		}

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedSectionSet(*name, config.Remote, config.RemoteUser, "", *user)
			mutator.NamedSectionSet(*name, config.Remote, config.ShowPrivate, "", strconv.FormatBool(showPrivate))
			mutator.NamedSectionSet(*name, config.Remote, config.ReleasePattern, "", realReleasePattern)
			mutator.NamedSectionSet(*name, config.Remote, config.RepositoryPattern, "", realRepositoryPattern)
			mutator.NamedSectionSet(*name, config.Remote, config.MilestonePattern, "", realMilestonePattern)
			mutator.NamedSectionSet(*name, config.Remote, config.DownloadUrl, "", realDownloadUrl)
		})
	}
}

func cmdRemoteRemove(cmd *cli.Cmd) {
	cmd.Spec = "NAME"

	var (
		name = cmd.StringArg("NAME", "", "The name of the remote definition")
		yes  = cmd.BoolOpt("y yes", false, "Accept all questions with yes")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		readDeleteConfirm := func() bool {
			if *yes {
				return true
			}
			return readYesNoQuestion(fmt.Sprintf("The configuration %s is about to be deleted. Do you "+
				"really want to continue?", *name), false)
		}

		if t := configuration.NamedSection(*name, config.Remote); len(t) > 0 {
			if !readDeleteConfirm() {
				// Stop execution
				fmt.Println("Configuration not changed")
				return
			}
		}

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedDelete(*name, config.Remote)
		})
	}
}
