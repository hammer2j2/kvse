package kvschema_test

import (
	"testing"

	"github.com/hammer2j2/kvse/pkg/kvschema"
	"github.com/stretchr/testify/assert"
)

func TestTransformMap(t *testing.T) {
	// Test case 1: Test map[string]interface{}
	raw1 := []byte(`{"key1": "value1", "key2": 123, "key3": true}`)
	expected1 := map[string]string{
		"key1": "value1",
		"key2": "123",
		"key3": "true",
	}
	result, chanErr := kvschema.Transform(raw1, []string{})
	err := <-chanErr
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, expected1, result)
}

// TODO: handle this case gracefully
// // Test case 2: Test []interface{}
// raw2 := []byte(`[{"person": {"firstname": "John", "age": 30}}]`)
// expected2 := map[string]string{
// 	"person.firstname": "John",
// 	"person.age":       "30",
// }
// result = kvschema.Transform(raw2)
// assert.Equal(t, expected2, result)

func TestTransformNameReplacement(t *testing.T) {
	// // Test case 3: name replacement
	raw3 := []byte(`{"persons": [{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]}`)
	expected3 := map[string]string{
		"John.age":  "30",
		"John.name": "John",
		"Jane.age":  "25",
		"Jane.name": "Jane",
		"persons":   "John,Jane",
	}
	result, chanErr := kvschema.Transform(raw3, []string{})
	err := <-chanErr
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, expected3, result)
}

func TestTransformStringArrays(t *testing.T) {
	// // Test case 4: Test string arrays
	raw4 := []byte(`{"persons": [{"name": "John", "age": 30, "experience": ["java","node","go"]}, {"name": "Jane", "age": 25}]}`)
	expected4 := map[string]string{
		"John.age":        "30",
		"John.name":       "John",
		"John.experience": "java,node,go",
		"Jane.age":        "25",
		"Jane.name":       "Jane",
		"persons":         "John,Jane",
	}
	result, chanErr := kvschema.Transform(raw4, []string{})
	err := <-chanErr
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, expected4, result)
}

func TestTransformBool(t *testing.T) {
	// // Test case 5: Test bool
	raw5 := []byte(`{"persons": [{"name": "John", "age": 30, "available": true}, {"name": "Jane", "age": 25}]}`)
	expected5 := map[string]string{
		"John.age":       "30",
		"John.name":      "John",
		"John.available": "true",
		"Jane.age":       "25",
		"Jane.name":      "Jane",
		"persons":        "John,Jane",
	}
	result, chanErr := kvschema.Transform(raw5, []string{})
	err := <-chanErr
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, expected5, result)
}

func TestTransformNil(t *testing.T) {
	// // Test case 6: Test nil
	raw6 := []byte(`{"persons": [{"name": "John", "age": 30, "available": null}, {"name": "Jane", "age": 25}]}`)
	expected6 := map[string]string{
		"John.age":       "30",
		"John.name":      "John",
		"John.available": "null",
		"Jane.age":       "25",
		"Jane.name":      "Jane",
		"persons":        "John,Jane",
	}
	result, chanErr := kvschema.Transform(raw6, []string{})
	err := <-chanErr
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, expected6, result)
}

func TestDuplicateKeyErrors(t *testing.T) {
	// // Test case 6: Test nil
	raw6 := []byte(`{"persons": [{"name": "John", "age": 30, "available": null}, {"name": "Jane", "age": 25}],"pets": [{"name": "John", "age": 5, "available": true}]}`)
	_, chanErr := kvschema.Transform(raw6, []string{})
	err := <-chanErr
	assert.ErrorContains(t, err, "Cannot overwrite existing value")
}
