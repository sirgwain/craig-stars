package cmd

import (
	"fmt"

	"github.com/sirgwain/craig-stars/cs"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:                "create",
	Short:              "Create a resource",
	Long:               `Create a resource.`,
	PersistentPreRunE:  dbPreRun,
	PersistentPostRunE: dbPostRun,
}

func init() {
	rootCmd.AddCommand(createCmd)
	addCreateUserCommand()
}

func addCreateUserCommand() {

	var username string
	var email string
	var password string
	var role = "user"

	// createUserCmd creates a new user in the database
	createUserCmd := &cobra.Command{
		Use:   "user",
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := cs.NewUser(username, password, email, cs.UserRoleFromString(role))
			if err != nil {
				return err
			}

			err = client.CreateUser(user)
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
	createUserCmd.Flags().StringVar(&role, "role", "", "role for user")
	createUserCmd.MarkFlagRequired("username")
	createUserCmd.MarkFlagRequired("password")
	createUserCmd.MarkFlagRequired("email")

	createCmd.AddCommand(createUserCmd)
}
