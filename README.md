# Simple http Go server

This is a simple HttpServer in Go

This server is being built as I progress through [boot.dev](https://boot.dev)'s
course [Learn HTTP Servers in Go](https://www.boot.dev/courses/learn-http-servers-golang).

## Functionality

So far we start a server on port _8080_
now we serve the files in the folder we run the server from /app path

We now have an API /api with the method:

/api/healthz  

available endpoints

### GET /api/healthz

allways returns 200/OK when the server is runing

### POST /api/users

allows the creation of users, requieres 2 fields email and password:

```
{
    "email": "<string>",    // unique user email
    "password": "<string>"  // user's password
}
```

on successfull register returns a User object wich contains all not private
user information:

```
{
    "id": "<UUID_string>",      // unique identifier string formatted UUID
    "created_at": "<UTC_date>", // original creation of the user account
    "updated_at": "<UTC_date>", // last time the user profile was modified
    "email": "<string>"         // user's unique email
}
```

### POST /api/login

allows registered users to get back a bearer token for use in other api calls,
requieres:

```
{
    "email": "<string>",        // unique user email
    "password": "<string>",     // user's password
    "expires_in_seconds": int   // optional defaults to 3600 if absent,
                                // 0 <= num <= 3600 s expiration time of the token
}
```

on succesful validation of the user credentials returns:

```
{
    "id": "<UUID_string>",      // unique identifier string formatted UUID
    "created_at": "<UTC_date>", // original creation of the user account
    "updated_at": "<UTC_date>", // last time the user profile was modified
    "email": "<string>",        // user's unique email
    "token": "<string>"         // Bearer token
}
```

### POST /api/chirps

allows registered users to publish a chirp to the server, requieres atentication:

`Authorization` header must be set with `Bearer <bearer_token>` to be able
to publish a chirp

```
{
    "body":"<text>" // text of the chirp to be published, up to 140 characters
                    // some words may be censored on the published chirp
}
```

on succes returns code 201 and the body:

```
{
    "id": "<UUID_string>",      // unique identifier of the chirp string-formatted UUID
    "created_at": "<UTC_date>", // original creation of the chirp
    "updated_at": "<UTC_date>", // last time the chirp was modified
    "body": "<string>",        // published text of the chirp
    "user_id": "<string>"       // unique identifier of the user who created the chirp string-formatted UUID
}
```

### GET /api/chirps

allows acces to all published chirps, requieres no parameters, returns a list
of all published chirps so far:

```
[
    {
        "id": "<UUID_string>",      // unique identifier of the chirp string-formatted UUID
        "created_at": "<UTC_date>", // original creation of the chirp
        "updated_at": "<UTC_date>", // last time the chirp was modified
        "body": "<string>",        // published text of the chirp
        "user_id": "<string>"       // unique identifier of the user who created the chirp string-formatted UUID
    },
    ...
]
```

### GET /api/chirps/{chirpID}

allows acces to a specific chirp, requieres the chir's unique id as the path
`/{chirpID}` parameter, if found returns the chirp fields:

```
{
    "id": "<UUID_string>",      // unique identifier of the chirp string-formatted UUID
    "created_at": "<UTC_date>", // original creation of the chirp
    "updated_at": "<UTC_date>", // last time the chirp was modified
    "body": "<string>",        // published text of the chirp
    "user_id": "<string>"       // unique identifier of the user who created the chirp string-formatted UUID
}
```

### GET /admin/metrics

allows acces to counters for the number of times that selected endpoints
have been requested.


### POST /admin/reset

allows the reset of the countes in metrics, and on _dev_ mode purgues the
entire user database.
