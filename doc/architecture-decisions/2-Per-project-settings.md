# 2. Per-project settings

Date: 2020-04-02

## Status

Accepted

## Context

Each source project stands on its own, and the ADR is local to any of
the projects.

## Decision

Directory structure changed to:

```
$HOME
├── .adr/
|   ├── projects.json
/some/src/path/project-1
|    ├── .adr/
|    |   ├── config.json
|    |   ├── template.md
|    ├── basedir (configured in config.json)
|    |   ├── 1.Decision-documented.md
|    |   ├── 2.Decision-documented.md
/some/src/path/project-2
|    ├── .adr/
|    |   ├── config.json
|    |   ├── template.md
|    ├── other/basedir (configured in config.json)
|    |   ├── 1.Decision-documented.md
|    |   ├── 2.Decision-documented.md

```

where `projects.json` in this case would be
```json
{
 "projects": [
  "/some/src/path/project-1",
  "/some/src/path/project-2"
 ]
}
```

## Consequences

`adr new` command should also support finding out what project is run from, by
current working directory.

`adr new` also supports being run from any subdirectory within a project. 
