package cmd

import (
	"fmt"

	"github.com/sirgwain/craig-stars/config"
	"github.com/sirgwain/craig-stars/cs"
	"github.com/sirgwain/craig-stars/db"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource",
	Long:  `Create a resource.`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	addCreateUserCommand()
}

func addCreateUserCommand() {

	var username string
	var email string
	var password string
	var role cs.Role = cs.RoleUser

	// createUserCmd represents the createUser command
	createUserCmd := &cobra.Command{
		Use:   "user",
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			user := cs.NewUser(username, password, email, role)
			db := db.NewClient()
			cfg := config.GetConfig()
			db.Connect(cfg)

			err := db.CreateUser(user)
			if err != nil {
				return err
			}

			fmt.Printf("user %s (%d) created\n", user.Username, user.ID)
			return nil
		},
	}

	createUserCmd.Flags().StringVarP(&username, "username", "u", "", "username to create")
	createUserCmd.Flags().StringVarP(&password, "password", "p", "", "password for user")
	createUserCmd.Flags().StringVarP(&email, "email", "e", "", "email for user")
	createUserCmd.Flags().Var(&role, "role", "role for user")
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("password")
	createUserCmd.MarkFlagRequired("email")

	createCmd.AddCommand(createUserCmd)
}
