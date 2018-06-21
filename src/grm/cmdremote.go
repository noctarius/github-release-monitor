package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"github.com/zieckey/goini"
	"strconv"
)

func cmdremote(cmd *cli.Cmd) {
	cmd.Command("add", "Adds a remoteGithub user", cmdremoteadd)
}

func cmdremoteadd(cmd *cli.Cmd) {
	cmd.Spec = "NAME USER [ -p=<private> ] [ --release-pattern=<release-pattern> ] [ --repository-pattern=<repository-pattern> ]"

	var (
		name              = cmd.StringArg("NAME", "", "Defines the name of this remote")
		user              = cmd.StringArg("USER", "", "Defines the user to be registered")
		private           = cmd.BoolOpt("p private", false, "Analyze private repositories, default: false")
		releasePattern    = cmd.StringOpt("release-pattern", "", "A pattern to match tag names")
		repositoryPattern = cmd.StringOpt("repository-pattern", "", "A pattern to match repository names")
		milestonePattern  = cmd.StringOpt("milestone-pattern", "", "A pattern to match milestone names")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No NAME variable given")
		}

		if *user == "" {
			log.Fatal("No USER variable given")
		}

		showPrivate := *private
		section := buildRemoteSection(*name)

		changeConfig(func(config *goini.INI) {
			config.SectionSet(section, keyRemoteUser, *user)
			config.SectionSet(section, keyShowPrivate, strconv.FormatBool(showPrivate))
			config.SectionSet(section, keyReleasePattern, *releasePattern)
			config.SectionSet(section, keyRepositoryPattern, *repositoryPattern)
			config.SectionSet(section, keyMilestonePattern, *milestonePattern)
		})
	}
}
