# CLI BACKUP TOOL

A command-line tool for backing up files with optional encryption.

## Installation
1. Clone the repository :

```bash
git clone https://github.com/sprixter/version-control-system.git
cd version-control-system
```

2. Make it a go module :

```bash
go init module version-control-system
```

3. Build the executable :

```bash
go build -o vcs
```

4. Run the CLI tool :

```bash
./vcs -source=<source_directory>
```

## Usage

```bash
Usage of ./vcs:
  -config string
    	Path to the configuration file. (default "config.json")
  -encrypt
    	Enable encryption.
  -source string
    	Source directory to backup. (default ".")
  -update-backup-dir
    	Update backup directory.
```

### Encrypt files
To enable encryption while backing up files :

```bash
  ./vcs -encrypt=<true/false>
```

### Change Backup Directory

```bash
  ./vcs -update-backup-dir=true
```

You then receive a prompt to enter the new backup directory :

```bash
  Enter the new backup directory :
  /your/new/backup/directory
```

## Configuration
The tool reads its configuration from a JSON file (config.json by default). The configuration file should have the following format:

```json
{
  "backupRoot": "/your/backup/directory"
}
```
Edit the config.json file to suitable backup directory before running the code.
