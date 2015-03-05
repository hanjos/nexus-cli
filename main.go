package main

import (
	"github.com/hanjos/nexus"
	"github.com/hanjos/nexus/credentials"
	"github.com/hanjos/nexus/search"
	
	count "github.com/hanjos/nexus-cli/artifact-count"

	"github.com/codegangsta/cli"

	"fmt"
	"reflect"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "nexus-cli"
	app.Usage = "Runs some queries against a Sonatype Nexus server."
	app.Version = "0.1.0"
	app.Author = "Humberto Anjos"
	app.Email = "h.anjos@acm.org"
	app.Commands = []cli.Command {
		{
			Name: "artifact-count",
			ShortName: "ac",
			Usage: "Counts the artifacts and GAVs in a given repository.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "url, u",
					Usage: "Nexus' URL.",
				},
				cli.StringFlag{
					Name: "repository, r",
					Usage: "The repository to inspect.",
				},
			},
			Action: func(c *cli.Context) {
				if(len(os.Args) == 0) {
					cli.ShowAppHelp(c)
					return
				}

				url := c.String("url")
				repo := c.String("repository")
				hasErrors := false
				
				if repo == "" {
					hasErrors = true
					fmt.Println("[ERROR] No repository given!")
				}

				if url == "" {
					hasErrors = true
					fmt.Println("[ERROR] No URL given!")
				}

				if hasErrors {
					return
				}

				fmt.Printf("%v @ %v\n", repo, url)
				n := nexus.New(url, credentials.None)

				artifacts, err := n.Artifacts(search.ByRepository(repo))
				if err != nil {
					fmt.Printf("[ERROR] %v: %v\n", reflect.TypeOf(err), err)
					return
				}

				fmt.Printf("\t%d artifacts\n", len(artifacts))
				gavs := count.GavsOf(artifacts)

				fmt.Printf("\t%d GAVs\n", len(gavs))
			},
		},
	}
	app.Action = func (c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}