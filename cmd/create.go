package cmd

import (
	"fmt"

	"github.com/sirgwain/craig-stars/appcontext"
	"github.com/sirgwain/craig-stars/game"

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
	var password string
	var role game.Role = game.RoleUser

	// createUserCmd represents the createUser command
	createUserCmd := &cobra.Command{
		Use:   "user",
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := appcontext.Initialize()
			user := game.NewUser(username, password, role)
			err := ctx.DB.SaveUser(user)
			if err != nil {
				return err
			}

			fmt.Printf("user %s (%d) created\n", user.Username, user.ID)
			return nil
		},
	}

	createUserCmd.Flags().StringVarP(&username, "username", "u", "", "username to create")
	createUserCmd.Flags().StringVarP(&password, "password", "p", "", "password for user")
	createUserCmd.Flags().Var(&role, "role", "role for user")
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("password")

	createCmd.AddCommand(createUserCmd)
}
