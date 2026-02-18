package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run: func(cmd *cobra.Command, args []string) {
		cr := dto.Credentials{}

		cr.Login = prompt("login: ")
		if err := cr.ValidateLogin(); err != nil {
			fmt.Println(err.Error())
			return
		}

		cr.Password = promptPassword("password: ")
		if err := cr.ValidatePassword(); err != nil {
			fmt.Println(err.Error())
			return
		}

		err := authService.Register(cr)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the GophKeeper system",
	Run: func(cmd *cobra.Command, args []string) {
		login := prompt("login: ")
		password := promptPassword("password: ")

		if login == "" || password == "" {
			fmt.Println("login and password can not be empty")
			return
		}

		cr := dto.Credentials{Login: login, Password: password}
		err := authService.Login(cr)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
}

func prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptPassword(label string) string {
	fmt.Print(label)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	fmt.Println()
	return string(bytePassword)
}
