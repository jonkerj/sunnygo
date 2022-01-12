# Webconnect
This document describes my efforts in revese engineering the SMA webconnect
protocl. It started when I got a SMA Sunnyboy inverter, which desperately needed
to be scraped into InfluxDB.

## Basic operation
The inverter is accessed through HTTPS (probably with an invalid certificate).
You authenticate, get a session ID, which is used later on in subsequent
requests. In the end, you'll want to sign out, as the session table appears to
be limited.

### Requests
Every request (including authentication), is a HTTP POST to some URI, with a
JSON object as a body. People call this REST nowadays. You'll have to set a
header `Content-Type: application/json; charset=utf-8` as well.

### Authentication
Request to `/dyn/login.json` with this body:
```json
{
  "right": ROLE,
  "pass": PASSWORD
}
```
There are, AFAIK, two roles: `usr` (User) and `istl` (Installer). If login
succeeded, you will get a JSON response containing a session ID:

```json
{
  "result": {
  	"sid": "Shgjsjh_287"
  }
}
```

### Data requests
After authentication, you'll need to pass the session ID as a query parameter
(`sid=Shgjsjh_287`).

## Instantaneous values
For `sunnygo`, the instantaneous values are the most interesting. They can be
retrieved using the URI `/dyn/getAllOnValues.json`. The web UI issues the
following body, which can be left out:
```json
{
  "destDev": []
}
```

The response is a seemingly strange data structure, which can later be explained
by the actual data model and language file. I left a lot out for brevity.

```json
{
  "result": {
    "0123-12345678": {
      "6100_40263F00": {
        "1": [
          {
            "val": 2169
          }
        ]
      },
      "6100_40265F00": {
        "1": []
      },
      "6100_00652C00": {
        "1": [
          {
            "val": null
          },
          {
            "val": null
          },
          {
            "val": null
          }
        ]
      },
      "6180_104A9B00": {
        "1": [
          {
            "val": "255.255.255.0"
          }
        ]
      },
      "6180_08652500": {
        "1": [
          {
            "val": [
              {
                "tag": 303
              }
            ]
          }
        ]
      }
    }
  }
}
```

The first level in the dict distinguishes devices. You'll probably get results
of a single device, in this case identified by `0123-12345678`. Let's stick to
what I call "fields" for now: `6100_40263F00`, `6100_40265F00`, `6100_00652C00`,
`6180_104A9B00` and `6180_08652500`.

These fields have values embedded in them. The key is to interpret the data,
which can be done using the next chapter.

## Data model
The data model can be retrieved without authentication using HTTP GET to
`/data/ObjectMetadata_Istl.json` (probably also `_Usr.json`, but it will contain
less). The result is a JSON representation of each possible field known to this
firmware. There are probably more fields in the model as there are in your
inverter.

The following is an excerpt of the data model, only containing above fields:

```json
{
  "6100_40263F00": {
    "Prio": 1,
    "TagId": 416,
    "TagIdEvtMsg": 10030,
    "Unit": 18,
    "DataFrmt": 0,
    "Scale": 1,
    "Typ": 0,
    "WriteLevel": 5,
    "TagHier": [
      835,
      230
    ],
    "Min": true,
    "Max": true,
    "Sum": true,
    "Avg": true,
    "Cnt": true,
    "SumD": true
  },
  "6100_40265F00": {
    "Prio": 1,
    "TagId": 413,
    "TagIdEvtMsg": 10043,
    "Unit": 16,
    "DataFrmt": 0,
    "Scale": 1,
    "Typ": 0,
    "WriteLevel": 5,
    "TagHier": [
      835,
      230
    ],
    "Min": true,
    "Max": true,
    "Sum": true,
    "Avg": true,
    "Cnt": true,
    "SumD": true
  },
  "6100_00652C00": {
    "Prio": 3,
    "TagId": 3314,
    "TagIdEvtMsg": 12025,
    "Unit": 1,
    "DataFrmt": 0,
    "Scale": 1,
    "Typ": 0,
    "WriteLevel": 5,
    "TagHier": [
      834,
      67,
      4067
    ],
    "Min": true,
    "Max": true,
    "Avg": true,
    "Cnt": true,
    "Hidden": true,
    "MinD": true,
    "MaxD": true
  },
  "6180_104A9B00": {
    "Prio": 2,
    "TagId": 1713,
    "TagIdEvtMsg": 10883,
    "DataFrmt": 8,
    "Typ": 2,
    "WriteLevel": 5,
    "TagHier": [
      839,
      1708
    ]
  },
  "6180_08652500": {
    "Prio": 3,
    "TagId": 240,
    "TagIdEvtMsg": 12027,
    "DataFrmt": 18,
    "Typ": 1,
    "WriteLevel": 5,
    "TagHier": [
      834,
      309,
      4065
    ],
    "Hidden": true
  }
}
```

This is what I know about this structure:
- `TagId` is the ID of the name of the field in the language table
- `TagIdEvtMsg` the same, but for the corresponding event message
- `DataFrmt`: Data format. By looking at the de-obfuscated Web UI code, I found
  the following values in my inverter:
  - 0-3 fixed point (0-3 decimals)
  - 5: time, seconds after midnight
  - 6: datetime as string
  - 7: duration
  - 8: UTF-8
  - 18: taglist
  - 26: fixed (4 decimals)
- `Typ`: Data type. From the Web UI sources:
  - 0: scalar
  - 1: status
  - 2: text
- `TagHier`: Tag hierarchy. This described where in the hierarchy a field is
  organized. The numbers correspond to IDs in the language table.
- `Scale`: for scalars, this is the scale to multiply the value with.
- `Unit`: also for scalars, the unit. Look it up in the language table.

## Language
Same as object model, HTTP GET to `/data/l10n/en-US.json`. It is a simple map,
from IDs to strings:

```json
{
  "1": "%",
  "16": "var",
  "18": "W",
  "67": "DC measurements",
  "230": "Grid measurements",
  "240": "Condition",
  "303": "Off",
  "309": "Operation",
  "413": "Reactive power",
  "416": "Power",
  "834": "DC Side",
  "835": "AC Side",
  "839": "System communication",
  "1708": "Speedwire",
  "1713": "Current subnet mask",
  "3314": "Signal strength of the selected network",
  "4065": "PV module control",
  "4067": "PV module electronics",
  "10030": "Power",
  "10043": "Reactive power",
  "10883": "Current speedwire subnet mask",
  "12025": "Received signal strength",
  "12027": "Condition"
}
```

## Example
Combining all this, this is how you the example values should be interpreted:
- `6100_40263F00`:
  - Organization: "AC Side / Grid measurements"
  - Name: "Power"
  - Type: scalar
  - Data format: integer (fixed point with 0 decimals)
  - Value: 2169 * 1 W
- `6100_40265F00`:
  - Organization: "AC Side / Grid measurements"
  - Name: "Reactive power"
  - Type: scalar
  - Data format: integer (fixed point with 0 decimals)
  - Value: n/a
- `6180_104A9B00`:
  - Organization: "DC Side / DC measurements / PV module electronics"
  - Name: "Signal strength of the selected network"
  - Type: scalar
  - Data format: integer (fixed point with 0 decimals)
  - Value: n/a, n/a, n/a
- `6180_104A9B00`:
  - Organization: "System communication / Speedwire"
  - Name: "Current subnet mask"
  - Type: text
  - Data format: UTF-8
  - Value: "255.255.255.0"
- `6180_08652500`:
  - Organization: "DC Side / Operation / PV module control"
  - Name: "Condition"
  - Type: status
  - Data format: taglist 
  - Value: "Off"
