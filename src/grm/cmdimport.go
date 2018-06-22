package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"grm/config"
	"fmt"
	"strings"
	"github.com/zieckey/goini"
)

func cmdImport(cmd *cli.Cmd) {
	cmd.Spec = "NAME IMPORTFILE"

	var (
		name       = cmd.StringArg("NAME", "", "The name of the remote definition")
		importFile = cmd.StringArg("IMPORTFILE", "", "The path and filename of the config to import")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		if *importFile == "" {
			log.Fatal("No import file specified")
		}

		readOverride := func() bool {
			line := readLine(fmt.Sprintf("A configuration for %s already exists, properties might get "+
				"overridden. Do you really want to continue? [yes|No]", *name), false)

			line = strings.ToLower(line)

			if line == "yes" || line == "y" || line == "true" {
				return true
			}
			return false
		}

		if t := configuration.NamedSection(*name, config.Remote); len(t) > 0 {
			if !readOverride() {
				// Stop execution
				println("Configuration not changed")
				return
			}
		}

		importer := goini.New()
		if err := importer.ParseFile(*importFile); err != nil {
			log.Fatal("Error opening the import file", err)
		}

		if values, ok := importer.GetKvmap(goini.DefaultSection); ok {
			configuration.ApplyChanges(func(mutator config.Mutator) {
				for k, v := range values {
					realKey := config.KeyLookup(k)
					specifier := config.ExtractSpecifier(k)
					if realKey.Exportable() {
						mutator.NamedSectionSet(*name, config.Remote, realKey, specifier, v)
					}
				}
			})

			println("Import successful")
		} else {
			println("Import failed, nothing to import")
		}
	}
}
