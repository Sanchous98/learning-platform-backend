package graphql

import (
	"fmt"
	tools "github.com/bhoriuchi/graphql-go-tools"
	"github.com/graphql-go/graphql"
)

func ResolveSchema(schemaContent []byte) *graphql.Schema {
	schema, err := tools.MakeExecutableSchema(tools.ExecutableSchema{
		TypeDefs: string(schemaContent),
		Resolvers: tools.ResolverMap{
			"Query": &tools.ObjectResolver{
				Fields: tools.FieldResolveMap{
					"user": func(p graphql.ResolveParams) (interface{}, error) {
						// lookup data
						return map[string]interface{}{"email": "bar"}, nil
					},
				},
			},
		},
		SchemaDirectives: tools.SchemaDirectiveVisitorMap{
			"IsGranted": &tools.SchemaDirectiveVisitor{
				VisitFieldDefinition: IsGranted,
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to build schema, error: %v", err))
	}

	return &schema
}
