package klusterhelper

import (
	"fmt"
	"reflect"

	"sigs.k8s.io/yaml"
)

func clean(obj any) any {
	switch v := obj.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			cleaned := clean(value)
			if !isEmpty(cleaned) {
				result[key] = cleaned
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	case []any:
		var result []any
		for _, item := range v {
			cleaned := clean(item)
			if !isEmpty(cleaned) {
				result = append(result, cleaned)
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return obj
	}
}

func isEmpty(obj any) bool {
	if obj == nil {
		return true
	}

	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

func marshalCleanYAML(obj any) ([]byte, error) {
	// k8s only use omitempty on json tags, but not yaml tags
	// so we need to round about remove empty fields
	// also, because the structs are not pointers, they're not really empty so json marshalling
	// ignores it too. so we're stuck with this ugly approach
	var generic map[string]any
	b, err := yaml.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}
	yaml.Unmarshal(b, &generic)

	cleaned := clean(generic)

	d, err := yaml.Marshal(cleaned)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cleaned object: %w", err)
	}
	return d, nil
}
