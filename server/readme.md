## Running with Docker
Run the game using a previously built docker image with the following command.
```
$ docker run -p 8080:8080 guessing-game
```

## Building with Docker
Build the game and create a docker image using the following command.
```
$ docker build . -t guessing-game
```

## HTTP server endpoints

### Create a new game
To create a new game, use the `new` endpoint. This endpoint returns the game ID which can be used with other requests to reference a particular game.
```
GET /new
```
Example
```
$ curl localhost:8080/new

fcb24a31-7e6f-49fe-a3e1-42760aeb5410
```

### Make a guess
Make a guess using the `guess` endpoint. This endpoint submits a guess for the specified game and returns the result.
```
POST /guess?id=[uuid]
```
Body data
```
[integer]
```
Example
```
$ curl --data "50" localhost:8080/guess?id=fcb24a31-7e6f-49fe-a3e1-42760aeb5410

too high
```

### Game status
Retrieve the same status using the `isactive` endpoint. When a game is ongoing, the endpoint will return `true`. Once the number has been successfully guessed, the endpoint will return `false`.
```
GET /isactive?id=[uuid]
```
Example
```
$ curl localhost:8080/isactive?id=fcb24a31-7e6f-49fe-a3e1-42760aeb5410

true
```