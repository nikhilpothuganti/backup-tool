# Overview of the CLI Tool

## Description
This CLI tool is designed to back up files from a specified source directory to a backup directory. It supports encryption for added security and maintains a log of backup activities.

## Features

### Creating Backups

1. Running with no flag copies the source directory into a backup directory.

2. Once a directory has been copied, the program tracks which files were changed, so that only changed files are recopied in the backup directory next time. This feature is achieved by calculating file hash values of the files in the source directory and comparing with the backup directory.

3. There is a flag to encrypt the files while backing up. This is achieved by using AES encryption with a 32-byte encryption key.

4. A log file is created in the backup directory which logs all the backup activities.

5. The tool uses an efficient backing mechanism such that if a subdirectory has already been backed up, only the new directories/files are copied.

### Configuration Parameters
The tool uses a JSON configuration file (config.json) to store the backup root directory path. You can update this file using the -update-backup-dir flag.

## Functions
### loadConfig
```bash
loadConfig(configFile string) (*Config, error)
```
Description : Loads the configuration from the specified JSON file.

Parameters : 
- 'configFile' : Path to the configuration file.

Returns : 
- '*Config' : Pointer to a 'Config' struct containing the loaded configuration.
- 'error' :  An error if the configuration file could not be loaded.

### saveConfig
```bash
saveConfig(config *Config, configFile string) error
```

Description :  Saves the updated configuration to the specified JSON file.

Parameters : 
- 'config' : Pointer to a 'Config' struct containing the updated configuration.
- 'configFile' : Path to the configuration file.

Returns : 
- 'error' : An error if the configuration file could not be updated.

### encrypt
```bash
encrypt(data []byte, key []byte) ([]byte, error)
```

Description : Encrypts the given byte slice using AES encryption.

Parameters :
- 'data' : Data to be encrypted.
- 'key' : Encryption key (32 bytes for AES-256).

Returns :
- '[]byte' : Encrypted Data.
- 'error' : An error in case the encryption fails.

### calculateHash
```bash
calculateHash(filePath string) (string, error)
```

Description : Calculates the SHA-256 checksum for the file at the specified path.

Parameters :
- 'filePath' : Path to the file.

Returns :
- 'string' : SHA-256 checksum of the file.
- 'error' : An error if checksum calculation fails.







