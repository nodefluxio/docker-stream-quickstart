#### Create Enrollement

```json
endpoint : /api/enrollment
method : POST
content-type :  multipart/form-data
request body :
  - name : testing
  - images : (multiple binary file)

response :
- 200
  {
    "ok": true,
    "message": "Successfully enroll new person",
    "enrollment":{
      "id":1,
      "name": "bambang",
      "face_id": 106596746326441980,
      "created_at": "2021-01-05T09:03:32.014123Z",
      "updated_at": "2021-01-05T09:03:32.014123Z",
      "deleted_at": null,
      "faces":[
        {
          "id": 1,
          "variation": 1234567876756,
          "image": "base64"
        }
      ]
    }
  }
```

#### Create Detail Enrollement

```json
endpoint : /api/enrollment/:id_enrollment
method : GET
content-type : application/json

response :
- 200
  {
    "ok": true,
    "message": "Succesfully get detail enrolled faces",
    "enrollment":{
      "id":1,
      "name": "bambang",
      "face_id": 106596746326441980,
      "created_at": "2021-01-05T09:03:32.014123Z",
      "updated_at": "2021-01-05T09:03:32.014123Z",
      "deleted_at": null,
      "faces":[
        {
          "id": 1,
          "variation": 1234567876756,
          "image": "base64"
        }
      ]
    }
  }
```

#### Get Enrollement

```json
endpoint: /api/enrollment
method : GET
query param :
 - limit : 10
 - page : 1
 - sort : timestamp:desc
content-type : application/json
response :
- 200
  {
    "ok": true,
    "message": "succesfully get enrolled person",
    "results": {
      "limit": 10,
      "current_page": 1,
      "total_data": 1,
      "total_page": 1,
      "enrollments":[
        {
          "id":1,
          "name": "bambang",
          "face_id": 106596746326441980,
          "created_at": "2021-01-05T09:03:32.014123Z",
          "updated_at": "2021-01-05T09:03:32.014123Z",
          "deleted_at": null,
          "faces":[
            {
              "id": 1,
              "variation": 1234567876756,
              "image": "base64"
            }
          ]
        }
      ]
    }
  }

```

#### PUT Enrollement

```json
endpoint: /api/enrollment/:id_enrollment
method : PUT
content-type : multipart/form-data
request body :
  - id : [1,2,3]
  - images : (multiple binary file)
  - name : "bambang"*
response :
- 200
  {
    "ok": true,
    "message": "succesfully update enrolled person",
    "enrollment":{
      "name": "bambang"
    }
  }

```

#### Delete Enrollement

```json
endpoint: /api/enrollment/:id_enrollment
method : DELETE
content-type : application/json
response :
- 200
  {
    "ok": true,
    "message": "succesfully enrolled person",
  }

```
