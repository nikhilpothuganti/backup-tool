package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	BackupRoot string `json:"backupRoot"`
}

func loadConfig(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// saveConfig saves the updated configuration to the file
func saveConfig(config *Config, configFile string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// encrypts a given byte slice using AES encryption
func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// calculates sha256 checksum for a file
func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func main() {
	configFile := "config.json"
	sourceDir := flag.String("source", ".", "Source directory to backup")
	updateConfig := flag.Bool("update-backup-dir", false, "Update configuration parameters.")
	encryptFlag := flag.Bool("encrypt", false, "Enable encryption.")
	flag.Parse()
	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	if *updateConfig {
		// Prompt for new backup root directory
		fmt.Println("Enter new backup directory:")
		fmt.Scanln(&config.BackupRoot)

		err = saveConfig(config, configFile)
		if err != nil {
			fmt.Println("Error saving config:", err)
			os.Exit(1)
		}
		fmt.Println("Configuration updated successfully!")
	}

	// Create/Append to log file
	LOG_FILE := filepath.Join(config.BackupRoot, "backup.log")
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Encryption key (32 bytes for AES-256)
	key := []byte("01234567890123456789012345678901")

	// Walk through the source directory and copy files
	err = filepath.Walk(*sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(*sourceDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(config.BackupRoot, relPath)

		if info.IsDir() {
			err = os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				if *encryptFlag {
					srcFile, err := os.Open(path)
					if err != nil {
						return err
					}
					defer srcFile.Close()
					data, err := io.ReadAll(srcFile)
					if err != nil {
						return err
					}
					// encrypt file
					encryptedData, err := encrypt(data, key)
					if err != nil {
						return err
					}
					destFile, err := os.Create(destPath)
					if err != nil {
						return err
					}
					defer destFile.Close()
					_, err = destFile.Write(encryptedData)
					if err != nil {
						return err
					}
					log.Printf("Encrypted and copied %s to %s\n", path, destPath)

				} else {
					// Copy file without encryption
					srcFile, err := os.Open(path)
					if err != nil {
						return err
					}
					defer srcFile.Close()
					destFile, err := os.Create(destPath)
					if err != nil {
						return err
					}
					defer destFile.Close()
					_, err = io.Copy(destFile, srcFile)
					if err != nil {
						return err
					}
					log.Printf("Copied %s to %s\n", path, destPath)
				}

			} else {
				srcChecksum, err := calculateHash(path)
				if err != nil {
					fmt.Println(err)
					return err
				}

				dstChecksum, err := calculateHash(destPath)
				if err != nil {
					fmt.Println(err)
					return err
				}
				if srcChecksum != dstChecksum {
					if *encryptFlag {
						srcFile, err := os.Open(path)
						if err != nil {
							return err
						}
						defer srcFile.Close()
						data, err := io.ReadAll(srcFile)
						if err != nil {
							return err
						}
						// encrypt file
						encryptedData, err := encrypt(data, key)
						if err != nil {
							return err
						}
						destFile, err := os.Create(destPath)
						if err != nil {
							return err
						}
						defer destFile.Close()
						_, err = destFile.Write(encryptedData)
						if err != nil {
							return err
						}
						log.Printf("Encrypted and copied %s to %s\n", path, destPath)

					} else {
						// Copy file without encryption
						srcFile, err := os.Open(path)
						if err != nil {
							return err
						}
						defer srcFile.Close()
						destFile, err := os.Create(destPath)
						if err != nil {
							return err
						}
						defer destFile.Close()
						_, err = io.Copy(destFile, srcFile)
						if err != nil {
							return err
						}
						log.Printf("Copied %s to %s\n", path, destPath)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through source directory")
		return
	}

	fmt.Printf("Backup completed.")
}
