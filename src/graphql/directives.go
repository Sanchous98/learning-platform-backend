package graphql

import "github.com/graphql-go/graphql"

func IsGranted(field *graphql.Field, args map[string]interface{}) {
	resolveFunc := field.Resolve
	field.Resolve = func(p graphql.ResolveParams) (interface{}, error) {
		// TODO: Check rights before resolving
		result, err := resolveFunc(p)
		if err != nil {
			return result, err
		}
		return make(map[string]interface{}), nil
	}
}
