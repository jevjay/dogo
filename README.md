# dogo
Dogo is a Slack chatbot, which can be extended with custom Docker images.

**Status:** Pre-Alpha.

The aim of this bot is to provide a user with an easy way to extend its functionality with a simple "plug'n'play" model.

By providing instructions via configuration file you can "teach" Dogo to perform "virtually" any task and greatly boost your team efficiency.

## Usage

To build locally:

```
git clone https://github.com/jevjay/dogo.git
cd dogo
go get -v -u github.com/nlopes/slack
go get -v -u gopkg.in/yaml.v2
go get -v -u github.com/docker/docker/api
go get -v -u github.com/docker/docker/client
go build -o dogo .
```

To build container locally:

```
git clone https://github.com/jevjay/dogo.git
cd dogo
docker build -t dogo .
```

## Contributing

When contributing to this repository, please first discuss the change you wish to make via issue, email, or any other method with the owners of this repository before making a change.

Please note we have a code of conduct, please follow it in all your interactions with the project
