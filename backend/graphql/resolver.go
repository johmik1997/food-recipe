// package graphql

// import (
// 	"foodrecipe/database"
// 	"github.com/graphql-go/graphql"
// )

// type Resolver struct {
// 	DB *database.DB
// }

// // LoginResolver resolves the login mutation
// func (r *Resolver) LoginResolver(p graphql.ResolveParams) (interface{}, error) {
// 	// Implement your login logic here
// 	return nil, nil
// }

// // MeResolver resolves the me query
// func (r *Resolver) MeResolver(p graphql.ResolveParams) (interface{}, error) {
// 	// Implement your logic to fetch the current user here
// 	return nil, nil
// }

// // RegisterResolver resolves the register mutation
// func (r *Resolver) RegisterResolver(p graphql.ResolveParams) (interface{}, error) {
// 	// Implement your register logic here
// 	return nil, nil
// }
// var userType = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "User",
// 	Fields: graphql.Fields{
// 		"id": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"name": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"email": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"created_at": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"updated_at": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 	},
// })
// var loginResponseType = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "LoginResponse",
// 	Fields: graphql.Fields{
// 		"token": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"user": &graphql.Field{

// 			Type: userType,
// 		},
// 	},
// })
// var registerResponseType = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "RegisterResponse",
// 	Fields: graphql.Fields{
// 		"user": &graphql.Field{
// 			Type: userType,
// 		},
// 	},
// })
// var registerInputType = graphql.NewInputObject(graphql.InputObjectConfig{
// 	Name: "RegisterInput",
// 	Fields: graphql.InputObjectConfigFieldMap{
// 		"name": &graphql.InputObjectFieldConfig{
// 			Type: graphql.NewNonNull(graphql.String),
// 		},
// 		"email": &graphql.InputObjectFieldConfig{
// 			Type: graphql.NewNonNull(graphql.String),
// 		},
// 		"password": &graphql.InputObjectFieldConfig{
// 			Type: graphql.NewNonNull(graphql.String),
// 		},
// 	},
// })
// var loginInputType = graphql.NewInputObject(graphql.InputObjectConfig{
// 	Name: "LoginInput",
// 	Fields: graphql.InputObjectConfigFieldMap{
// 		"email": &graphql.InputObjectFieldConfig{
// 			Type: graphql.NewNonNull(graphql.String),
// 		},
// 		"password": &graphql.InputObjectFieldConfig{
// 			Type: graphql.NewNonNull(graphql.String),
// 		},
// 	},
// })
// var loginQuery = `
// mutation Login($email: String!, $password: String!) {
// 	login(email: $email, password: $password) {
// 		token
// 		user {
// 			id
// 			name
// 			email
// 			created_at
// 			updated_at
// 		}
// 	}
// }
// `
// var registerQuery = `
// mutation Register($input: RegisterInput!) {
// 	register(input: $input) {
// 		user {
// 			id
// 			name
// 			email
// 			created_at
// 			updated_at
// 		}
// 	}
// }
// `
// var meQuery = `
// query Me {
// 	me {
// 		id
// 		name
// 		email
// 		created_at
// 		updated_at
// 	}
// }
// `

// func NewSchema(db *database.DB) (graphql.Schema, error) {
// 	resolver := &Resolver{DB: db}

// 	return graphql.NewSchema(graphql.SchemaConfig{
// 		Query:    resolver.QueryType(),
// 		Mutation: resolver.MutationType(),
// 	})
// }

// func (r *Resolver) QueryType() *graphql.Object {
// 	return graphql.NewObject(graphql.ObjectConfig{
// 		Name: "Query",
// 		Fields: graphql.Fields{
// 			"me": &graphql.Field{
// 				Type: userType,
// 				Resolve: r.MeResolver,
// 			},
// 		},
// 	})
// }

//	func (r *Resolver) MutationType() *graphql.Object {
//		return graphql.NewObject(graphql.ObjectConfig{
//			Name: "Mutation",
//			Fields: graphql.Fields{
//				"login": &graphql.Field{
//					Type: loginResponseType,
//					Args: graphql.FieldConfigArgument{
//						"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
//						"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
//					},
//					Resolve: r.LoginResolver,
//				},
//				"register": &graphql.Field{
//					Type: registerResponseType,
//					Args: graphql.FieldConfigArgument{
//						"input": &graphql.ArgumentConfig{Type: registerInputType},
//					},
//					Resolve: r.RegisterResolver,
//				},
//			},
//		})
//	}
package graphql