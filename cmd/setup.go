/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		type Cluster struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		var clusters []Cluster
		for {
			var name, url, cont string
			fmt.Print("Enter cluster name: ")
			fmt.Scanln(&name)
			fmt.Print("Enter cluster API URL: ")
			fmt.Scanln(&url)
			clusters = append(clusters, Cluster{Name: name, URL: url})

			fmt.Print("Add another cluster? (y/n): ")
			fmt.Scanln(&cont)
			if cont != "y" && cont != "Y" {
				break
			}
		}
		// Save clusters to a file in JSON format
		configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			return
		}
		filePath := configDir + "/clusters.json"
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(clusters); err != nil {
			fmt.Println("Error saving clusters:", err)
			return
		}
		fmt.Printf("Clusters saved to %s\n", filePath)

		fmt.Println("setup called")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
