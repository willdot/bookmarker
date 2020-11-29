# Bookmarker
Bookmarker is a CLI app powered by Firesearch that allows you to store url bookmarks with tags.
The bookmarks can then be searched for by using the tags provided when the bookmark was created.

For example, a bookmark for `https://github.com` can be stored with the tags `git` `source control` etc. Then searching for the bookmark using `source` or `git` will return the bookmark with the URL which can be opened.

## Usage

### Installation
`go get github.com/willdot/bookmarker`

### Firesearch
Currently this CLI tool has only been tested with Firesearch in dev. To set up it up in dev, follow the instructions [here](https://firesearch.dev/docs/tutorial) to get the Firebase emulator started and the Firesearch Docker container up and running.

To use against production Firesearch, follow the instructions [here](https://firesearch.dev/docs/deploy) and ensure that the following environment variables are set before running:
`ENDPOINT`
`INDEXPATH`
`SECRET`

TODO: Try this out using production Firesearch

### Add a bookmark
`bookmarker add --url=github.com --name=github --tags="git, source control"`
Note: The tags must be wrapped in `""` and comma separated.

### Find a bookmark
`bookmarker find git`

## Shoutouts
Firesearch is the power behind this project. Check it out at https://firesearch.dev It's so easy to get up and running with excellent API documentation.
