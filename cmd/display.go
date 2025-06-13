/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// displayCmd represents the display command
var displayCmd = &cobra.Command{
	Use:   "display",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("display called")
		// Display all clusters with styled output
		displayAllClustersStyled()

	},
}

func init() {
	rootCmd.AddCommand(displayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// displayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// displayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getClusterInfoStyled(clusterName, kubeconfigPath string) string {
	// Styles
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("231"))

	// Version
	versionCmd := exec.Command("sh", "-c", fmt.Sprintf("oc version --short --kubeconfig=%s", kubeconfigPath))
	versionOut, _ := versionCmd.CombinedOutput()
	version := strings.TrimSpace(string(versionOut))

	// Nodes
	nodesCmd := exec.Command("sh", "-c", fmt.Sprintf("oc get nodes --no-headers --kubeconfig=%s", kubeconfigPath))
	nodesOut, _ := nodesCmd.CombinedOutput()
	nodes := strings.Split(strings.TrimSpace(string(nodesOut)), "\n")
	readyCount := 0
	for _, node := range nodes {
		if strings.Contains(node, " Ready ") {
			readyCount++
		}
	}
	nodeStatus := fmt.Sprintf("%d/%d ready", readyCount, len(nodes))

	// Pods
	podsCmd := exec.Command("sh", "-c", fmt.Sprintf("oc get pods --all-namespaces --no-headers --kubeconfig=%s", kubeconfigPath))
	podsOut, _ := podsCmd.CombinedOutput()
	pods := strings.Split(strings.TrimSpace(string(podsOut)), "\n")
	notReady := 0
	for _, pod := range pods {
		if !strings.Contains(pod, " Running ") && !strings.Contains(pod, " Completed ") {
			notReady++
		}
	}
	podStatus := fmt.Sprintf("%d not ready / %d total", notReady, len(pods))

	// Resource usage
	topCmd := exec.Command("sh", "-c", fmt.Sprintf("oc top nodes --no-headers --kubeconfig=%s", kubeconfigPath))
	topOut, _ := topCmd.CombinedOutput()
	resourceUsage := strings.TrimSpace(string(topOut))

	// Events
	eventsCmd := exec.Command("sh", "-c", fmt.Sprintf("oc get events --all-namespaces --sort-by=.lastTimestamp --kubeconfig=%s | tail -n 5", kubeconfigPath))
	eventsOut, _ := eventsCmd.CombinedOutput()
	events := strings.TrimSpace(string(eventsOut))

	// Compose styled output
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("Cluster: %s\n", clusterName)))
	b.WriteString(labelStyle.Render("Version: ") + valueStyle.Render(version) + "\n")
	b.WriteString(labelStyle.Render("Nodes: ") + valueStyle.Render(nodeStatus) + "\n")
	b.WriteString(labelStyle.Render("Pods: ") + valueStyle.Render(podStatus) + "\n")
	b.WriteString(labelStyle.Render("Resource Usage:\n") + valueStyle.Render(resourceUsage) + "\n")
	b.WriteString(labelStyle.Render("Last 5 Events:\n") + valueStyle.Render(events) + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("-", 40)) + "\n")
	return b.String()
}

func displayAllClustersStyled() {
	kubeDir := os.Getenv("HOME") + "/.kube/oc-multicluster-tui"
	files, _ := filepath.Glob(kubeDir + "/*.kubeconfig")
	for _, file := range files {
		clusterName := strings.TrimSuffix(filepath.Base(file), ".kubeconfig")
		fmt.Print(getClusterInfoStyled(clusterName, file))
	}
}
