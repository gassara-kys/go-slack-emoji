package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/h2non/filetype"
	"github.com/urfave/cli"
)

const (
	uploadURLFormat = "https://%s.slack.com/api/emoji.add"
	sleepCount      = 10
)

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "upload",
		Usage: "emoji upload command",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "upload-token",
				Usage:  "specific upload token(like xoxs-xxxx)",
				EnvVar: "SLACK_UPLOAD_TOKEN,UPLOAD_TOKEN",
			},
			cli.StringFlag{
				Name:   "upload-workspace, w",
				Usage:  "specific upload workspace",
				EnvVar: "SLACK_UPLOAD_WORKSPACE,UPLOAD_WORKSPACE,UPLOAD_WS",
			},
		},
		Action: func(c *cli.Context) error {
			return action(c, &upload{Out: os.Stdout})
		},
	})
}

type upload struct {
	Out         io.Writer
	uploadToken string
	workspace   string
}

func (u *upload) Run(c *cli.Context, token string) error {
	return u.Main(token, c.String("upload-token"), c.String("upload-workspace"))
}

func (u *upload) Main(token, uploadToken, uploadWorkspace string) error {
	u.uploadToken = uploadToken
	u.workspace = uploadWorkspace
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(dir + file.Name())
		if err != nil {
			return err
		}
		if !filetype.IsImage(buf) {
			continue
		}
		err = u.upload(dir + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *upload) upload(filename string) error {
	url := fmt.Sprintf(uploadURLFormat, u.workspace)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", strings.TrimSuffix(fi.Name(), path.Ext(fi.Name())))
	writer.WriteField("mode", "data")
	writer.WriteField("token", u.uploadToken)

	fw, err := writer.CreateFormFile("image", fi.Name())
	if err != nil {
		return err
	}
	fw.Write(fileContents)
	contentType := writer.FormDataContentType()
	if err = writer.Close(); err != nil {
		return err
	}

	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// onece retry post request after sleep a few secound.
		retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return err
		}
		fmt.Fprintf(u.Out, "retly alter %d secnods...\n", retryAfter)
		time.Sleep(time.Duration(retryAfter) * time.Second)
		resp, err := http.Post(url, contentType, body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected StatusCode, want=200, got=%d", resp.StatusCode)
	}

	fmt.Fprintf(u.Out, "uploaded file=%s\n", filename)
	return nil
}
