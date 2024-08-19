## fitanonymize

CLI util to process `.fit` files, stripping GPS coordinates and renaming them with a random 8-character alphanumeric name.

# WARNING 

Data is **not** completely deidentified. The tool was intended only to preserve
location privacy for a single or multiple users, and Does not currently remove
ANT+ device IDs or headunit information from `file_id` messages. A third-party
_will_ potentially be able to group data from a user if this is not removed.

Happy to accept a PR or add a flag for such if it's important for your use case.

## Features

- Find all `.fit` files in a specified directory or from a space-delimited list of files.
- Remove GPS data from `.fit` activity, session, and lap records.
- Save the processed files with a randomly generated name in specified output directory.

# Development

- `make build`: for current platform
- `make cross-compile`: build binaries for windows, darwin, and linux

# Usage

`./fitanonymize -files <files_or_directory> [-output <output_directory>]`

## Flags

- `-files`: (Required) A space-delimited set of .fit filenames or a directory
  containing .fit files.

- `-output`:  (Optional) The directory where the processed files will be saved.
  Defaults to the current directory.

## Examples

1. Process a directory of files:

    `./fitanonymize -files /path/to/files -output /path/to/output`

2. Process specific files:
  
    `./fitanonymize -files "file1.fit file2.fit file3.fit"`
