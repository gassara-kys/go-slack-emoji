package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/h2non/filetype"
	"github.com/slack-go/slack"
	"github.com/urfave/cli"
)

const dir = "./images/"

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "download",
		Usage: "emoji download command",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "output-file, f",
				Usage: "output download files",
			},
		},
		Action: func(c *cli.Context) error {
			return action(c, &download{Out: os.Stdout})
		},
	})
}

type download struct {
	Out io.Writer
}

func (d *download) Run(c *cli.Context, token string) error {
	return d.Main(token, c.Bool("output-file"))
}

func (d *download) Main(token string, outputFile bool) error {
	api := slack.New(token)
	emoji, err := api.GetEmoji()
	if err != nil {
		return err
	}
	r := regexp.MustCompile(`^alias*`)
	for name, url := range emoji {
		fmt.Fprintf(d.Out, "%s: %s\n", name, url)
		if r.MatchString(url) || !outputFile {
			continue
		}
		if err := d.downloadFile(name, url); err != nil {
			return err
		}
	}
	return nil
}

func (d *download) downloadFile(name, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("not expected status code: want=200, got=%d", response.StatusCode)
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if !filetype.IsImage(buf) {
		return fmt.Errorf("not image file")
	}
	kind, _ := filetype.Image(buf)
	filename := dir + name + "." + kind.Extension
	if err := ioutil.WriteFile(filename, buf, 0664); err != nil {
		return err
	}
	return nil
}
