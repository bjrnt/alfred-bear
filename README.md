# alfred-bear

> Allows you to quickly search for and open Bear notes from Alfred. It will search tags, note content, titles, etc. to try to bring you the best match.

![usage example screenshot](/image.png)

## Installation

1. Grab the latest release [here](https://github.com/bjrnt/alfred-bear/releases/) and install the workflow file.
2. Set `BEAR_TOKEN` inside the workflow's environment variables (select the workflow and press the `[x]` icon in the upper right). The token can be found in by going to `Help > API Token > Copy Token` in Bear.

## Usage

Open Alfred and type `b` and try typing a query. You can also configure a hotkey for it inside the workflow settings.

## Maintenance

### Issues and Feature Requests

Feel free to open an issue for the project if you have encountered a problem or have a feature request for the workflow.

### Building

The project can be built, linked to Alfred, and released using [jason0x34/go-alfred](https://github.com/jason0x43/go-alfred). Commands for this can be found in [./.vscode/tasks.json](./.vscode/tasks.json).
