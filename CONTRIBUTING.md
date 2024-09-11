# Contributing Guide

## Pre-requisites

- Git
- NodeJS 20+
- Yarn
- go 1.22.1
- Mage
- Docker

## How to run this locally

- Clone the repo locally `git clone https://github.com/grafana/infinity-libs infinity-libs` and `cd infinity-libs`
- Install packages `yarn`
- Run the backend tests with `yarn test`

## Releasing

To create a new version of a package for release, follow these steps:

1. Check out the commit you want to tag by running: `git checkout <COMMIT_SHA>`.
   - **Note:** Ensure that this commit includes the updated version in the `package.json` of the package you wish to release.   
2. Tag the commit with: `git tag lib/go/<PACKAGE_NAME>/<VERSION>` (e.g., **lib/go/jsonframer/v1.1.1**).
   - **Note:** We are using lightweight tags, so no additional options are necessary.
3. Push the tag to the remote repository with: `git push origin lib/go/<PACKAGE_NAME>/<VERSION>`.
4. Verify that the tag was created successfully [here](https://github.com/grafana/infinity-libs/tags).
