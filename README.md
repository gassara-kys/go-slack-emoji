# go-slack-emoji

slack emoji tools

## get USER_TOEKN( like `xoxs-XXXXXX`)
- go to your slack.com > https://`<your-workspace>`.slack.com/customize/emoji
- you can get your USER_TOEKN.
  - open browser console and exec `boot_data.api_token`.

# set your token
- set the token
```bash
# create env.sh
$ cp env_sample.sh env.sh

# edit SLACK_TOKEN
$ vi  env.sh
```

## download
- specific `SLACK_TOKEN` in env.sh file.
- exec `make download` command
- download emoji file in `images/` directory
```bash
# download 
$ make download
```

## upload
- specific `SLACK_UPLOAD_TOKEN` and `SLACK_UPLOAD_WORKSPACE` in env.sh file.
- exec `make upload` command
- upload `images/` files to SLACK_UPLOAD_WORKSPACE.
```bash
# download 
$ make upload
```
