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

type Cluster struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup clusters for oc-multicluster-tui",
	Long: `The setup command allows you to configure clusters for oc-multicluster-tui.
You can append new clusters, reset the clusters file, delete specific clusters, or list all configured clusters.
For more information, use the --help flag.`,

	Run: func(cmd *cobra.Command, args []string) {

		configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
		filePath := configDir + "/clusters.json"

		if cmd.Flag("help").Changed {
			helpSetup()
			return
		}

		if cmd.Flag("list").Changed {
			listClusters()
			return
		} else if _, err := os.Stat(filePath); err == nil && cmd.Flag("reset").Changed {
			resetClustersFile()
		} else if _, err := os.Stat(filePath); err == nil && cmd.Flag("append").Changed {
			fmt.Println("Appending to existing clusters file...")
			appendToClustersFile()
			return
		} else if _, err := os.Stat(filePath); err == nil && cmd.Flag("delete-cluster").Changed {
			clusterName, _ := cmd.Flags().GetString("delete-cluster")
			deleteClusterFromFile(clusterName)
			return
		} else if _, err := os.Stat(filePath); err == nil && (!cmd.Flag("append").Changed || !cmd.Flag("reset").Changed) {
			fmt.Println("Config file already exists, skipping setup.")
			fmt.Println("Use --append to add new clusters or --reset to reset the clusters file.")
			return
		} else {
			fmt.Println("Creating new clusters file...")
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
		saveClustersToFile(clusters)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	setupCmd.Flags().BoolP("append", "a", true, "Append to existing clusters file")
	setupCmd.Flags().BoolP("reset", "r", true, "Reset the clusters file to default state")
	setupCmd.Flags().StringP("delete-cluster", "d", "", "Specify a cluster name to delete from the clusters file")
	setupCmd.Flags().BoolP("list", "l", true, "List all clusters in the clusters file")
	setupCmd.Flags().BoolP("help", "h", false, "Help message for setup command")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func resetClustersFile() {
	// Reset the clusters file to default state
	configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
	filePath := configDir + "/clusters.json"
	if err := os.Remove(filePath); err != nil {
		fmt.Println("Error resetting clusters file:", err)
		return
	}
	fmt.Printf("Clusters file reset to default state at %s\n", filePath)
}
func appendToClustersFile() {
	// Append a new cluster to the existing clusters file
	configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
	filePath := configDir + "/clusters.json"

	// Load existing clusters
	var clusters []Cluster
	file, err := os.Open(filePath)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		_ = decoder.Decode(&clusters)
	}

	// Interactive input loop
	for {
		var name, url, cont string
		fmt.Print("Enter cluster name: ")
		fmt.Scanln(&name)
		//make sure cluster name is unique
		for _, cluster := range clusters {
			if cluster.Name == name {
				fmt.Println("Cluster name already exists. Please enter a unique name.")
				fmt.Print("Enter cluster name: ")
				fmt.Scanln(&name)
			}
		}
		fmt.Print("Enter cluster API URL: ")
		fmt.Scanln(&url)
		// make sure URL is unique
		for _, cluster := range clusters {
			if cluster.URL == url {
				fmt.Println("Cluster URL already exists. Please enter a unique URL.")
				fmt.Print("Enter cluster API URL: ")
				fmt.Scanln(&url)
			}
		}
		clusters = append(clusters, Cluster{Name: name, URL: url})

		fmt.Print("Add another cluster? (y/n): ")
		fmt.Scanln(&cont)
		if cont != "y" && cont != "Y" {
			break
		}
	}

	// Save updated clusters list
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening clusters file for writing:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(clusters); err != nil {
		fmt.Println("Error saving clusters:", err)
		return
	}
	fmt.Printf("Clusters appended and saved to %s\n", filePath)
}
func saveClustersToFile(clusters []Cluster) {
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
}
func deleteClusterFromFile(clusterName string) {
	// Delete a cluster from the clusters file
	configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
	filePath := configDir + "/clusters.json"

	// Load existing clusters
	var clusters []Cluster
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening clusters file:", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&clusters); err != nil {
		fmt.Println("Error decoding clusters:", err)
		return
	}

	// Filter out the cluster to delete
	var updatedClusters []Cluster
	for _, cluster := range clusters {
		if cluster.Name != clusterName {
			updatedClusters = append(updatedClusters, cluster)
		}
	}

	// Save updated clusters list
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error opening clusters file for writing:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(updatedClusters); err != nil {
		fmt.Println("Error saving clusters:", err)
		return
	}
	fmt.Printf("Cluster '%s' deleted and changes saved to %s\n", clusterName, filePath)
}

func listClusters() {
	// List all clusters in the clusters file
	configDir := os.Getenv("HOME") + "/.config/oc-multicluster-tui"
	filePath := configDir + "/clusters.json"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening clusters file:", err)
		return
	}
	defer file.Close()

	var clusters []Cluster
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&clusters); err != nil {
		fmt.Println("Error decoding clusters:", err)
		return
	}

	if len(clusters) == 0 {
		fmt.Println("No clusters found.")
		return
	}

	fmt.Println("Clusters:")
	for _, cluster := range clusters {
		fmt.Printf("- %s: %s\n", cluster.Name, cluster.URL)
	}
}
func helpSetup() {
	// Display help message for setup command
	fmt.Println("Usage: oc-multicluster-tui setup [flags]")
	fmt.Println("Flags:")
	fmt.Println("  -a, --append            Append to existing clusters file")
	fmt.Println("  -r, --reset             Reset the clusters file to default state")
	fmt.Println("  -d, --delete-cluster    Specify a cluster name to delete from the clusters file")
	fmt.Println("  -l, --list              List all clusters in the clusters file")
	fmt.Println("  -h, --help              Help message for setup command")
}
