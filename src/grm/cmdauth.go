package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"grm/config"
	"fmt"
)

func cmdAuth(cmd *cli.Cmd) {
	cmd.Spec = "NAME|--all [ -u=<username> ] [ -p=<password> ] [ --yes ]"

	var (
		name     = cmd.StringArg("NAME", "", "The name of the remote definition")
		username = cmd.StringOpt("u username", "", "The username to access Github")
		password = cmd.StringOpt("p password", "", "The password to access Github")
		yes      = cmd.BoolOpt("y yes", false, "Accept all questions with yes")
		all      = cmd.BoolOpt("all", false, "Re-authorize all remote definitions")
	)

	cmd.Action = func() {
		if *name == "" && !*all {
			log.Fatal("No remote name specified")
		}

		readOverride := func(definition string) bool {
			if *yes {
				return true
			}
			return readYesNoQuestion(fmt.Sprintf("You already have an authorization configuration for remote "+
				"definition '%s', are you sure to override?", definition), false)
		}

		definitions := []string{*name}
		if *all {
			definitions = configuration.NamedSections(config.Remote)
		}

		for _, definition := range definitions {
			specifier := config.ExtractSpecifier(definition)

			if configuration != nil {
				_, oku := configuration.NamedSectionGet(specifier, config.Remote, config.Username, "")
				_, okp := configuration.NamedSectionGet(specifier, config.Remote, config.Password, "")

				if oku && okp {
					if !readOverride(specifier) {
						// Stop execution
						fmt.Println("Configuration not changed")
						return
					}
				}
			}

			fmt.Println(fmt.Sprintf("Configure authorization information for remote definition: %s", specifier))
			realUsername := *username
			if realUsername == "" {
				realUsername = readLine("Username:", false, "")
			}

			realPassword := *password
			if realPassword == "" {
				realPassword = readLine("Password:", true, "")
			}

			encryptedPassword, salt := encrypt(realPassword, machineKey)

			configuration.ApplyChanges(func(mutator config.Mutator) {
				mutator.NamedSectionSet(specifier, config.Remote, config.Username, "", realUsername)
				mutator.NamedSectionSet(specifier, config.Remote, config.Password, "", encryptedPassword)
				mutator.NamedSectionSet(specifier, config.Remote, config.Salt, "", salt)
			})
		}
	}
}
