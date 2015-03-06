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
	"sort"
)

// FLAGS
var urlFlag = cli.StringFlag { Name: "url, u", Usage: "Nexus' URL." }
var repoFlag = cli.StringFlag { Name: "repository, r", Usage: "The repository to inspect." }

// ACTIONS
func countArtifactsAction(c *cli.Context) {
	if(len(os.Args) == 0) {
		cli.ShowSubcommandHelp(c)
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
}

func listArtifactsAction(c *cli.Context) {
	if(len(os.Args) == 0) {
		cli.ShowSubcommandHelp(c)
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
	for _, artifact := range artifacts {
		fmt.Printf("%v\n", artifact)
	}
}

func repoAsSimpleString(r *nexus.Repository) string {
	if r == nil {
		return ""
	}

	return r.Name + " (" + r.ID + "): " + r.Type
}

type repoSort []*nexus.Repository

func (rs repoSort) Len() int { return len(rs) }
func (rs repoSort) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

// hosted < virtual < proxy
func (rs repoSort) Less(i, j int) bool { 
	a, b := rs[i], rs[j]

	if a.RemoteURI != "" && b.RemoteURI == "" {
		return false
	}

	if a.RemoteURI == "" && b.RemoteURI != "" {
		return true
	}

	if a.Type == "hosted" && b.Type != "hosted" {
		return true
	}

	if a.Type != "hosted" && b.Type == "hosted" {
		return false
	}

	return a.Name < b.Name
}

func listReposAction(c *cli.Context) {
	if(len(os.Args) == 0) {
		cli.ShowAppHelp(c)
		return
	}

	url := c.String("url")
	
	if url == "" {
		fmt.Println("[ERROR] No URL found: use the flag --url")
		return
	}

	n := nexus.New(url, credentials.None)

	repos, err := n.Repositories()
	if err != nil {
		fmt.Printf("[ERROR] %v: %v\n", reflect.TypeOf(err), err)
		return
	}

	fmt.Printf("%d repositories in %v:\n", len(repos), url)
	sort.Sort(repoSort(repos))
	for _, repo := range(repos) {
		fmt.Printf("\t%v\n", repoAsSimpleString(repo))
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "nexus-cli"
	app.Usage = "Runs some queries against a Sonatype Nexus server."
	app.Version = "0.1.0"
	app.Author = "Humberto Anjos"
	app.Email = "h.anjos@acm.org"
	app.Commands = []cli.Command {
		{
			Name: "artifacts",
			ShortName: "a",
			Usage: "Commands concerning the artifacts and GAVs in a given repository.",
			Flags: []cli.Flag{ urlFlag, repoFlag },
			Subcommands: []cli.Command {
				{
					Name: "count",
					ShortName: "c",
					Usage: "Counts the artifacts and GAVs in a given repository.",
					Flags: []cli.Flag{ urlFlag, repoFlag },
					Action: countArtifactsAction,
				},
				{
					Name: "list",
					ShortName: "l",
					Usage: "Lists the artifacts and GAVs in a given repository.",
					Flags: []cli.Flag{ urlFlag, repoFlag },
					Action: listArtifactsAction,
				},
			},
			Action: func (c *cli.Context) {
				cli.ShowCommandHelp(c, c.Command.Name)
			},
		},
		{
			Name: "repositories",
			ShortName: "repos",
			Usage: "Commands concerning repositories.",
			Flags: []cli.Flag{ urlFlag },
			Subcommands: []cli.Command {
				{
					Name: "list",
					ShortName: "l",
					Usage: "Lists the repositories.",
					Flags: []cli.Flag{ urlFlag },
					Action: listReposAction,
				},
			},
			Action: func (c *cli.Context) {
				cli.ShowSubcommandHelp(c)
			},
		},
	}
	app.Action = func (c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}