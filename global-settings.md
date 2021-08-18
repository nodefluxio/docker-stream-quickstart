#### Create Or Update Settings (Similarity)

```json
endpoint : /api/settings
method : POST
content-type : application/jsons
request body :
  {
    "similarity":0.5
  }

response :
- 200
  {
    "global_setting": {
        "id": 3,
        "similarity": 0.5,
        "created_at": "2021-03-30T09:03:35.490362Z"
    },
    "message": "Successfully Create or Update result",
    "ok": true
  }
```

#### Get Settings (Similarity)

```json
endpoint : /api/settings
method : GET
content-type : application/json

response :
- 200
  {
    "message": "Succesfully get global setting",
    "ok": true,
    "setting": {
        "id": 0,
        "similarity": 0.7,
        "created_at": "0001-01-01T00:00:00Z"
  }
```
