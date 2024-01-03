package cli

import (
 "fmt"
 "os"

 "github.com/spf13/cobra"
)

var version = "0.0.6"

var rootCmd = &cobra.Command{
    Use:  "mem",
    Version: version,
    Long: `Memphis Utility`,
    Run: func(cmd *cobra.Command, args []string) {
        
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
        os.Exit(1)
    }
}