package docs

import "github.com/swaggo/swag"

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate: docTemplate,
	})
}

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API для управления WhatsApp Web через Selenium",
        "title": "WhatsApp Parser API",
        "version": "1.0"
    },
    "host": "localhost:8081",
    "basePath": "/",
    "schemes": ["http"],
    "paths": {
        "/session": {
            "post": {
                "description": "Создает новую сессию и возвращает QR код для авторизации",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Создать новую сессию WhatsApp",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "session_id": {
                                    "type": "string"
                                },
                                "qr_code": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/session/{id}": {
            "post": {
                "description": "Восстанавливает сохраненную сессию WhatsApp по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Восстановить существующую сессию",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID сессии",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/session/{id}/message": {
            "post": {
                "description": "Отправляет сообщение через WhatsApp используя указанную сессию",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "message"
                ],
                "summary": "Отправить сообщение",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID сессии",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные сообщения",
                        "name": "message",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/SendMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "SendMessageRequest": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Hello, World!"
                },
                "phone_number": {
                    "type": "string",
                    "example": "1234567890"
                }
            }
        }
    }
}` 