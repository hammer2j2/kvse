package kvse

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	version             = "0.0.2"
	defaultConfigFile   = "config.yaml"
	programName         = "kvse"
	usage               = "read <fact> [ -f <file> ]"
	ConfigFileOptionStr = "file"
	HideTrunkOptionStr  = "t"
	DebugOptionStr      = "d"
	ValuesOnlyOptionStr = "v"
)

var RootCmd = &cobra.Command{
	Version: version,
	Use:     programName,
	Short:   programName + " transforms a " + programName + " yaml file to a flat map of key-value pairs",
	Long: `
   
One can use kvse to transform and retrieve configuration in a single command`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(cmd.OutOrStdout(), "Error: missing command argument\n")
			if err := cmd.Help(); err != nil {
				logger.Fatal("calling Help()", err)
			}
		}
	},
}

func init() {
	logger = log.New(os.Stderr, programName+" ", log.LstdFlags)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}
