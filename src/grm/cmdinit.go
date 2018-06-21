package main

import (
	"github.com/jawher/mow.cli"
	"strings"
	"github.com/zieckey/goini"
)

func cmdinit(cmd *cli.Cmd) {
	cmd.Spec = "[ -u=<username> ] [ -p=<password> ]"

	var (
		username = cmd.StringOpt("u username", "", "Specify the username to access Github")
		password = cmd.StringOpt("p password", "", "Specify the password to access Github")
	)

	cmd.Action = func() {
		readOverride := func() bool {
			line := readLine("You already have a configuration file, are you sure to override? [yes/No] ", false)
			line = strings.ToLower(line)

			if line == "yes" || line == "y" || line == "true" {
				return true
			}
			return false
		}

		if config != nil {
			_, oku := config.SectionGet(sectionCredentials, keyUsername)
			_, okp := config.SectionGet(sectionCredentials, keyPassword)

			if oku && okp {
				if !readOverride() {
					// Stop execution
					println("Configuration file not overridden, stopping")
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

		changeConfig(func(config *goini.INI) {
			config.SectionSet(sectionCredentials, keyUsername, realUsername)
			config.SectionSet(sectionCredentials, keyPassword, encryptedPassword)
			config.SectionSet(sectionCredentials, keySalt, salt)
		})
	}
}
