package cmd

import (
	"github.com/misrab/clai/internal/webui"
	"github.com/spf13/cobra"
)

var (
	webuiPort      int
	webuiNoBrowser bool

	webuiCmd = &cobra.Command{
		Use:   "webui",
		Short: "Start the web UI",
		Long:  "Start a local web server and open the clai web interface in your browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			return webui.Start(webuiAssets, webuiPort, !webuiNoBrowser)
		},
	}
)

func init() {
	webuiCmd.Flags().IntVarP(&webuiPort, "port", "p", 8080, "Port to run web server on")
	webuiCmd.Flags().BoolVar(&webuiNoBrowser, "no-browser", false, "Don't auto-open browser")
	rootCmd.AddCommand(webuiCmd)
}
