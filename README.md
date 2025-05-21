# Go Error Handling Demo

This project demonstrates various techniques for handling errors gracefully in Go

It includes examples of:

- Basic error returning and checking.

- Creating errors with `errors.New` and `fmt.Errorf`.

- Error wrapping using `fmt.Errorf` with the `%w` verb (Go 1.13+).

- Inspecting error chains with `errors.Is()` for sentinel errors.

- Inspecting error chains with `errors.As()` for custom error types.

- Defining and using custom error types with an `Unwrap()` method.

The application is designed to be run within a Docker container.

## Project Structure

```
error-handling-demo/
├── main.go         # The Go application code
├── go.mod          # Go module definition
├── go.sum          # Go module checksums
└── Dockerfile      # Docker instructions to build the image

```

## Prerequisites

- **Docker:** Ensure Docker is installed and running on your system. Download from [docker.com](https://www.docker.com/get-started).

- **Go (Optional, for local development):** Go version 1.13 or later if you wish to run or modify the Go code directly without Docker. Download from [golang.org](https://golang.org/dl/).

## Docker Instructions

These instructions will guide you through building the Docker image and running the containerized application.

### 1. Clone the Repository (or Create Files)

If you have this project in a Git repository, clone it:

```
git clone https://github.com/munich-gophers/go-error-handling-demo
cd error-handling-demo

```

Otherwise, ensure you have the `main.go`, `go.mod`, and `Dockerfile` in a directory named `error-handling-demo`.

### 2. Prepare Go Modules (if necessary)

If you've just created the `go.mod` file or made changes, it's good practice to tidy the modules. This step is also handled within the Docker build if `go.sum` is missing, but doing it locally ensures your `go.sum` is up-to-date.

```
# Navigate to the error-handling-demo directory
go mod tidy

```

This will generate/update the `go.sum` file.

### 3. Build the Docker Image

Navigate to the `error-handling-demo` directory in your terminal (the directory containing the `Dockerfile`). Then, run the following command to build the Docker image:

```
docker build -t go-error-handling-demo .

```

- `-t go-error-handling-demo`: Tags the image with the name `go-error-handling-demo`. You can choose a different name/tag.

- `.`: Specifies that the build context (location of `Dockerfile` and source files) is the current directory.

### 4. Run the Docker Container

Once the image is built successfully, you can run the application in a Docker container:

```
docker run --rm go-error-handling-demo

```

- `--rm`: Automatically removes the container when it exits. This is useful for keeping your system clean after running short-lived demo applications.

- `go-error-handling-demo`: The name of the image you built in the previous step.

## Expected Output

Running the container will execute the `main` function in `main.go`, which demonstrates several error handling scenarios. The output will look similar to this:

```
--- Processing request with config 'missing.json' and data query '123' ---
Detailed Config Error: Operation 'open file' on file 'missing.json' failed.
  Underlying cause: File does not exist. Attempting fallback...
Main Error Handler: failed during config stage: config error during 'open file' for file 'missing.json': failed to open: file does not exist

--- Processing request with config 'invalid.json' and data query '123' ---
Detailed Config Error: Operation 'parse content' on file 'invalid.json' failed.
Main Error Handler: failed during config stage: config error during 'parse content' for file 'invalid.json': invalid JSON structure

--- Processing request with config 'valid.json' and data query 'notfound_db' ---
Successfully loaded config: valid.json
Data Fetch Error: The requested data was not found in the database (sql.ErrNoRows).
Main Error Handler: failed during data fetching stage: database query for ID 'notfound_db' failed: sql: no rows in result set

--- Processing request with config 'valid.json' and data query 'custom_resource_err' ---
Successfully loaded config: valid.json
Data Fetch Error: A specific custom resource was not found (ErrResourceNotFound).
Main Error Handler: failed during data fetching stage: could not retrieve resource data: resource not found

--- Processing request with config 'valid.json' and data query '123' ---
Successfully loaded config: valid.json
Successfully processed request. Data: Sample Data

```

## Running Locally (Without Docker - Optional)

If you have Go installed (version 1.13+), you can also run the application directly:

1. Navigate to the `error-handling-demo` directory.

2. Ensure dependencies are up-to-date:

   ```
   go mod tidy

   ```

3. Run the `main.go` file:

   ```
   go run main.go

   ```

This will produce the same output as running the Docker container.
