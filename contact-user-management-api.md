### List User 
```json
endpoint : /api/manage-users
method : GET
query param :
 - limit : 10
 - page : 1
 - sort : timestamp:desc
content-type : application/json
response :
{
  "ok": true,
  "message": "successfully get list of user",
  "results": {
    "limit": 10,
    "current_page": 1,
    "total_data": 1,
    "total_page": 1,
    "users": [
      {
      "id": "4",
      "email": "test@test1.com",
      "username": "testtest1",
      "fullname": "test",
      "role": "operator",
      "avatar": null,
      "site_id": [1,2,3,4],
      "created_at": "2021-04-07T05:00:25.859Z",
      "updated_at": "2021-04-07T05:00:25.859Z",
      "deleted_at": null,
    }
    ]
  }
}
```

### Get Detail User 
```json
endpoint : /api/manage-users/:id
method : GET
content-type : application/json
response :
{
  "ok": true,
  "message": "Successfully get detail user",
  "user":{
    "id": "4",
    "email": "test@test1.com",
    "username": "testtest1",
    "fullname": "test",
    "role": "operator",
    "avatar": null,
    "site_id": [1,2,3,4],
    "created_at": "2021-04-07T05:00:25.859Z",
    "updated_at": "2021-04-07T05:00:25.859Z",
    "deleted_at": null,
  }
}
```

#### Add User

```json
endpoint : /api/manage-users
method : POST
content-type : application/json
request body :
{
  "email": "test@test1.com",
  "username": "testtest1",
  "fullname": "test",
  "role": "operator",
  "avatar": "base64",
  "site_id": [1,2,3,4],
  "password": "password",
  "re_password": "password",
}


response :
- 201
  {
    "ok": true,
    "message": "Successfully add new user",
    "user":{
      "id": "4",
      "email": "test@test1.com",
      "username": "testtest1",
      "fullname": "test",
      "role": "operator",
      "avatar": null,
      "site_id": [1,2,3,4],
      "created_at": "2021-04-07T05:00:25.859Z",
      "updated_at": "2021-04-07T05:00:25.859Z",
      "deleted_at": null,
    }
  }
```

#### update user info

```json
endpoint: /api/manage-users/:id
method : PUT
content-type : application/json
request body :
{
  "fullname": "test",
  "role": "operator",
  "avatar": "base64",
  "site_id": [1,2,3,4],
}
response :
- 200
  {
    "ok": true,
    "message": "succesfully update user data",
  }

```

#### update user info

```json
endpoint: /api/manage-users/:id/change-password
method : PUT
content-type : application/json
request body :
{
  "email": "test@test1.com",
  "username": "testtest1",
  "password": "password",
  "re_password": "password",

}
response :
- 200
  {
    "ok": true,
    "message": "succesfully change password user",
  }

```


#### delete user info

```json
endpoint: /api/manage-users/:id
method : DELETE
content-type : application/json
response :
- 200
  {
    "ok": true,
    "message": "succesfully delete user data",
  }

```