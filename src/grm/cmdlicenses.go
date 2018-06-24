package main

import (
	"github.com/jawher/mow.cli"
	"fmt"
)

func cmdLicenses(cmd *cli.Cmd) {
	cmd.Spec = ""
	cmd.Action = func() {
		fmt.Println("Github-Release-Monitor uses (directly or indirectly) the following open source libraries:")

		fmt.Println(" * dateparse:")
		fmt.Println("\tLicense: MIT")
		fmt.Println("\tURL: https://github.com/adaddon/dateparse")

		fmt.Println(" * go-github:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: https://github.com/google/go-github")

		fmt.Println(" * go-querystring:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: https://github.com/google/go-querystring")

		fmt.Println(" * mow.cli:")
		fmt.Println("\tLicense: MIT")
		fmt.Println("\tURL: https://github.com/jawher/mow.cli")

		fmt.Println(" * go-isatty:")
		fmt.Println("\tLicense: MIT")
		fmt.Println("\tURL: https://github.com/mattn/go-isatty")

		fmt.Println(" * mpb:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: https://github.com/vbauerster/mpb")

		fmt.Println(" * ewma:")
		fmt.Println("\tLicense: MIT")
		fmt.Println("\tURL: https://github.com/vividcortex/ewma")

		fmt.Println(" * goini:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: https://github.com/zieckey/goini")

		fmt.Println(" * x/crypto/ssh:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: http://golang.org/x/crypto/ssh")

		fmt.Println(" * x/sys/unix:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: http://golang.org/x/sys/unix")

		fmt.Println(" * x/sys/windows:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: http://golang.org/x/sys/windows")

		fmt.Println(" * x/sys/windows:")
		fmt.Println("\tLicense: BSD 3-Clause")
		fmt.Println("\tURL: http://golang.org/x/sys/windows")
	}
}
