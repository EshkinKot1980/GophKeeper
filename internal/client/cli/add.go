package cli

import (
	"encoding/json"
	"fmt"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add secret to system",
}

var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Add credentials to system",
	Run: func(cmd *cobra.Command, args []string) {
		name := prompt("Enter secret name: ")
		if name == "" {
			fmt.Println("name can not be empty")
			return
		}

		login := prompt("login: ")
		password := promptPassword("password: ")

		if login == "" || password == "" {
			fmt.Println("login and password can not be empty")
			return
		}

		data, err := json.Marshal(dto.Credentials{Login: login, Password: password})
		if err != nil {
			fmt.Printf("failed to decote credentials to json: %s\n", err.Error())
			return
		}

		err = secretServise.Upload(
			dto.SecretRequest{
				Name:     name,
				DataType: dto.SecretTypeCredentials,
				Meta:     []dto.MetaData{},
			},
			data,
		)
		if err != nil {
			fmt.Printf("failed to send data to server: %s\n", err.Error())
			return
		}
		// fmt.Println(string(data))
	},
}

func init() {
	addCmd.AddCommand(credentialsCmd)
	rootCmd.AddCommand(addCmd)
}
