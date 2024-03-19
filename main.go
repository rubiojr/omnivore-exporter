package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/go-shiori/obelisk"
	"github.com/rubiojr/omnivore-go"
	"github.com/urfave/cli/v2"
)

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

func exportURL(url string, filePath string, debug, compress bool) error {
	req := obelisk.Request{
		URL: url,
	}

	arc := obelisk.Archiver{EnableLog: debug}
	arc.Validate()

	result, _, err := arc.Archive(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to archive: %v", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %v", err)
	}
	defer f.Close()

	if compress {
		gz := gzip.NewWriter(f)
		defer gz.Close()
		_, err = gz.Write(result)
	} else {
		_, err = f.Write(result)
	}
	return err
}

func exportMonolith(url string, filePath string, debug, compress bool) error {
	cmdName := "monolith"

	_, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("command '%s' not found in PATH: %w", cmdName, err)
	}

	var cmdArgs []string
	if debug {
		cmdArgs = append(cmdArgs, "--silent")
	}

	cmdArgs = append(cmdArgs, "--output", filepath.Clean(filePath), url)

	cmd := exec.Command(cmdName, cmdArgs...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "Omnivore Export",
		Usage: "A tool to export Omnivore documents.",
		Commands: []*cli.Command{
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export all documents to HTML files.",
				Action:  exportAll,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "debug",
					},
					&cli.BoolFlag{
						Name: "no-color",
					},
					&cli.BoolFlag{
						Name: "compress",
					},
					&cli.BoolFlag{
						Name:     "use-monolith",
						Aliases:  []string{"m"},
						Usage:    "Use monolith to archive if available in PATH.",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "output-dir",
						Aliases:  []string{"o"},
						Usage:    "Output directory for the exported files.",
						Required: false,
						Value:    "omnivore-exports",
					},
					&cli.StringSliceFlag{
						Name:     "labels",
						Aliases:  []string{"l"},
						Usage:    "Export only articles labeled with label.",
						Required: false,
					},
					&cli.StringSliceFlag{
						Name:     "skip-labels",
						Usage:    "Export only articles labeled with label.",
						Required: false,
						Value:    cli.NewStringSlice("omnivore-exporter-skip", "Newsletter"),
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: "+err.Error())
	}

}

func exportAll(c *cli.Context) error {
	debug := c.Bool("debug")
	if c.Bool("no-color") {
		color.NoColor = true
	}

	compress := c.Bool("compress")
	useMonolith := c.Bool("use-monolith")
	outputDir := c.String("output-dir")
	exportLabels := c.StringSlice("labels")

	labels := ""
	if len(exportLabels) > 0 {
		fmt.Println("Exporting only articles with labels:", exportLabels)
		labels = labelsToQuery(exportLabels)
	}

	skipLabels := c.StringSlice("skip-labels")
	if len(labels) == 0 && len(skipLabels) > 0 {
		fmt.Println("Exporting all articles except those with labels:", skipLabels)
		labels = skipLabelsQuery(skipLabels)
	}

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create base dir: %v", err)
	}
	token := getAPIToken()
	if token == "" {
		return fmt.Errorf("OMNIVORE_API_TOKEN is not set")
	}

	if compress && !useMonolith {
		fmt.Println("Compressing enabled")
	}

	if useMonolith {
		fmt.Println("Using: monolith")
	} else {
		fmt.Println("Using: obelisk")
	}

	query := fmt.Sprintf("in:all %s sort:saved", labels)
	fmt.Println("Search query:", query)
	fmt.Println("Exporting to folder", outputDir, "...")
	client := omnivore.NewClient(omnivore.Opts{Token: token})
	a, err := client.Search(
		omnivore.SearchOpts{
			Query: query,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to search: %v", err)
	}

	exportFn := exportURL
	if useMonolith {
		exportFn = exportMonolith
	}

	counter := 0
	for _, searchItem := range a {
		title := searchItem.Title
		filePath := filepath.Join(outputDir, title+".html")
		if compress {
			filePath += ".gz"
		}
		if fileExists(filePath) {
			skip("Skipping %s %s (exists)", searchItem.ID, title)
			continue
		}
		info("Exporting '%s'...", title)
		url := searchItem.Url
		err = exportFn(url, filePath, debug, compress)
		if err != nil {
			fail("Failed to export '%s' (ignoring)", title)
			continue
		}
		counter++
	}
	fmt.Printf("\n%d documents %s \n", counter, green("exported."))

	return nil
}

func getAPIToken() string {
	return os.Getenv("OMNIVORE_API_TOKEN")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func labelsToSlice(labels []omnivore.Label) []string {
	l := []string{}
	for _, label := range labels {
		l = append(l, strings.ToLower(label.Name))
	}
	return l
}

func labelsToQuery(labels []string) string {
	l := []string{}
	for _, label := range labels {
		l = append(l, fmt.Sprintf("label:%s", label))
	}
	return strings.Join(l, " OR ")
}

func skipLabelsQuery(labels []string) string {
	l := []string{}
	for _, label := range labels {
		l = append(l, fmt.Sprintf("-label:%s", label))
	}
	return strings.Join(l, " AND ")
}

func debug(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

func info(msg string, args ...interface{}) {
	fmt.Printf("* "+msg+"\n", args...)
}

func skip(msg string, args ...interface{}) {
	fmt.Printf(yellow("* ")+msg+"\n", args...)
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, red("* ")+msg+"\n", args...)
}
