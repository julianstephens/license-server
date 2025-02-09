// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/products": {
            "get": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "retrieves all products",
                "tags": [
                    "products"
                ],
                "summary": "Get all products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-array_model_Product"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "creates a new product",
                "tags": [
                    "products"
                ],
                "summary": "Add a product",
                "parameters": [
                    {
                        "description": "new product info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/admin/products/:id": {
            "get": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "retrieve a specific product",
                "tags": [
                    "products"
                ],
                "summary": "Get a product",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "updates a specific product",
                "tags": [
                    "products"
                ],
                "summary": "Update a product",
                "parameters": [
                    {
                        "description": "updated product info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Product"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "deletes a specific product",
                "tags": [
                    "products"
                ],
                "summary": "Delete a product",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_Product"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/admin/products/:id/key": {
            "get": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "creates an ed25519 key pair for a specific product and version",
                "tags": [
                    "products"
                ],
                "summary": "Add a product key pair",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_ProductKeyPair"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/admin/users": {
            "get": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "retrieves all users",
                "tags": [
                    "users"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-array_model_User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "creates a new user",
                "tags": [
                    "users"
                ],
                "summary": "Add a user",
                "parameters": [
                    {
                        "description": "new user info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/admin/users/:id": {
            "get": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "retrieve a specific user",
                "tags": [
                    "users"
                ],
                "summary": "Get a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_UserWithScopes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "updates a specific user",
                "tags": [
                    "users"
                ],
                "summary": "Update a user",
                "parameters": [
                    {
                        "description": "updated user info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "deletes a specific user",
                "tags": [
                    "users"
                ],
                "summary": "Delete a user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/admin/users/:id/scopes": {
            "put": {
                "description": "update scopes for a specific user",
                "tags": [
                    "users"
                ],
                "summary": "Update a user's scope",
                "parameters": [
                    {
                        "description": "scopes to modify",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.ScopeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_UserWithScopes"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "registers new application user",
                "tags": [
                    "auth"
                ],
                "summary": "Register a user",
                "parameters": [
                    {
                        "description": "new user info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/auth/token": {
            "post": {
                "description": "creates a new API key",
                "tags": [
                    "auth"
                ],
                "summary": "Create a token",
                "parameters": [
                    {
                        "description": "returning user info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_DisplayAPIKey"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/licenses/:id/revoke": {
            "delete": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "revokes a license with id",
                "tags": [
                    "licenses"
                ],
                "summary": "Revoke a license",
                "parameters": [
                    {
                        "description": "license id",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LicenseRevokeRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/licenses/issue": {
            "post": {
                "security": [
                    {
                        "ApiKey": []
                    }
                ],
                "description": "issues a new product license",
                "tags": [
                    "licenses"
                ],
                "summary": "Issue a license",
                "parameters": [
                    {
                        "description": "new license info",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.LicenseRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPResponse-model_ActivationData"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.AuthRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "controller.ScopeRequest": {
            "type": "object",
            "properties": {
                "add": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "remove": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "httputil.HTTPResponse-array_model_Product": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Product"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-array_model_User": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.User"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_ActivationData": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.ActivationData"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_DisplayAPIKey": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.DisplayAPIKey"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_Product": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.Product"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_ProductKeyPair": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.ProductKeyPair"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_User": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.User"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "httputil.HTTPResponse-model_UserWithScopes": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.UserWithScopes"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "model.ActivationData": {
            "type": "object",
            "properties": {
                "expiration_date": {
                    "type": "integer"
                },
                "issue_date": {
                    "type": "integer"
                },
                "license_id": {
                    "type": "string"
                },
                "product": {
                    "type": "string"
                },
                "product_key": {
                    "type": "string"
                },
                "refresh_date": {
                    "type": "integer"
                }
            }
        },
        "model.DisplayAPIKey": {
            "type": "object",
            "properties": {
                "authentication_scopes": {
                    "type": "string"
                },
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "integer"
                },
                "expires_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "model.LicenseRequest": {
            "type": "object",
            "required": [
                "key"
            ],
            "properties": {
                "key": {
                    "type": "string"
                },
                "machine": {
                    "type": "string"
                }
            }
        },
        "model.LicenseRevokeRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "model.Product": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "integer"
                },
                "features": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ProductFeature"
                    }
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "model.ProductFeature": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "product_id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "model.ProductKeyPair": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "private_key": {
                    "type": "string"
                },
                "product_id": {
                    "type": "string"
                },
                "public_key": {
                    "type": "string"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "model.UserWithScopes": {
            "type": "object",
            "properties": {
                "authentication_scopes": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "user": {
                    "$ref": "#/definitions/model.User"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKey": {
            "type": "apiKey",
            "name": "X-API-KEY",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http"},
	Title:            "License Server API",
	Description:      "REST API for managing software licenses",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
