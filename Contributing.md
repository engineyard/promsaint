# Contribution Guidelines

The typical git workflow applies for contributions. Create your work in a new branch
and create a PR against master when ready. A maintainer has to review and accept the 
PR before it can be merged.

## Building

`make`

The binaries are created in the `bin/` folder.

## Releasing

The release artifact is a `tar.gz` archive with the following structure:
```
promsaint/
|- promsaint
|- promsaint-cli
`- promsaint-smoke
```
The name pattern of the archive is `promsaint-<ts>-<arch>.tar.gz`
* `<ts>` timestamp id in the format `YYYYMMDD<ID>`. `<ID>` is a three digit wide running id (starting at `001`)
* `<arch>` is the architecture of the build (e.g. `amd64`)

The artifact is distributed via `distfiles.engineyard.com`.

**TODO:** this can and should be improved.
