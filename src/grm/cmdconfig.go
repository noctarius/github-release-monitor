package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"grm/config"
	"fmt"
)

func cmdConfig(cmd *cli.Cmd) {
	cmd.Command("set", "Sets a configuration parameter", cmdConfigSet)
	cmd.Command("get", "Gets a configuration parameter", cmdConfigGet)
	cmd.Command("remove", "Removes a configuration parameter", cmdConfigRemove)
	cmd.Command("list", "Lists all configuration parameters", cmdConfigList)
}

func cmdConfigSet(cmd *cli.Cmd) {
	cmd.Spec = "NAME KEY VALUE [ --repository=<repository> ]"

	var (
		name       = cmd.StringArg("NAME", "", "The already defined remote user")
		key        = cmd.StringArg("KEY", "", "The property key to configure")
		value      = cmd.StringArg("VALUE", "", "The property's new value")
		repository = cmd.StringOpt("repository", "", "Set as repository specific override")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		if *key == "" {
			log.Fatal("No key specified")
		}

		realKey := config.KeyLookup(*key)
		if realKey == nil {
			log.Fatal(fmt.Sprintf("Unknown key specified: %s", *key))
		}

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedSectionSet(*name, config.Remote, realKey, *repository, *value)
		})
	}
}

func cmdConfigGet(cmd *cli.Cmd) {
	cmd.Spec = "NAME KEY [ --repository=<repository> ]"

	var (
		name       = cmd.StringArg("NAME", "", "The already defined remote user")
		key        = cmd.StringArg("KEY", "", "The property key to configure")
		repository = cmd.StringOpt("repository", "", "Set as repository specific override")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		if *key == "" {
			log.Fatal("No key specified")
		}

		realKey := config.KeyLookup(*key)
		if realKey == nil {
			log.Fatal(fmt.Sprintf("Unknown key specified: %s", *key))
		}

		if *repository != "" {
			if v, ok := configuration.NamedSectionGet(*name, config.Remote, realKey, *repository); ok {
				println(fmt.Sprintf("Configured value for key '%s' => %s", *key, v))
			}

		} else {
			if v, ok := configuration.NamedSectionGet(*name, config.Remote, realKey, ""); ok {
				println(fmt.Sprintf("Default value for key '%s' => %s", *key, v))
			}

			println("Existing overrides:")
			values := configuration.NamedSectionGetOverrides(*name, config.Remote, realKey)
			for k, v := range values {
				println(fmt.Sprintf("\t%s => %s", k, v))
			}
		}
	}
}

func cmdConfigRemove(cmd *cli.Cmd) {
	cmd.Spec = "NAME KEY [ --repository=<repository> ]"

	var (
		name       = cmd.StringArg("NAME", "", "The already defined remote user")
		key        = cmd.StringArg("KEY", "", "The property key to configure")
		repository = cmd.StringOpt("repository", "", "Set as repository specific override")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		if *key == "" {
			log.Fatal("No key specified")
		}

		realKey := config.KeyLookup(*key)
		if realKey == nil {
			log.Fatal(fmt.Sprintf("Unknown key specified: %s", *key))
		}

		configuration.ApplyChanges(func(mutator config.Mutator) {
			mutator.NamedSectionDelete(*name, config.Remote, realKey, *repository)
		})
	}
}

func cmdConfigList(cmd *cli.Cmd) {
	cmd.Spec = "NAME"

	var (
		name = cmd.StringArg("NAME", "", "The already defined remote user")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		println("Available configuration properties:")
		values := configuration.NamedSection(*name, config.Remote)
		for k, v := range values {
			println(fmt.Sprintf("%s => %s", k, v))
		}
	}
}
