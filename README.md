# Catan Map Generator (CMG)

Little tool for generating Catan maps, hopefully more fair than manual.

## Badges

[![GitHub release](https://img.shields.io/github/release/joostvdg/cmg.svg)]()
[![license](https://img.shields.io/github/license/joostvdg/cmg.svg)]()
[![Go Report Card](https://goreportcard.com/badge/joostvdg/cmg)](https://goreportcard.com/report/github.com/joostvdg/cmg)

## Jenkins X

```
jx import
```

## Google Cloud Function

See the wrapper project [github.com/joostvdg/cmg-gcf](https://github.com/joostvdg/cmg-gcf) for running it as a Cloud Function.

## Heroku

* get Heroku CLI: `brew tap heroku/brew && brew install heroku` 
    * linux `sudo snap install --classic heroku`
    * [download installer](https://cli-assets.heroku.com/heroku-x64.exe) for windows
* create app with Heroku: `heroku create`
* install govendor `go get -u github.com/kardianos/govendor`
* add `vendor` directory to `.gitignore` file
* init govendor `govendor init`
    * need to make sure you're in a directory explicitly in the `GOPATH`, symlinks don't work...
* fetch local files for govendor:  `govendor fetch +local`
* publish to Heroku app: `git push heroku master`