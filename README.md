# Autorclone

> Note: this software is an alpha, work in progress.

## Overview

Autorclone is a wrapper over [rclone](https://rclone.org/) attempting to automate repetitive tasks of synchronizing sources to destinations

## Features

- Synchronizes a (for now) local directory to one or more predefined [rclone `remotes`](https://rclone.org/docs/)
- TODO: runs in background as a daemon, executing defined jobs at specified intervals. Job definitions follow [crontab](https://crontab.guru/) notation

## Usage

`autorclone -h`

```ini
Usage: autorclone.exe <command>

Flags:
  -h, --help                    Show context-sensitive help.
      --log-level=info          Set log level to one of: panic, fatal, error,
                                warn, info, debug, trace
      --rclone-path="rclone"    Path to rclone binary, if empty will use rclone
                                from PATH env
      --rclone-sync-args="-v --min-size 0.001 --multi-thread-streams 0 --retries 1 --human-readable --track-renames --log-format shortfile sync"
                                Rclone default sync arguments
                                ($AUTO_RCLONE_SYNC_ARGS)

Commands:
  upload <source-directory> <destination1 [destination2] [...]> ...
    Synchronize source to rclone remote destination(s). Use 'rclone config show'
    to list them.

  daemon
    Run as a background program, executing schelduled jobs

Run "autorclone.exe <command> --help" for more information on a command.
```

## License

[MIT](./LICENSE.md)
