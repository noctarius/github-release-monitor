package main

import (
	"github.com/jawher/mow.cli"
	"strings"
	"log"
	"grm/config"
	"fmt"
)

func cmdAuth(cmd *cli.Cmd) {
	cmd.Spec = "NAME [ -u=<username> ] [ -p=<password> ]"

	var (
		name     = cmd.StringArg("NAME", "", "The name of the remote definition")
		username = cmd.StringOpt("u username", "", "The username to access Github")
		password = cmd.StringOpt("p password", "", "The password to access Github")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No remote name specified")
		}

		readOverride := func() bool {
			line := readLine("You already have an authorization configuration, are you sure "+
				"to override? [yes|No] ", false)

			line = strings.ToLower(line)

			if line == "yes" || line == "y" || line == "true" {
				return true
			}
			return false
		}

		if configuration != nil {
			_, oku := configuration.NamedSectionGet(*name, config.Remote, config.Username, "")
			_, okp := configuration.NamedSectionGet(*name, config.Remote, config.Password, "")

			if oku && okp {
				if !readOverride() {
					// Stop execution
					fmt.Println("Configuration not changed")
					return
				}
			}
		}

		realUsername := *username
		if realUsername == "" {
			realUsername = readLine("Username:", false)
		}

		realPassword := *password
		if realPassword == "" {
			realPassword = readLine("Password:", true)
		}

		encryptedPassword, salt := encrypt(realPassword, machineKey)

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedSectionSet(*name, config.Remote, config.Username, "", realUsername)
			mutator.NamedSectionSet(*name, config.Remote, config.Password, "", encryptedPassword)
			mutator.NamedSectionSet(*name, config.Remote, config.Salt, "", salt)
		})
	}
}
