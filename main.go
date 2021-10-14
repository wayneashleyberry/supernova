package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

const prefix = "SUPERNOVA"

type specification struct {
	GithubAccessToken string `envconfig:"GITHUB_ACCESS_TOKEN" required:"true"`
	GithubUsername    string `envconfig:"GITHUB_USERNAME" required:"true"`
}

func getSpecification() specification {
	var s specification

	envconfig.MustProcess(prefix, &s)

	return s
}

func main() {
	rootCmd := &cobra.Command{
		Use: "supernova",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "read",
		Short: "Print a list of your GitHub stars",
		RunE: func(cmd *cobra.Command, args []string) error {
			return readStars()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "env",
		Short: "Print the required environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			var s specification

			return envconfig.Usage(prefix, &s)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Unstar everything on GitHub",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteStars()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func deleteStars() error {
	s := getSpecification()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.GithubAccessToken},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	page := 0
	empty := false

	for !empty {
		sr, _, err := client.Activity.ListStarred(ctx, s.GithubUsername, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return fmt.Errorf("list starred: %w", err)
		}

		if len(sr) == 0 {
			empty = true
		}

		for _, r := range sr {
			url := r.GetRepository().GetHTMLURL()
			url = strings.Replace(url, "https://github.com/", "", 1)

			parts := strings.Split(url, "/")

			_, err := client.Activity.Unstar(ctx, parts[0], parts[1])
			if err != nil {
				return fmt.Errorf("unstar: %w", err)
			}
		}

		page++
	}

	return nil
}

func readStars() error {
	s := getSpecification()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: s.GithubAccessToken},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	page := 0
	empty := false

	starredRepositories := []*github.StarredRepository{}

	for !empty {
		sr, _, err := client.Activity.ListStarred(ctx, s.GithubUsername, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return fmt.Errorf("list starred: %w", err)
		}

		starredRepositories = append(starredRepositories, sr...)

		if len(sr) == 0 {
			empty = true
		}

		page++
	}

	for _, starredRepository := range starredRepositories {
		fmt.Println(starredRepository.GetRepository().GetHTMLURL())
	}

	return nil
}
