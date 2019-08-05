# gswitch
## Overview
_gswitch_ is a simple tool to switch between [GCP](https://cloud.google.com/gcp/) projects. Easily.

_gswitch_ lets you:

1. Authenticate in GCP using your Google's User Account or Service Account.
1. Reconfigure _gcloud_ and _kubectl_ tools to a different GCP project without running complicated shell commands.


## Installation
### Requirements:
1. _gcloud_ is [installed](https://cloud.google.com/sdk/install) and available in PATH
1. Optional: _kubectl_ is [installed](https://kubernetes.io/docs/tasks/tools/install-kubectl/) and available in PATH


### Instructions for Linux and Windows:
1. [Download](https://github.com/RomanSharapov/gswitch/releases) the latest _gswitch_ release.
1. Unpack the downloaded archive.
1. Now put _gswitch_ binary to any folder from PATH and you're good to go!


### Install instructions for Mac:
1. Install a brew repository that contains formulae for _gswitch_:

        brew tap romansharapov/repo

1. Now you can install _gswitch_:

        brew install gswitch


## Examples
        Usage of gswitch:
            --no-auth               Do not authenticate user. Use previous identity
            --no-launch-browser     Do not launch a browser for authorization. If enabled or DISPLAY variable is not set, prints a URL to standard output to be copied. Disabled by default
        -p, --project string        Project ID for gcloud configuration. If project ID set, you can't use service account auth
            --use-service-account   Use service account instead of Google account. This option overwrites any other options!

### Scenario #1:
If you have a service account from a GCP project, just run _gswitch_ with `--use-service-account` parameter and it'll read `GOOGLE_APPLICATION_CREDENTIALS` environment variable, read account's JSON and switch your _gcloud_ tool to the project specified in the account file. If there're Kubernetes clusters available in the project, it'll switch _kubectl_ tool to the first discovered cluster in the project.

        gswitch --use-service-account

### Scenario #2:
If you want to authenticate with your Google Account and switch to some specific GCP project, just run `gswitch` without parameters and it'll obtain access credentials for your user account via a web-based authorization flow. If you'd rather authorize without a web browser but still interact with the command line, use the `--no-launch-browser` flag.

        gswitch --no-launch-browser

Additionally, you can specify `--project` parameter to switch to some specific project straight away. By default, it'll use the current project.

        gswitch --project <gcp_project_name>

If you just need to switch between project without reauthenticating yourself, use `--no-auth` parameter to use the previous identity.

        gswitch --project <gcp_project_name> --no-auth

If there're Kubernetes clusters available in the project, it'll switch _kubectl_ tool to the first discovered cluster in the project.
