package kvse_test

import (
	"fmt"
	fs "io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hammer2j2/kvse/pkg/kvse"

	"github.com/stretchr/testify/assert"
)

var (
	err                    error
	cwd                    string
	exPath                 string
	executable             string
	expected_filename      = "test_expected_file.yaml"
	expected_file_dir      string
	expected_file_path     string
	non_existent_filename  = "non-existent_file.yaml"
	test_data_dir          = "testdata"
	test_filename          = "testconfig.yaml"
	test_file_path         string
	non_existent_file_path string
	filePerm               fs.FileMode = 0755
	testData                           = []byte(`myproj:
  envs:
    - name: dev
      regions:
      - name: us-west-2
      - name: us-east-1
      account: 1234
    - name: qa
      regions:
      - us-west-2
      - us-east-1
      account: 4567
#	`)
)

// add a setup and teardown
func setupTest() func() {
	// Setup code here
	executable, err = os.Executable()
	if err != nil {
		panic(err)
	}
	exPath = filepath.Dir(executable)
	cwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("calling os.Getwd(): %v", err))
	}
	expected_file_path = filepath.Join(exPath, expected_filename)
	non_existent_file_path = filepath.Join(exPath, non_existent_filename)

	test_file_path = filepath.Join(cwd, test_data_dir, test_filename)

	// create one expected file to start with
	if err := os.WriteFile(expected_file_path, testData, filePerm); err != nil {
		panic(fmt.Sprintf("cannot write file: %s with error: %v", expected_file_path, err))
	}

	// make sure non existent file doesn't exist
	if _, err := os.Stat(non_existent_file_path); err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Sprintf("unexpected non_existent_file: %s with error: %v", non_existent_file_path, err))
		}
	} else {
		panic(fmt.Sprintf("unexpected non_existent_file: %s with error: %v", non_existent_file_path, err))
	}

	// tear down
	return func() {
		// tear-down code here
		for _, path := range []string{cwd, exPath} {
			expected_file_path := filepath.Join(path, expected_filename)
			if _, err := os.Stat(expected_file_path); err != nil {
				if !os.IsNotExist(err) {
					panic(fmt.Sprintf("unexpected non_existent_file: %s with error: %v", non_existent_file_path, err))
				}
			} else { // found it so remove it
				if err := os.Remove(expected_file_path); err != nil {
					panic(fmt.Sprintf("cannot remove expected file: %s\n", expected_file_path))
				}
			}
		}
	}
}
func TestReadFileDefaultPath(t *testing.T) {
	defer setupTest()()
	// when read is called with filename only
	// GetFilePath can find it
	expected := expected_filename
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj", expected, []string{})
	_, err := f.Read()
	assert.NoError(t, err)
}

func TestReadCustomFilePath(t *testing.T) {
	defer setupTest()()
	// when read is called with alternate config file
	// GetFilePath can find it
	expected := test_file_path
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj", expected, []string{})
	_, err := f.Read()
	assert.NoError(t, err)
}

func TestReadCustomFilePathNotFound(t *testing.T) {
	defer setupTest()()
	// when read is called with alternate config file that doesn't exist
	// GetFilePath returns wrapped ENOENT error
	expected := non_existent_file_path
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj", expected, []string{})
	expectedError := fmt.Errorf("Cannot find config file: %s: %v", expected,
		error(syscall.ENOENT))
	_, actualError := f.Read()
	assert.ErrorAs(t, actualError,
		&expectedError, "Read should return error of 'no such file or directory'")
}

func TestGetFilePathExecutableDir(t *testing.T) {
	var actual string
	defer setupTest()()
	// when expected file does not contain a directory path
	// and the config file is in the executable directory
	// GetFilePath can find it
	expected := expected_file_path
	actual, err := kvse.GetFilePath(expected_filename)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
func TestGetFilePathCWD(t *testing.T) {
	var actual string
	defer setupTest()()
	cwd, err = os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("calling os.Getwd(): %v", err))
	}

	expected := filepath.Join(cwd, expected_filename)
	if err := os.Rename(expected_file_path, expected); err != nil {
		panic(fmt.Errorf("cannot move expected file: %s with error: %v\n", expected_file_path, err))
	}
	// when expected file does not contain a directory path
	// and the config file is in the current directory
	// GetFilePath can find it
	actual, err := kvse.GetFilePath(expected_filename)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetFilePathNotExist(t *testing.T) {
	// when an expected file does not exist
	// GetFilePath returns file not found error
	_, err = kvse.GetFilePath(non_existent_filename)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "no such file or directory", "the err should contain 'no such file or directory'")
}

func TestGetFileFullPath(t *testing.T) {
	var actual string
	defer setupTest()()
	// when expected file contains a directory path
	// GetFilePath can find it
	expected := expected_file_path
	actual, err := kvse.GetFilePath(expected_file_path)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetSpecificKey(t *testing.T) {
	defer setupTest()()
	// when read is called with a specific key
	// GetFilePath can find it
	expected := make(map[string]string)
	expected["myproj.dev.account"] = "1234"
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj.dev.account", expected_file_path, []string{})
	result, err := f.Read()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetKeySet(t *testing.T) {
	defer setupTest()()
	// when read is called with a specific key
	// GetFilePath can find it
	expected := map[string]string{
		"myproj.dev.account":        "1234",
		"myproj.dev.us-west-2.name": "us-west-2",
		"myproj.dev.us-east-1.name": "us-east-1",
		"myproj.dev.name":           "dev",
		"myproj.dev.regions":        "us-west-2,us-east-1",
	}
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj.dev", expected_file_path, []string{"debug", "includeNames"})
	result, err := f.Read()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetKeySetWithScalarList(t *testing.T) {
	defer setupTest()()
	// when read is called with a specific key
	// GetFilePath can find it
	expected := map[string]string{
		"myproj.qa.account": "4567",
		"myproj.qa.name":    "qa",
		"myproj.qa.regions": "us-west-2,us-east-1",
	}
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj.qa", expected_file_path, []string{"debug", "includeNames"})
	result, err := f.Read()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetKeySetWithScalarListAndNoIncludeNames(t *testing.T) {
	defer setupTest()()
	// when read is called with a specific key
	// GetFilePath can find it
	expected := map[string]string{
		"myproj.qa.account": "4567",
		"myproj.qa.regions": "us-west-2,us-east-1",
	}
	f := kvse.NewFactRequest()
	f.SetupRequest("myproj.qa", expected_file_path, []string{"debug"})
	result, err := f.Read()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
