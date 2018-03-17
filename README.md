# magicbot
bot for posting magic cards to slack

## Installation
Magicbot vendors its dependencies (with [dep]), so it shouldn't require any external commands to install. Simply checkout the repo and run:

```shell
$ go install
```

Magicbot comes with a Docker image as well, which can be built with:

```shell
$ docker build -t magicbot .
```

In the future these images may be published to Docker Hub.

[dep]: https://github.com/golang/dep
