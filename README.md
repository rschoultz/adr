[![Build Status](https://travis-ci.com/rschoultz/adr.svg?branch=master)](https://travis-ci.com/rschoultz/adr)

# ADR Go
A minimalist command line tool written in Go to work with 
[Architecture Decision Records](http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions) 
(ADRs).

Greatly inspired by the [adr-tools](https://github.com/npryce/adr-tools) 
with all of the added benefits of using the Go instead of Bash.

# Quick start
## Installing adr
Go to the [releases page](https://github.com/marouni/adr/releases) 
and grab one of the binaries that corresponds to your platform.

Alternatively, if you have a Go developement environment setup you can install it directly using :
```bash
go get github.com/rschoultz/adr && go install github.com/rschoultz/adr
```

## Initializing adr
Before creating any new ADR you need to initiatilize the source directory that  
that will host your ADRs and use the `init` sub-command to initialize the configuration1 :

```bash
cd some-project-dir
adr init docs/architecture/decisions
```

This will create a subdirectory `docs/architecture/decisions` to hold your decisions,
and also an `.adr` directory in the project with the files `config.json` 
and `template.md`, that the `adr` command will reference.  

## Creating a new ADR

As simple as :
```bash
adr new my awesome proposition
```
This will create a new numbered ADR in the initialized folder:
`docs/architecture/decisions/xxx-my-new-awesome-proposition.md`.

Next, just open the file in your preferred markdown editor and starting writing your ADR.
If you have configured environment variables VISUAL or EDITOR, that editor will be called
automatically.

## Table of contents

Generate a markdown table of contents by

```bash
adr generate toc
```

and using flag --prefix (or -p) you can prefix the links.


### Add both .adr and the decisions to source control

```bash
git add .adr 
git add docs/architecture/decisions
git commit -m "Added Architecture Decisions"
```

If you forgot to add commit the .adr directory to source control and lost
it, you can do "init" and the directory will re-appear. 
You should update the counter in the config.json file. 
Existing decisions will not be affected. 
