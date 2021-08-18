```json
method : GET
url : /api/events
query param :
 - limit : 10
 - page : 1
 - filter[type] : unrecognized / recognized
 - filter[timestamp_from] : 2020-07-27T00:00:00+07:00
 - filter[timestamp_to] : 2020-07-28T00:00:00+07:00
 - sort : timestamp:desc
 - search : bambang
response body :
{
  "ok": true,
  "message": "successfully fetch event data",
  "results": {
    "limit": 10,
    "current_page": 1,
    "total_data": 1,
    "total_page": 1,
    "events": [
      {
        "timestamp": "string",
        "data": [
          {
            "id": 0,
            "primary_image": "string",
            "secondary_image": "string",
            "label": "string",
            "result": "string",
            "location": "string",
            "timestamp": "string"
          }
        ]
      }
    ]
  }
}
```