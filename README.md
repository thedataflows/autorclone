# Autorclone

> Note: this software is an alpha, work in progress.

## Overview

Autorclone is a wrapper over [rclone](https://rclone.org/) attempting to automate repetitive tasks of synchronizing sources to destinations. You still need to define rclone remotes with rclone utility.

Requires rclone to be installed and present in the path or can be manually specified.

## Features

- Synchronizes a rclone source (or local directory/file) to a destination that may be another directory or rclone predefined remote (`rclone listremotes`). Synchronization means that destination files and directories **will be deleted (if not present on the source)**. By default, a backup will be done by rclone before deletion/modification.
- TODO: runs in background as a daemon, executing defined jobs at specified intervals. Job definitions follow [crontab](https://crontab.guru/) notation

## Usage

`autorclone -h`

```ini
Usage: autorclone.exe <command>

Flags:
  -h, --help                    Show context-sensitive help.
      --log-level=info          Set log level to one of: panic, fatal, error,
                                warn, info, debug, trace
      --rclone-path="rclone"    Path to rclone binary, by default will try
                                rclone from PATH env
      --rclone-version="v1.57.0"
                                Rclone release to be downloaded if not in PATH
      --rclone-sync-args="sync -v --min-size 0.001 --multi-thread-streams 0 --retries 1 --human-readable --track-renames --links --ignore-errors --log-format shortfile"
                                Rclone default sync arguments
                                ($AUTORCLONE_SYNC_ARGS)

Commands:
  sync <source> <destination1 [destination2] [...]> ...
    Synchronize source to rclone destination(s). Use 'rclone config show' to
    list them.

  run
    Manually run predefined sync jobs. Without any argument, will run all jobs
    in the predefined job definition file

  daemon
    Run as a background program, executing schelduled jobs

  version
    Show version and exit

Run "autorclone.exe <command> --help" for more information on a command.
```

## Tips

1. Cleanup backups on destination
    - Run initial sync: `autorclone sync source destination` (may result in some .rclonebak files)
    - Review the backup files and keep what might be overwritten or deleted by mistake
    - Run sync again with no backup suffix: `autorclone sync --backup-suffix "" source destination` (since destination has .rclonebak files, they will be deleted)
2. Set `--backup-suffix ""` when the destination is a cloud storage

## License

[MIT](./LICENSE.md)
