/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to OpenShift clusters using CLI tokens",
	Long: `The login command allows users to authenticate with OpenShift clusters using CLI tokens.
This command is essential for accessing cluster resources and managing multiple clusters from a single interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Log in in to OpenShift clusters.

		fmt.Println("login called")

		//read clusters file and get the list of clusters
		ClustersFile, err := readClustersFile()
		if err != nil {
			fmt.Printf("Error reading clusters file: %v\n", err)
			return
		}

		fmt.Println("Available clusters:")
		for i, cluster := range ClustersFile {
			fmt.Printf("%d: %s (%s)\n", i+1, cluster.Name, cluster.URL)
		}
		//log in to each cluster and save each kubeconfig to $HOME/.kube/oc-multicluster-tui/<cluster_name>.kubeconfig
		for _, cluster := range ClustersFile {
			fmt.Printf("Logging in to cluster %s at %s...\n", cluster.Name, cluster.URL)
			fmt.Println("Please enter the CLI token for cluster ", cluster.Name, ": ")
			var token string
			fmt.Scanln(&token)
			loginCmdStr := fmt.Sprintf("oc login %s --token=%s --kubeconfig=$HOME/.kube/oc-multicluster-tui/%s.kubeconfig", cluster.URL, token, cluster.Name)
			out, err := exec.Command("sh", "-c", loginCmdStr).CombinedOutput()
			if err != nil {
				fmt.Printf("Failed to log in to cluster %s: %v\nOutput: %s\n", cluster.Name, err, string(out))
				continue
			}
			fmt.Printf("Successfully logged in to cluster %s\n", cluster.Name)
			fmt.Printf("kubeconfig saved to $HOME/.kube/oc-multicluster-tui/%s.kubeconfig\n", cluster.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func readClustersFile() ([]Cluster, error) {
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
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Clusters file not found at %s. Run oc-multicluster-tui setup to create Clusters file.", filePath)
			return nil, fmt.Errorf("clusters file not found: %w", err)
		}
		fmt.Printf("Error reading clusters file: %v\n", err)
		return nil, err
	}
	// If no clusters are found, return an empty slice
	if len(clusters) == 0 {
		fmt.Println("No clusters found in the clusters file. Please run oc-multicluster-tui setup to add clusters.")
		return clusters, nil
	}
	// Return the list of clusters
	fmt.Printf("Clusters loaded from %s\n", filePath)
	return []Cluster{}, nil
}
