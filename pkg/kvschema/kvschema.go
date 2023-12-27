package kvschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrNameNotFound = errors.New("name not found in list")
	logger          *log.Logger
)

const (
	separator   = "."
	packageName = "kvschema"
)

func init() {
	logger = log.New(os.Stderr, packageName+" ", log.LstdFlags)
}

type Map map[string]json.RawMessage
type Array []json.RawMessage

var m map[string]string
var debug bool = false
var errChan chan error = make(chan error)

func Transform(j []byte, options []string) (map[string]string, chan error) {
	var r []byte
	go addError(errChan, nil)
	if slices.Contains(options, "debug") {
		debug = true
	}
	m = make(map[string]string) // reset this for each call
	r = json.RawMessage(j)
	walkJson(r, "")
	return m, errChan
}

func walkJson(raw json.RawMessage, stack string) {
	var val interface{}
	sep := separator
	json.Unmarshal(raw, &val)

	switch val.(type) {
	case map[string]interface{}:
		var cont Map
		if stack == "" {
			sep = ""
		}
		json.Unmarshal(raw, &cont)
		for key, value := range cont {
			walkJson(value, stack+sep+key)
		}
		return
	case []interface{}:
		var val []interface{}
		var nameList []string
		var scalarList []string
		json.Unmarshal(raw, &val)

		parentStack := stack              // save this from updateStack calls later
		for count, element := range val { // for each element of the list
			switch element.(type) {
			case map[string]interface{}: // if it's a map, special handling to extract the 'name' schema element
				var jsonObj Map
				if stack == "" {
					sep = ""
				}
				b, err := json.Marshal(element)
				if err != nil {
					go addError(errChan, fmt.Errorf("Cannot marshal element: %v", err))
					logger.Printf("Error: Cannot marshal element: %v", err)
				}
				json.Unmarshal(b, &jsonObj)
				nameKey, err := getNameFromMap(b)
				if err != nil {
					panic(err)
				}
				// comprise a list of names in the original stack element
				nameList = append(nameList, nameKey)
				if count == 0 {
					safeWrite(parentStack, strings.Join(nameList, ","))
				} else {
					m[parentStack] = strings.Join(nameList, ",") // update parent stack with current list
				}

				stack = shiftStack(stack, nameKey) // get the new name-based stack

				for key, value := range jsonObj {
					walkJson(value, stack+sep+key)
				}
			case string: // a list of scalar strings so we need to join them into 1 scalar
				if debug {
					logger.Printf("type default: %T\n", element)
				}
				s := element.(string)
				scalarList = append(scalarList, s)
				m[stack] = strings.Join(scalarList, ",")
			default:
				if debug {
					logger.Printf("type default has no routine: %T\n", element)
				}
			}
		}
		return
	case float64:
		safeWrite(stack, strconv.FormatFloat(val.(float64), 'f', -1, 64))
	case string:
		safeWrite(stack, val.(string))
	case bool:
		safeWrite(stack, strconv.FormatBool(val.(bool)))
	case nil:
		safeWrite(stack, "null")
	}
}

func safeWrite(key string, value string) {
	if _, ok := m[key]; ok {
		go addError(errChan, fmt.Errorf("Cannot overwrite existing value: %s=%s with %s\n", key, m[key], value))
		logger.Printf("Error: Cannot overwrite existing value: %s=%s with %s\n", key, m[key], value)
		return
	}
	m[key] = value
}

// getNameFromMap searches for an object with a key of 'name' and
// returns the value and any error
func getNameFromMap(raw json.RawMessage) (string, error) {
	var jsonMap Map
	json.Unmarshal(raw, &jsonMap)
	for key, rawValue := range jsonMap {
		if key == "name" {
			var val interface{}
			json.Unmarshal(rawValue, &val)
			strVal := val.(string)
			return strVal, nil
		}
	}
	return "", fmt.Errorf("Cannot find name key in object list")
}

// shiftStack takes a stack string (separator separated string) and a nameKey and returns a new stack
// with the nameKey overwriting the last element in the separated stack list.  The purpose in overwriting
// elements in the stack found in the source tree is to remove 'schema elements' from the resulting map keys.
// returns a new assembled stack string
func shiftStack(stack, nameKey string) string {
	stackArr := strings.Split(stack, separator)
	stackArr[len(stackArr)-1] = nameKey
	newStack := strings.Join(stackArr, separator)
	if _, ok := m[newStack]; ok {
		go addError(errChan, fmt.Errorf("Cannot shift into existing stack: %s\n", newStack))
		logger.Printf("Error: Cannot shift into existing stack: %s\n", newStack)
		return stack
	}
	return newStack
}

func addError(c chan error, e error) {
	if e != nil {
		_ = <-c // remove the nil from the channel
	}
	c <- e
}
