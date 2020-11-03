package graphql

import "github.com/graphql-go/graphql"

func IsGranted(field *graphql.Field, args map[string]interface{}) {
	resolveFunc := field.Resolve
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		result, err := resolveFunc(p)
		if err != nil {
			return result, err
		}
		data := result.(map[string]interface{})
		data["description"] = args["value"]
		return data, nil
	}
}
