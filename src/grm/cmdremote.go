package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"strconv"
	"grm/config"
)

func cmdRemote(cmd *cli.Cmd) {
	cmd.Command("add", "Adds a remote Github user", cmdRemoteAdd)
	cmd.Command("remove", "Removes a remote Github user", cmdRemoteRemove)
}

func cmdRemoteAdd(cmd *cli.Cmd) {
	cmd.Spec = "NAME USER [ -p=<private> ] [ --release-pattern=<release-pattern> ] [ --repository-pattern=<repository-pattern> ] [ --milestone-pattern=<milestone-pattern> ] [ --download-url=<download-url> ]"

	var (
		name              = cmd.StringArg("NAME", "", "Defines the name of this remote")
		user              = cmd.StringArg("USER", "", "Defines the user to be registered")
		private           = cmd.BoolOpt("p private", false, "Analyze private repositories, default: false")
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
	// TODO
}