/*package parser

type Transformation struct {
	Field     string
	KeyName   string
	ValueName string
}

var Transformations = map[string][]Transformation{
	"dlllist": {
		{Field: "DLLList", KeyName: "path", ValueName: "status"},
	},
}

func ApplyTransformations(doc map[string]interface{}, eventType string) {
	transforms, ok := Transformations[eventType]
	if !ok {
		return
	}

	for _, t := range transforms {
		if obj, ok := doc[t.Field].(map[string]interface{}); ok {
			doc[t.Field] = convertObjectToArray(obj, t.KeyName, t.ValueName)
		}
	}
}

func convertObjectToArray(obj map[string]interface{}, keyName, valueName string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(obj))
	for k, v := range obj {
		result = append(result, map[string]interface{}{
			keyName:   k,
			valueName: v,
		})
	}
	return result
}
*/