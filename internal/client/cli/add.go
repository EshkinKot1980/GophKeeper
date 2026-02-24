package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/spf13/cobra"
)

const (
	// Название записи метаданных для имени файла
	MetaFileName = "FileName"
	// Название записи метаданных для директории файла (абсолютный путь)
	MetaFilePath = "FilePath"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a secret to the system",
}

var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Adds credentials to the system",
	RunE:  addCredentials,
}

var fileCmd = &cobra.Command{
	Use:   "file <file_path>",
	Short: "Adds a file to the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		return addFile(os.Stdout, args[0])
	},
}

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Adds text data to the system",
	RunE:  addText,
}

func addCredentials(cmd *cobra.Command, args []string) error {
	name, err := prompt.SecretName()
	if err != nil {
		return err
	}

	credentials, err := prompt.Credentials()
	if err != nil {
		return err
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to encode credentials to json: %w", err)
	}

	err = secretService.Upload(
		dto.SecretRequest{
			Name:     name,
			DataType: dto.SecretTypeCredentials,
			Meta:     []dto.MetaData{},
		},
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to send data to server: %w", err)
	}
	return nil
}

func addFile(out io.Writer, path string) error {
	// Проверяем доступность и размер файла
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileName := info.Name()
	if info.IsDir() {
		return fmt.Errorf("failed add file: the file \"%s\" is directory", info.Name())
	}

	size := info.Size()
	if size > cfg.FileMaxSize {
		return fmt.Errorf("file size (%d bytes) exceeds the limit of %d bytes", size, cfg.FileMaxSize)
	}

	// Формируем метаданные
	meta := []dto.MetaData{}
	meta = append(meta, dto.MetaData{Name: MetaFileName, Value: fileName})

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute file path: %w", err)
	}
	meta = append(meta, dto.MetaData{Name: MetaFilePath, Value: filepath.Dir(absPath)})

	name, err := prompt.SecretName()
	if err != nil {
		return err
	}

	// Читаем файл и отправляем в зашифрованном виде на сервер
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Fprintln(out, "sending the file to the server")
	err = secretService.Upload(
		dto.SecretRequest{
			Name:     name,
			DataType: dto.SecretTypeFile,
			Meta:     meta,
		},
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to send data to server: %w", err)
	}
	return nil
}

func addText(cmd *cobra.Command, args []string) error {
	name, err := prompt.SecretName()
	if err != nil {
		return err
	}

	text, err := prompt.Text()
	if err != nil {
		return err
	}

	err = secretService.Upload(
		dto.SecretRequest{
			Name:     name,
			DataType: dto.SecretTypeText,
			Meta:     []dto.MetaData{},
		},
		[]byte(text),
	)
	if err != nil {
		return fmt.Errorf("failed to send data to server: %w", err)
	}

	return nil
}

func init() {
	addCmd.AddCommand(credentialsCmd)
	addCmd.AddCommand(fileCmd)
	addCmd.AddCommand(textCmd)
	rootCmd.AddCommand(addCmd)
}
