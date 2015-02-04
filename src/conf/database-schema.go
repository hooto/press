package conf

const dsBase = `
{
    "engine": "MyISAM",
    "charset": "utf8",
    "tables": [
        {
            "name": "spec",
            "columns": [
                {
                    "name": "id",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "state",
                    "type": "int16"
                },
                {
                    "name": "userid",
                    "type": "string",
                    "length": "10"
                },
                {
                    "name": "title",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "comment",
                    "type": "string",
                    "length": "200"
                },
                {
                    "name": "version",
                    "type": "string",
                    "length": "10",
                    "default": "1"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "state",
                    "type": 1,
                    "cols": ["state"]
                },
                {
                    "name": "userid",
                    "type": 1,
                    "cols": ["userid"]
                },
                {
                    "name": "created",
                    "type": 1,
                    "cols": ["created"]
                },
                {
                    "name": "updated",
                    "type": 1,
                    "cols": ["updated"]
                }
            ]
        },
        {
            "name": "nodex",
            "columns": [
                {
                    "name": "id",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "specid",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "name",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "state",
                    "type": "int16",
                    "default": "0"
                },
                {
                    "name": "userid",
                    "type": "string",
                    "length": "10"
                },
                {
                    "name": "title",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "comment",
                    "type": "string",
                    "length": "200"
                },
                {
                    "name": "fields",
                    "type": "string-text"
                },
                {
                    "name": "terms",
                    "type": "string-text"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "specid",
                    "type": 1,
                    "cols": ["specid"]
                },
                {
                    "name": "state",
                    "type": 1,
                    "cols": ["state"]
                },
                {
                    "name": "userid",
                    "type": 1,
                    "cols": ["userid"]
                },
                {
                    "name": "created",
                    "type": 1,
                    "cols": ["created"]
                },
                {
                    "name": "updated",
                    "type": 1,
                    "cols": ["updated"]
                }
            ]
        },
        {
            "name": "termx",
            "columns": [
                {
                    "name": "id",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "specid",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "name",
                    "type": "string",
                    "length": "30"
                },
                {
                    "name": "state",
                    "type": "int16",
                    "default": "0"
                },
                {
                    "name": "userid",
                    "type": "string",
                    "length": "10"
                },
                {
                    "name": "type",
                    "type": "string",
                    "length": "10"
                },
                {
                    "name": "title",
                    "type": "string",
                    "length": "100"
                },
                {
                    "name": "created",
                    "type": "datetime"
                },
                {
                    "name": "updated",
                    "type": "datetime"
                }
            ],
            "indexes": [
                {
                    "name": "PRIMARY",
                    "type": 3,
                    "cols": ["id"]
                },
                {
                    "name": "specid",
                    "type": 1,
                    "cols": ["specid"]
                },
                {
                    "name": "state",
                    "type": 1,
                    "cols": ["state"]
                },
                {
                    "name": "userid",
                    "type": 1,
                    "cols": ["userid"]
                },
                {
                    "name": "created",
                    "type": 1,
                    "cols": ["created"]
                },
                {
                    "name": "updated",
                    "type": 1,
                    "cols": ["updated"]
                }
            ]
        }
    ]
}
`

const (
	dsTplNodeModels = `
{
    "columns": [
        {
            "name": "id",
            "type": "string",
            "length": "30"
        },
        {
            "name": "state",
            "type": "int16"
        },
        {
            "name": "userid",
            "type": "string",
            "length": "10"
        },
        {
            "name": "title",
            "type": "string",
            "length": "100"
        },
        {
            "name": "created",
            "type": "datetime"
        },
        {
            "name": "updated",
            "type": "datetime"
        }
    ],
    "indexes": [
        {
            "name": "PRIMARY",
            "type": 3,
            "cols": ["id"]
        },
        {
            "name": "state",
            "type": 1,
            "cols": ["state"]
        },
        {
            "name": "userid",
            "type": 1,
            "cols": ["userid"]
        },
        {
            "name": "created",
            "type": 1,
            "cols": ["created"]
        },
        {
            "name": "updated",
            "type": 1,
            "cols": ["updated"]
        }
    ]
}
`
	dsTplTermModels = `
{
    "name": "template",
    "columns": [
        {
            "name": "id",
            "type": "uint32",
            "IncrAble": true
        },
        {
            "name": "state",
            "type": "int16"
        },
        {
            "name": "userid",
            "type": "string",
            "length": "10"
        },
        {
            "name": "title",
            "type": "string",
            "length": "100"
        },
        {
            "name": "created",
            "type": "datetime"
        },
        {
            "name": "updated",
            "type": "datetime"
        }
    ],
    "indexes": [
        {
            "name": "PRIMARY",
            "type": 3,
            "cols": ["id"]
        },
        {
            "name": "state",
            "type": 1,
            "cols": ["state"]
        },
        {
            "name": "userid",
            "type": 1,
            "cols": ["userid"]
        },
        {
            "name": "created",
            "type": 1,
            "cols": ["created"]
        },
        {
            "name": "updated",
            "type": 1,
            "cols": ["updated"]
        }
    ]
}
`
)
