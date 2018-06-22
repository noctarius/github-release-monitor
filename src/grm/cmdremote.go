package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"strconv"
	"grm/config"
	"fmt"
	"strings"
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

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedSectionSet(*name, config.Remote, config.RemoteUser, "", *user)
			mutator.NamedSectionSet(*name, config.Remote, config.ShowPrivate, "", strconv.FormatBool(showPrivate))
			mutator.NamedSectionSet(*name, config.Remote, config.ReleasePattern, "", *releasePattern)
			mutator.NamedSectionSet(*name, config.Remote, config.RepositoryPattern, "", *repositoryPattern)
			mutator.NamedSectionSet(*name, config.Remote, config.MilestonePattern, "", *milestonePattern)
			mutator.NamedSectionSet(*name, config.Remote, config.DownloadUrl, "", *downloadUrl)
		})
	}
}

func cmdRemoteRemove(cmd *cli.Cmd) {
	cmd.Spec = "NAME"

	var (
		name = cmd.StringArg("NAME", "", "The name of the remote definition")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		readDeleteConfirm := func() bool {
			line := readLine(fmt.Sprintf("The configuration %s is about to be deleted. Do you "+
				"really want to continue? [yes|No]", *name), false)

			line = strings.ToLower(line)

			if line == "yes" || line == "y" || line == "true" {
				return true
			}
			return false
		}

		if t := configuration.NamedSection(*name, config.Remote); len(t) > 0 {
			if !readDeleteConfirm() {
				// Stop execution
				println("Configuration not changed")
				return
			}
		}

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedDelete(*name, config.Remote)
		})
	}
}
