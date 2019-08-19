# Catan Map Generator (CMG)

Little tool for generating Catan maps which are fair.

## Badges

[![GitHub release](https://img.shields.io/github/release/joostvdg/cmg.svg)]()
[![license](https://img.shields.io/github/license/joostvdg/cmg.svg)]()
[![Go Report Card](https://goreportcard.com/badge/joostvdg/cmg)](https://goreportcard.com/report/github.com/joostvdg/cmg)
[![DepShield Badge](https://depshield.sonatype.org/badges/joostvdg/cmg/depshield.svg)](https://depshield.github.io)
[![Coverage Status](https://coveralls.io/repos/github/joostvdg/cmg/badge.svg?branch=master)](https://coveralls.io/github/joostvdg/cmg?branch=master)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-78%25-brightgreen.svg?longCache=true&style=flat)</a>

## Live version

The application is running as a free heroku app to try: [https://catan-map-generator.herokuapp.com/api/map](https://catan-map-generator.herokuapp.com/api/map).

Look for the Code to see the API for for more commands (or the shell scripts such as `6p-game.sh`).

The UI is available at [https://catanmap.herokuapp.com/](https://catanmap.herokuapp.com/).

## Run on Jenkins X

```
jx import
```

## Run as Google Cloud Function

See the wrapper project [github.com/joostvdg/cmg-gcf](https://github.com/joostvdg/cmg-gcf) for running it as a Cloud Function.

## Run on Heroku

* get Heroku CLI: `brew tap heroku/brew && brew install heroku` 
    * linux `sudo snap install --classic heroku`
    * [download installer](https://cli-assets.heroku.com/heroku-x64.exe) for windows
* create app with Heroku: `heroku create`
* configure go.mod file (see below)
* publish to Heroku app: `git push heroku master`

### Go.mod

When using Go 1.11+, just use `gomod` for managing deps, don't use the others tools unless you have a good reason.

When using `gomod`, you have to add some annotations for Heroku to your `go.mod` file.
Docs state it defaults to building `.`, it seemed it was `./cmd/.` which was not correct for me.

Probably best to state it explicitly.

```bash
// +heroku goVersion go1.11
// +heroku install .
```

For more info, [read Heroku's docs on Go](https://elements.heroku.com/buildpacks/heroku/heroku-buildpack-go).
