package cli

import (
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user in the GophKeeper system",
	RunE:  register,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login the GophKeeper system",
	RunE:  login,
}

func login(cmd *cobra.Command, args []string) error {
	cr, err := prompt.Credentials()
	if err != nil {
		return err
	}
	return authService.Login(cr)
}

func register(cmd *cobra.Command, args []string) error {
	cr, err := prompt.RegisterCredentials()
	if err != nil {
		return err
	}
	return authService.Register(cr)
}

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
}
