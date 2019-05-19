workflow "GitHub Release" {
  on = "push"
  resolves = ["Release Mac", "Release Windows", "Release Linux"]
}

action "Build" {
  uses = "joostvdg/actions/golang-build@master"
  args = "linux/amd64 darwin/amd64 windows/amd64"
  env = {
    GO111MODULE = "on"
  }
}

action "Release Mac" {
  uses = "ngs/go-release.action@v1.0.1"
  needs = ["Build"]
  secrets = ["GITHUB_TOKEN"]
  env = {
    GOOS = "darwin"
    GOARCH = "amd64"
    GO111MODULE = "on"
  }
}

action "Release Windows" {
  uses = "ngs/go-release.action@v1.0.1"
  needs = ["Build"]
  secrets = ["GITHUB_TOKEN"]
  env = {
    GOOS = "windows"
    GOARCH = "amd64"
    GO111MODULE = "on"
  }
}

action "Release Linux" {
  uses = "ngs/go-release.action@v1.0.1"
  needs = ["Build"]
  secrets = ["GITHUB_TOKEN"]
  env = {
    GOOS = "linux"
    GOARCH = "amd64"
  }
}
