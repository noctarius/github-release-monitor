package main

import (
	"github.com/jawher/mow.cli"
	"log"
	"fmt"
	"github.com/zieckey/goini"
	"grm/config"
	"os"
)

func cmdExport(cmd *cli.Cmd) {
	cmd.Spec = "NAME [ --out=<outfile> ]"

	var (
		name = cmd.StringArg("NAME", "", "The name of the remote definition")
		out  = cmd.StringOpt("out", "", "The export path and filename, default: {NAME}.config")
	)

	cmd.Action = func() {
		if *name == "" {
			log.Fatal("No name specified")
		}

		outFile := *out
		if outFile == "" {
			outFile = fmt.Sprintf("%s.config", *name)
		}

		export := goini.New()
		values := configuration.NamedSection(*name, config.Remote)

		for k, v := range values {
			realKey := config.KeyLookup(k)
			if realKey.Exportable() {
				export.Set(k, v)
			}
		}

		file, err := os.Create(outFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not create export file '%s'", outFile), err)
		}

		export.Write(file)
		fmt.Println("Export successful")
	}
}
