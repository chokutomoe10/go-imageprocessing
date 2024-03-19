# image-processing

I'm using OpenCV library for image processing tasks, make sure OpenCV be installed on your system. Follow the installation instructions based on your OS [here](https://gocv.io/getting-started/).

## To Setup and Start

Install GoCV package:

```go
go get -u gocv.io/x/gocv
```

Create `compressed`, `converted`, `resized` folders inside the `images` folder and create `PNG`, `JPEG` folders inside `converted` folder.

Run the service:

```go
go run main.go
```

## Test the HTTP Routes

To test the HTTP routes, you can use application called Postman or other applications that have similar functionalities with it to test the HTTP routes.

### Convert

Go to `Body`, `form-data`, add `Images` for `Key`, upload PNG images from your local machine in `Value` (you can use the images inside the uploads folder), and send the request to `localhost:8080/convert` url with `POST` method.

### Resize

Go to `Body`, `form-data`, add `Images`, `SizeX`, `SizeY` for `Key`, upload images from your local machine in `Value` (you can use the images inside the uploads folder), and send the request to `localhost:8080/resize` url with `POST` method.

### Compress

Go to `Body`, `form-data`, add `Images`, `SizeJPEG` (for JPEG images), `SizeX` and `SizeY` (for PNG images) for `Key`, upload images from your local machine in `Value` (you can use the images inside the uploads folder), and send the request to `localhost:8080/compress` url with `POST` method

## Running the tests

```go
# run all tests
go test -v

# run all tests with test coverage
go test -v -cover
````