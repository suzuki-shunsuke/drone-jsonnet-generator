# drone-jsonnet-generator

[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/drone-jsonnet-generator/status.svg)](https://cloud.drone.io/suzuki-shunsuke/drone-jsonnet-generator)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/drone-jsonnet-generator/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/drone-jsonnet-generator)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/drone-jsonnet-generator)](https://goreportcard.com/report/github.com/suzuki-shunsuke/drone-jsonnet-generator)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/drone-jsonnet-generator.svg)](https://github.com/suzuki-shunsuke/drone-jsonnet-generator)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/drone-jsonnet-generator.svg)](https://github.com/suzuki-shunsuke/drone-jsonnet-generator/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/drone-jsonnet-generator/master/LICENSE)

CLI tool to convert [Drone](https://drone.io/) v0.8x format .drone.yml to v1 format .drone.jsonnet

When you migrate matrix build to multi machine pipeline,
[jsonnet](https://jsonnet.org) is recommended officially.

https://docs.drone.io/user-guide/pipeline/migrating/

> The above syntax can be quite verbose if you are testing a large number of variations.
> To simplify your configuration we recommend using jsonnet.

But it is bothersome to convert .drone.yml to .drone.jsonnet manually .

`drone-jsonnet-generator` automates this bothersome tasks.

## Note

* This tool assumes that matrix build is used at source .drone.yml . If you don't use matrix build, you should use `drone convert` instead of this tool
* This tool can't convert perfectly. You should fix generated .drone.jsonnet . This tool automate 95% tasks. You should fix the following points.
  * Remove pipeline name's double quotes (ex. `"'GO_VERSION:' + GO_VERSION"` -> `'GO_VERSION:' + GO_VERSION`)
  * Change matrix build variable `${VARIABLE_NAME}` to jsonnet variable (ex. `"golang:${GO_VERSION}"` -> `"golang:" + GO_VERSION`)

## How to use

Run the command `drone-jsonnet-generator gen` and fix generated .drone.jsonnet manually.

## Example

### Example 1. matrix build without include

source .drone.yml

```yaml
---
pipeline:
  build:
    image: golang:${GO_VERSION}
    commands:
    - echo hello
services:
  database:
    image: ${DATABASE}
matrix:
  GO_VERSION:
  - 1.4
  - 1.3
  DATABASE:
  - mysql:5.5
  - mysql:6.5
```

generated .drone.jsonnet

```jsonnet
local pipeline(GO_VERSION, DATABASE) = {
  "kind": "pipeline",
  "name": "'GO_VERSION:' + GO_VERSION + ' DATABASE:' + DATABASE",
  "platform": {
    "arch": "amd64",
    "os": "linux"
  },
  "services": [
    {
      "image": "${DATABASE}",
      "name": "database",
      "pull": "default"
    }
  ],
  "steps": [
    {
      "commands": [
        "echo hello"
      ],
      "image": "golang:${GO_VERSION}",
      "name": "build",
      "pull": "default"
    }
  ]
};

local array_DATABASE = [
  "mysql:5.5",
  "mysql:6.5"
];
local array_GO_VERSION = [
  "1.4",
  "1.3"
];

[
  pipeline(GO_VERSION, DATABASE) for DATABASE in array_DATABASE for GO_VERSION in array_GO_VERSION 
]
```

### Example 2. matrix build with include

source .drone.yml

```yaml
---
pipeline:
  build:
    image: golang:${GO_VERSION}
    commands:
    - echo hello
services:
  database:
    image: ${DATABASE}
matrix:
  include:
  - GO_VERSION: 1.4
    DATABASE: mysql:5.5
  - GO_VERSION: 1.4
    DATABASE: mysql:6.5
  - GO_VERSION: 1.3
    DATABASE: mysql:5.5
```

generated .drone.jsonnet

```jsonnet
local pipeline(GO_VERSION, DATABASE) = {
  "kind": "pipeline",
  "name": "'GO_VERSION:' + GO_VERSION + ' DATABASE:' + DATABASE",
  "platform": {
    "arch": "amd64",
    "os": "linux"
  },
  "services": [
    {
      "image": "${DATABASE}",
      "name": "database",
      "pull": "default"
    }
  ],
  "steps": [
    {
      "commands": [
        "echo hello"
      ],
      "image": "golang:${GO_VERSION}",
      "name": "build",
      "pull": "default"
    }
  ]
};

local args = [
  {
    "DATABASE": "mysql:5.5",
    "GO_VERSION": "1.4"
  },
  {
    "DATABASE": "mysql:6.5",
    "GO_VERSION": "1.4"
  },
  {
    "DATABASE": "mysql:5.5",
    "GO_VERSION": "1.3"
  }
];

[
  pipeline(arg.GO_VERSION, arg.DATABASE) for arg in args
]
```

## Install

`drone-jsonnet-generator` is written with Golang and binary is distributed at [release page](https://github.com/suzuki-shunsuke/drone-jsonnet-generator/releases), so installation is easy and no dependency is needed.

## Usage

```console
$ drone-jsonnet-generator --help
NAME:
   drone-jsonnet-generator - convert Drone v0.8x format .drone.yml to v1 format .drone.jsonnet

USAGE:
   drone-jsonnet-generator [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   suzuki-shunsuke https://github.com/suzuki-shunsuke

COMMANDS:
     gen  convert Drone v0.8x format .drone.yml to v1 format .drone.jsonnet
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

```console
$ drone-jsonnet-generator gen --help
NAME:
   drone-jsonnet-generator gen - convert Drone v0.8x format .drone.yml to v1 format .drone.jsonnet

USAGE:
   drone-jsonnet-generator gen [command options] [arguments...]

OPTIONS:
   --source value, -s value  source .drone.yml path (default: ".drone.yml")
   --target value, -t value  target .drone.jsonnet path (default: ".drone.jsonnet")
   --stdout                  output generated jsonnet to stdout
```

## Trouble shooting

### yaml: unmarshal errors: cannot unmarshal !!map into string

https://discourse.drone.io/t/drone-cli-drone-yaml-bug-drone-convert-command-is-faied-to-parse-drone-yml/5005

If you use `include` or `exclude` in the step conditions about matrix builds, it is failed to generate the jsonnet file.

```yaml
---
pipeline:
  test:
    image: alpine:3.9.4
    commands:
    - echo ${NAME}
    when:
      matrix:
        include:  # failed to generate the jsonnet file
          NAME: main
matrix:
  NAME:
  - main
  - foo
```

```console
$ drone-jsonnet-generator gen
yaml: unmarshal errors:
  line 7: cannot unmarshal !!map into string
```

It is difficult to solve the problem, but there is a workaround.

1. remove the conditions about matrix builds from .drone.yml temporarily
2. generate the jsonnet file
3. add the step conditions about matrix builds to the jsonnet file manually

In case the above .drone.yml,

1. remove the conditions about matrix builds from .drone.yml temporarily

```yaml
---
pipeline:
  test:
    image: alpine:3.9.4
    commands:
    - echo ${NAME}
matrix:
  NAME:
  - main
  - foo
```

2. generate the jsonnet file

```console
$ drone-jsonnet-generator gen
```

generated jsonnet file

```jsonnet
local pipeline(NAME) = {
  "kind": "pipeline",
  "name": "'NAME:' + NAME",
  "platform": {
    "arch": "amd64",
    "os": "linux"
  },
  "steps": [
    {
      "commands": [
        "echo ${NAME}"
      ],
      "image": "alpine:3.9.4",
      "name": "test",
      "pull": "default"
    }
  ]
};

local array_NAME = [
  "main",
  "foo"
];

[
  pipeline(NAME) for NAME in array_NAME 
]
```

3. add the step conditions about matrix builds to the jsonnet file manually

Unfortunately, it is difficult to automate this step.

```jsonnet
  "steps": [
    (if NAME == "main" then
    {
      "commands": [
        "echo ${NAME}"
      ],
      "image": "alpine:3.9.4",
      "name": "test",
      "pull": "default"
    })
  ]
```

```jsonnet
local pipeline(NAME) = {
  "kind": "pipeline",
  "name": "'NAME:' + NAME",
  "platform": {
    "arch": "amd64",
    "os": "linux"
  },
  "steps": [
    (if NAME == "main" then
    {
      "commands": [
        "echo ${NAME}"
      ],
      "image": "alpine:3.9.4",
      "name": "test",
      "pull": "default"
    })
  ]
};

local array_NAME = [
  "main",
  "foo"
];

[
  pipeline(NAME) for NAME in array_NAME 
]
```

## Contribution

See [CONTRIBUTING.md](CONTRIBUTING.md) .

## License

[MIT](LICENSE)
