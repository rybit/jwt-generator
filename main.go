package main

import (
	"log"

	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

var secret string

type JWTClaims struct {
	jwt.StandardClaims
	ID     string   `json:"id"`
	Email  string   `json:"email"`
	Groups []string `json:"groups"`
}

func main() {
	root := &cobra.Command{}
	root.PersistentFlags().StringVarP(&secret, "secret", "s", "", "secret value to use")

	gen := &cobra.Command{
		Use: "gen",
		Run: generate,
	}
	gen.Flags().StringSliceP("group", "g", []string{}, "groups to add to the token")
	gen.Flags().StringP("user", "u", "", "the user to have for the token")
	gen.Flags().StringP("email", "e", "", "the email to use for the token")
	root.AddCommand(gen)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func validateArgs(cmd *cobra.Command) {
	if secret == "" {
		log.Fatal("Must provide a secret")
	}
}

func generate(cmd *cobra.Command, args []string) {
	validateArgs(cmd)

	groups, err := cmd.Flags().GetStringSlice("group")
	if err != nil {
		log.Fatal("Failed to get groups")
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		log.Fatal("Failed to get email")
	}

	user, err := cmd.Flags().GetString("user")
	if err != nil || user == "" {
		log.Fatal("Must provide a user")
	}

	claims := &JWTClaims{
		ID:     user,
		Email:  email,
		Groups: groups,
	}

	result, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
