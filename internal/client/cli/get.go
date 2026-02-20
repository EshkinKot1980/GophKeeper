package cli

import (
	"encoding/json"
	"fmt"
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
		id, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("ID must be a number")
		}

		secret, info, err := secretService.GetSecretAndInfo(id)
		if err != nil {
			return err
		}

		if info.DataType == dto.SecretTypeCredentials {
			return printCredentials(secret, info)
		}

		return fmt.Errorf("unsuported secret type: %s", info.DataType)
	},
}

func printCredentials(secret []byte, info dto.SecretResponse) error {
	var cr dto.Credentials
	err := json.Unmarshal(secret, &cr)
	if err != nil {
		return fmt.Errorf("failed decode secret json: %w", err)
	}

	fmt.Println(info.Name)
	fmt.Println("--------------------------------")
	fmt.Printf("login:    %s\n", cr.Login)
	fmt.Printf("password: %s\n", cr.Password)
	fmt.Println("--------------------------------")
	fmt.Println("created:", info.Created.Format("2006-01-02 15:04:05"))

	return nil
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&filePath, "out", "o", "", "Path to save file")
}
