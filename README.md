# Autorclone

## Overview

Autorclone is a wrapper over [rclone](https://rclone.org/) attempting to automate repetitive tasks of synchronizing sources to destinations

## Features

- Synchronizes a (for now) local directory to one or more predefined [rclone `remotes`](https://rclone.org/docs/)
- TODO: runs in background as a daemon, executing defined jobs at specified intervals. Job definitions follow [crontab](https://crontab.guru/) notation

## Usage

`autorclone sync -h`

`autorclone daemon -h`

## License

[MIT](./LICENSE.md)
