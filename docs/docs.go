// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "api@utile.space"
        },
        "license": {
            "name": "utile.space API License",
            "url": "https://utile.space/api/"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/dns/aaaa/{domain}": {
            "get": {
                "description": "Resolves AAAA records (IPv6) of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "AAAA resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/caa/{domain}": {
            "get": {
                "description": "Resolves CAA records of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "CAA resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/cname/{domain}": {
            "get": {
                "description": "Resolves CNAME records of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "CNAME resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/dmarc/{domain}": {
            "get": {
                "description": "Resolves DMARC TXT records of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "DMARC resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/mx/{domain}": {
            "get": {
                "description": "Resolves MX records of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "MX resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/ns/{domain}": {
            "get": {
                "description": "Resolves the name servers of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "NS resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/ptr/{ip}": {
            "get": {
                "description": "Resolves a domain name for a given IP address",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "PTR resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "IP address",
                        "name": "ip",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/txt/{domain}": {
            "get": {
                "description": "Resolves TXT records of a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "TXT resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/dns/{domain}": {
            "get": {
                "description": "Resolves a given domain name",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dns"
                ],
                "summary": "DNS resolution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Domain to resolve",
                        "name": "domain",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DNSResolution"
                        }
                    }
                }
            }
        },
        "/d{dice}": {
            "get": {
                "description": "Endpoint to roll a dice of the given number of faces",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "dice"
                ],
                "summary": "Roll a dice",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of faces of the dice between 2 and 100",
                        "name": "dice",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.DieResult"
                        }
                    }
                }
            }
        },
        "/links": {
            "get": {
                "description": "Returns a page of recommended links by SonnyAD",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "links"
                ],
                "summary": "Get Recommended Links Page",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Start cursor for pagination",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search filter",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.LinksPage"
                        }
                    }
                }
            }
        },
        "/status": {
            "get": {
                "description": "Get the status of the API",
                "produces": [
                    "application/json",
                    "text/xml",
                    "application/yaml",
                    "text/plain"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Healthcheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Health"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.DNSResolution": {
            "type": "object",
            "properties": {
                "resolution": {},
                "type": {
                    "type": "string"
                }
            }
        },
        "api.DieResult": {
            "type": "object",
            "properties": {
                "die": {
                    "type": "integer"
                },
                "result": {
                    "type": "integer"
                }
            }
        },
        "api.Link": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Tag"
                    }
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "api.LinksPage": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Link"
                    }
                },
                "next": {
                    "type": "string"
                }
            }
        },
        "api.Tag": {
            "type": "object",
            "properties": {
                "color": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.Health": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "utile.space Open API",
	Description:      "The collection of free API from utile.space, the Swiss Army Knife webtool.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
