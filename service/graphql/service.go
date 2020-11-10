package graphql

import (
	"fmt"
	tools "github.com/bhoriuchi/graphql-go-tools"
	"github.com/graphql-go/graphql"
)

func ResolveSchema(schemaContent []byte) *graphql.Schema {
	schema, err := tools.MakeExecutableSchema(tools.ExecutableSchema{
		TypeDefs: string(schemaContent),
		SchemaDirectives: tools.SchemaDirectiveVisitorMap{
			"IsGranted": &tools.SchemaDirectiveVisitor{
				VisitFieldDefinition: IsGranted,
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to parse schema, error: %v", err))
	}

	return &schema
}
