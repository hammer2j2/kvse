package kvse

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hammer2j2/kvse/pkg/kvschema"
	"sigs.k8s.io/yaml"
)

const (
	DEFAULT_SECTION  = "default"
	CONFIG_FILE_NAME = "kvse.ini"
	programName      = "kvse"
)

var (
	kvseFileName string = "config.yaml"
	logger       *log.Logger
)

func init() {
	logger = log.New(os.Stderr, programName+" ", log.LstdFlags)
}

type Kvse interface {
	SetupRequest(string, string, []string)
	Read() (map[string]string, error)
}

type FactRequest struct {
	FactSpec   string
	ConfigFile string
	Options    []string
}

func (f *FactRequest) SetupRequest(
	factSpec string,
	configFile string,
	options []string) {
	f.FactSpec = factSpec
	f.ConfigFile = configFile
	f.Options = options
}

func NewFactRequest() *FactRequest {
	return &FactRequest{}
}

func (f FactRequest) Read() (value map[string]string, err error) {
	var envFile string
	debug := slices.Contains(f.Options, "debug")
	includeNames := slices.Contains(f.Options, "includeNames")

	if debug {
		logger.Printf("fact request: %v", f)
	}
	if envFile, err = GetFilePath(f.ConfigFile, strconv.FormatBool(debug)); err != nil {
		return value, fmt.Errorf("Cannot find config file: %s: %v", f.ConfigFile, errors.Unwrap(err))
	}
	if debug {
		logger.Printf("Reading %s file: %s", programName, envFile)
	}
	data, err := os.ReadFile(envFile)
	if err != nil {
		return value, err
	}

	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return value, err
	}

	kv, chanErr := kvschema.Transform(jsonData, f.Options)

	if err := <-chanErr; err != nil {
		return kv, err
	}

	regexp := regexp.MustCompile(`^` + regexp.QuoteMeta(f.FactSpec) + `.*$`)

	resultSet := make(map[string]string)

	for k, v := range kv {
		if debug {
			logger.Printf("%s: %s\n", k, v)
		}
		if f.FactSpec == "." || regexp.Match([]byte(k)) {
			if includeNames || !strings.HasSuffix(k, "name") { // exclude names unless wanted
				resultSet[k] = v
			}
		}
	}
	return resultSet, nil
}

// GetFilePath accepts a filename or full file path. If given a filename it looks
// for the file in the current and executing program directories
// It checks for ability to read the provided file or file path and
// returns the full path of the first location if found or the empty string and an error
func GetFilePath(params ...string) (_ string, err error) {
	var debug bool
	if len(params) > 1 {
		optList := params[1:] // strip away filename param
		if slices.Contains(optList, "debug") {
			debug = true
		}
	}
	filename := params[0]
	pathexp := regexp.MustCompile(`^(/|\.?/|~?/|\.?\\|~?\\)`)

	if debug {
		logger.Printf("checking for path expression in filename: %s", filename)
	}
	if found := pathexp.Match([]byte(filename)); found { // filename is a path
		if debug {
			logger.Println("filename is a path")
		}
		_, err = os.Stat(filename)
		if err == nil {
			return filename, nil
		} else {
			return "", err
		}
	}
	if debug {
		logger.Println("filename is not a path")
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("calling os.Getwd(): %v", err)
	}
	for _, path := range []string{exPath, cwd} {
		c := filepath.Join(path, filename)
		_, err = os.Stat(c)
		if err == nil {
			return c, nil
		}
	}
	return "", err
}
