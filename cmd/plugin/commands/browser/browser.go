package browser

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/ingress-nginx/cmd/plugin/commands/browser/fzf"
	"k8s.io/ingress-nginx/cmd/plugin/commands/ingresses"
	"k8s.io/ingress-nginx/cmd/plugin/request/k8sclient"
	"k8s.io/ingress-nginx/cmd/plugin/util"
)

// CreateCommand creates and returns this cobra subcommand
func CreateCommand(flags *genericclioptions.ConfigFlags) *cobra.Command {
	var allNamespaces bool
	cmd := &cobra.Command{
		Use:   "browser",
		Short: "Open url defined in ingress with browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				return err
			}
			util.PrintError(browser(flags, host, allNamespaces))
			return nil
		},
	}
	cmd.Flags().String("host", "", "Show just the ingress definitions for this hostname")
	cmd.Flags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "Find ingress definitions from all namespaces")
	return cmd
}

func browser(flags *genericclioptions.ConfigFlags, host string, allNamespaces bool) error {
	k8sclient.GlobalClient(flags) // try to see if there're any errors
	if fzf.IsInteractiveMode(os.Stdout) {
		return fzf.FzfRun(color.Error)
	}

	rows, err := ingresses.GetIngressRow(flags, allNamespaces)
	if err != nil {
		return err
	}
	var urls []string
	for _, row := range rows {
		if len(row.Host) == 0 || (len(host) > 0 && row.Host != host) {
			continue
		}
		scheme := "http"
		if row.TLS {
			scheme = "https"
		}
		urls = append(urls, fmt.Sprintf("%s://%s%s", scheme, row.Host, row.Path))
	}
	if len(urls) == 0 {
		return nil
	}
	for _, url := range urls {
		fmt.Fprintln(color.Output, url)
	}
	return nil
}
