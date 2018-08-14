# Measurement Tool

A server for recording measurements

### Interface

`GET /time`

Returns the current unix time on the server.

`POST /measurements`

Saves a list of measurements. Body should be a json formatted string:

```
[
    {
        "unixtime": 1234567890.1234567,
        "properties": [
            {
                "key": "temperature",
                "value": "15"
            },
            {
                "key": "location",
                "value": "livingroom"
            },
            {
                "key": "must always be a string",
                "value": "must always be a string"
            },
            {
                ... more properties ...
            }
        ]
    },
    {
        ... more measurements ...
    }
]
```
