# Instant Sonar

**Instantly analyse your code with SonarQube in Docker, with a single command.**

[![Build](https://github.com/LucaScorpion/instant-sonar/actions/workflows/build.yml/badge.svg)](https://github.com/LucaScorpion/instant-sonar/actions/workflows/build.yml)
[![Publish](https://github.com/LucaScorpion/instant-sonar/actions/workflows/publish.yml/badge.svg)](https://github.com/LucaScorpion/instant-sonar/actions/workflows/publish.yml)

[![asciicast of instant-sonar](https://asciinema.org/a/604152.svg)](https://asciinema.org/a/604152)

## Usage

Simply download the appropriate binary for your platform from the [latest release](https://github.com/LucaScorpion/instant-sonar/releases/latest) and execute it.

To analyse your current working directory:

```shell
instant-sonar
```

To analyse a different directory:

```shell
instant-sonar "path/to/project"
```

For all help and options:

```shell
instant-sonar --help
```

## What it Does

Instant Sonar will start a [SonarQube](https://hub.docker.com/_/sonarqube) container,
configure it so it can be accessed without having to log in,
and set up a new project.
It will then start a [Sonar Scanner](https://hub.docker.com/r/sonarsource/sonar-scanner-cli) container,
which will run the project analysis and send the results to SonarQube.
After this is done, it will output a link which you can use to view the analysis results.

Is something not working as intended?
Feel free to [create an issue](https://github.com/LucaScorpion/instant-sonar/issues/new)!
