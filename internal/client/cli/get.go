package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

var filePath string

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get secret from system by ID",
	Long:  "Downloads secret from the server and save it to disk for file type.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return get(os.Stdout, args[0])
	},
}

func get(out io.Writer, argID string) error {
	id, err := strconv.ParseUint(argID, 10, 64)
	if err != nil {
		return fmt.Errorf("id must be a number")
	}

	secret, info, err := secretService.GetSecretAndInfo(id)
	if err != nil {
		return err
	}

	switch info.DataType {
	case dto.SecretTypeCredentials:
		return outputCredentials(out, secret, info)
	case dto.SecretTypeFile:
		return saveFile(secret, info)
	case dto.SecretTypeText:
		return outputText(out, secret, info)
	}

	return fmt.Errorf("unsuported secret type: %s", info.DataType)
}

func outputCredentials(out io.Writer, secret []byte, info dto.SecretInfo) error {
	var cr dto.Credentials
	err := json.Unmarshal(secret, &cr)
	if err != nil {
		return fmt.Errorf("failed decode secret json: %w", err)
	}

	fmt.Fprintln(out, info.Name)
	fmt.Fprintln(out, "--------------------------------")
	fmt.Fprintf(out, "login:    %s\n", cr.Login)
	fmt.Fprintf(out, "password: %s\n", cr.Password)
	fmt.Fprintln(out, "--------------------------------")
	fmt.Fprintln(out, "created:", info.Created.Format("2006-01-02 15:04:05"))

	return nil
}

func saveFile(secret []byte, info dto.SecretInfo) error {
	// если не задан путь для вывода файла, товыводим его в текущую директорию
	if filePath == "" {
		cd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("faled to spot current directory: %w", err)
		}
		filePath = cd
	}

	pathInfo, err := os.Stat(filePath)

	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to get output file info: %w", err)
		}
		// если файл не существует, проверяем существует ли директория его содержащая
		fileDir := filepath.Dir(filePath)

		if _, err := os.Stat(fileDir); err != nil {
			return fmt.Errorf("failed to spot output directory")
		}
		//директория существует - оставляем путь к файлу как есть
	} else {
		// файл существует
		if pathInfo.IsDir() {
			// если это директория добавляем к пути имя из меты
			fileName := getFileName(info.Meta)
			if fileName == "" {
				return fmt.Errorf("failed to spot file name from metadata")
			}
			filePath = filepath.Join(filePath, fileName)
			// проверяем существут ли итоговый файл,
			// и если да запрашиваем у пользователя нужно ли его переписать
			_, err := os.Stat(filePath)
			if err == nil && !prompt.Overwrite(filePath) {
				return nil
			}
		} else {
			// если файл не директория, запрашиваем у пользователя нужно ли его переписать
			if !prompt.Overwrite(filePath) {
				return nil
			}
		}
	}

	err = os.WriteFile(filePath, secret, 0600)
	if err != nil {
		return fmt.Errorf("failed to to save file: %w", err)
	}

	return nil
}

func outputText(out io.Writer, secret []byte, info dto.SecretInfo) error {
	fmt.Fprintln(out, info.Name)
	fmt.Fprintln(out, "--------------------------------")
	fmt.Fprintln(out, string(secret))
	fmt.Fprintln(out, "--------------------------------")
	fmt.Fprintln(out, "created:", info.Created.Format("2006-01-02 15:04:05"))

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&filePath, "out", "o", "", "Path to save file")
}
