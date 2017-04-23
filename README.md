# process-slack-status

Change your Slack status when running different processes/programs.

## getting started...

This doesn't integrate as an official API user so is slightly more annoying to setup.

The easiest way to grab the data needed is to open the team chat with the web
client and open a JS console. Then grab the following,

```js
TS.model.api_token
// this is for --api-token

TS.model.api_url
// this needs the site name prepending, then pass to --api-url
// usually just /api/ so you need to pass https://whatever-team.slack.com/api/

TS.boot_data.version_uid
// this is for --version-uid
```

You'll need to create a config file, see `config.toml.example`. You just need to
add a [table](https://github.com/toml-lang/toml#table) for each program named
after the executable (lower-case), then provide the emoji and text to use. There
isn't any fancy ordering of which will come first so it is best to only setup
for programs you wouldn't normally run at the same time...

The `*` key is used to set a default status that will be used if none of the
matching programs are running, if it isn't set then your status won't
automatically "reset" once the program stops.

Then you should be able to run it like so,

```sh
$ go get hawx.me/code/process-slack-status
$ process-slack-status \
     --api-token 'xoxs-something-something' \
     --api-url 'https://my-team.slack.com/api/' \
     --version-uid 'aebcaaeaabebebabcba' \
     --config my-config.toml
```
