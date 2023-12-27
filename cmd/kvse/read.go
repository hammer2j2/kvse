package kvse

import (
	"fmt"
	"log"
	"os"

	pkgkvse "github.com/hammer2j2/kvse/pkg/kvse"
	"github.com/spf13/cobra"
)

var (
	err              error
	result           string
	configFile       string
	configFileFlag   string
	optionHideTrunk  bool
	hideTrunkFlag    bool
	optionDebug      bool
	debugFlag        bool
	optionValuesOnly bool
	valuesOnlyFlag   bool
	factRequest      = &pkgkvse.FactRequest{}
	logger           *log.Logger
)

func init() {
	logger = log.New(os.Stderr, programName+" ", log.LstdFlags)
}

func NewReadCommand(factRequest pkgkvse.Kvse) *cobra.Command {
	return &cobra.Command{
		Use:   usage,
		Short: "reads the given fact",
		Args:  cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {

			var options = []string{}
			var kv map[string]string
			configFile := getReadOptionConfigFile(cmd)
			options = getReadOptionTrunk(cmd, options)
			options = getReadOptionDebug(cmd, options)
			valuesOnly := getReadOptionValuesOnly(cmd, options)

			factRequest.SetupRequest(
				args[0],
				configFile,
				options)

			if kv, err = factRequest.Read(); err != nil {
				logger.Fatal("Unhandled error calling Read()", err)
			}
			for k, v := range kv {
				if valuesOnly {
					fmt.Printf("%s\n", v)
				} else {
					fmt.Printf("%s=%s\n", k, v)
				}
			}
		},
	}
}

func getReadOptionConfigFile(cmd *cobra.Command) string {
	configFile, err := cmd.Flags().GetString(ConfigFileOptionStr)
	if configFile == "" || err != nil {
		configFile = defaultConfigFile
	}
	return configFile
}

func getReadOptionTrunk(cmd *cobra.Command, options []string) []string {
	trunk, _ := cmd.Flags().GetBool(string(HideTrunkOptionStr[0]))
	if trunk {
		options = append(options, "hideTrunk")
	}
	return options
}

func getReadOptionDebug(cmd *cobra.Command, options []string) []string {
	debug, _ := cmd.Flags().GetBool(string(DebugOptionStr[0]))
	if debug {
		options = append(options, "debug")
	}
	return options
}

func getReadOptionValuesOnly(cmd *cobra.Command, options []string) bool {
	valuesOnly, _ := cmd.Flags().GetBool(string(ValuesOnlyOptionStr[0]))
	return valuesOnly
}

func init() {
	readCmd := NewReadCommand(factRequest)
	readCmd.PersistentFlags().StringVarP(&configFileFlag, ConfigFileOptionStr, string(ConfigFileOptionStr[0]), "", "configuration file path")
	readCmd.PersistentFlags().BoolVarP(&hideTrunkFlag, HideTrunkOptionStr, string(HideTrunkOptionStr[0]), false, "hide the trunk input to read from the results")
	readCmd.PersistentFlags().BoolVarP(&debugFlag, DebugOptionStr, string(DebugOptionStr[0]), false, "debug output")
	readCmd.PersistentFlags().BoolVarP(&valuesOnlyFlag, ValuesOnlyOptionStr, string(ValuesOnlyOptionStr[0]), false, "values only output")
	RootCmd.AddCommand(readCmd)
}
