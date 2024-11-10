# OMTI Tool

`omti` is a command-line tool designed to simplify GitHub repository management and database operations. It offers easy-to-use commands for tasks like creating repositories, pushing changes, pulling updates, and managing databases. This tool helps streamline workflows and boosts productivity.

## Installation

To install `omti`, use the following Go command:

```bash
go install github.com/LeTrongDat/omti@v1.0.1
```

## Usage

The basic structure for using `omti` commands is:

```bash
omti [command]
```

To view all available commands, use:

```bash
omti --help
```

## Commands Overview

| Command       | Description                                                                                              | Usage Example                                      |
|---------------|----------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| **completion** | Generates an autocompletion script for the specified shell.                                              | `omti completion`                                  |
| **db**         | Manages database-related operations, including backup and restore for PostgreSQL databases.              | `omti db [subcommand]`                             |
| **repo**       | Provides tools for managing GitHub repositories, such as creating repositories and tagging commits.      | `omti repo [subcommand]`                           |
| **help**       | Displays help information for any command.                                                               | `omti help`                                        |

### Database (`db`) Subcommands

| Subcommand | Description                                           | Usage Example                                                                                        |
|------------|-------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| **backup** | Backup a PostgreSQL database locally or over SSH.     | `omti db backup --remote <user>@<host>:<remote-db-port>`<br>Example: `omti db backup --remote admin@192.168.1.10:5432` |

#### Database Configuration Format

To specify database configuration, use this format:

```
<username>:<password>@<host>:<port>/<dbname>
```

**Example:**
```
postgres:v8hlDV0yMAHHlIurYupj@10.1.0.54:15432/golang
```

### Repository (`repo`) Subcommands

| Subcommand  | Description                                                              | Usage Example                                  |
|-------------|--------------------------------------------------------------------------|------------------------------------------------|
| **create**  | Create a new GitHub repository and push a local folder as the first commit. | `omti repo create`                             |
| **tag**     | Create a new tag for the latest commit of a branch in the repository.     | `omti repo tag`                                |

## Global Flags

| Flag              | Description                                        | Default Value |
|-------------------|----------------------------------------------------|---------------|
| `-h`, `--help`    | Display help for any command.                      |               |
| `--log-level`     | Set log level (`debug`, `info`, `warn`, `error`).  | `info`        |
