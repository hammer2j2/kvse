package kvse_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	kvse "github.com/hammer2j2/kvse/cmd/kvse"

	mock_kvse "github.com/hammer2j2/kvse/pkg/kvse/mocks"

	"github.com/stretchr/testify/assert"
)

func TestReadCmdDefaultFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// mock in the test
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKvse := mock_kvse.NewMockKvse(mockCtrl)
	mockKvse.EXPECT().SetupRequest(
		"myproj", "config.yaml", []string{}).
		DoAndReturn(
			func(factSpec string, configFile string, options []string) {
				assert.Equal(t, factSpec, "myproj", "config.yaml")
			})
	mockKvse.EXPECT().Read().Times(1)

	var configFlag string
	var hideTrunkFlag bool
	readCmd := kvse.NewReadCommand(mockKvse)
	readCmd.PersistentFlags().StringVarP(&configFlag, kvse.ConfigFileOptionStr, string(kvse.ConfigFileOptionStr[0]), "", "configuration file path")
	readCmd.PersistentFlags().BoolVarP(&hideTrunkFlag, kvse.HideTrunkOptionStr, string(kvse.HideTrunkOptionStr[0]), false, "hide the trunk input to read from the results")

	// given arguments to read root element "myproj" with an override config file
	os.Args = []string{"read", "myproj"}

	// then executing the command will result in the expected arguments
	readCmd.Execute()
}

func TestReadCmdWithOverrideFilePath(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// mock in the test
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKvse := mock_kvse.NewMockKvse(mockCtrl)
	mockKvse.EXPECT().SetupRequest(
		"myproj", "./overrideconfigfile.yaml", []string{}).
		DoAndReturn(
			func(factSpec string, configFile string, options []string) {
				assert.Equal(t, factSpec, "myproj")
			})
	mockKvse.EXPECT().Read().Times(1)

	var configFlag string
	var hideTrunkFlag bool
	readCmd := kvse.NewReadCommand(mockKvse)
	readCmd.PersistentFlags().StringVarP(&configFlag, kvse.ConfigFileOptionStr, string(kvse.ConfigFileOptionStr[0]), "", "configuration file path")
	readCmd.PersistentFlags().BoolVarP(&hideTrunkFlag, kvse.HideTrunkOptionStr, string(kvse.HideTrunkOptionStr[0]), false, "hide the trunk input to read from the results")

	// given arguments to read root element "myproj" with an override config file
	os.Args = []string{"read", "myproj", "-f", "./overrideconfigfile.yaml"}

	// then executing the command will result in the expected arguments
	readCmd.Execute()
}

func TestReadCmdWithHideOption(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// mock in the test
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockKvse := mock_kvse.NewMockKvse(mockCtrl)
	mockKvse.EXPECT().SetupRequest(
		"myproj", "./overrideconfigfile.yaml", []string{"hideTrunk"}).
		DoAndReturn(
			func(factSpec string, configFile string, options []string) {
				assert.Equal(t, factSpec, "myproj")
			})
	mockKvse.EXPECT().Read().Times(1)

	var configFlag string
	var hideTrunkFlag bool
	readCmd := kvse.NewReadCommand(mockKvse)
	readCmd.PersistentFlags().StringVarP(&configFlag, kvse.ConfigFileOptionStr, string(kvse.ConfigFileOptionStr[0]), "", "configuration file path")
	readCmd.PersistentFlags().BoolVarP(&hideTrunkFlag, kvse.HideTrunkOptionStr, string(kvse.HideTrunkOptionStr[0]), false, "hide the trunk input to read from the results")

	// given arguments to read root element "myproj" with an override config file and the hide trunk option
	os.Args = []string{"read", "myproj", "-f", "./overrideconfigfile.yaml", "-t"}

	// then executing the command will result in the expected arguments
	readCmd.Execute()
}
