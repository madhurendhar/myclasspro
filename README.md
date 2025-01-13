# goscraper

This project is a simple Golang API that fetches HTML content from a specified URL and parses it dynamically using regex. The API is designed for performance, speed, and accuracy, accommodating changes in the HTML structure.

## Project Structure

```
goscraper
├── src
│   ├── main.go          # Entry point of the application
│   ├── handlers
│   │   └── fetch.go     # Handles fetching HTML content
│   ├── services
│   │   └── parser.go     # Contains HTML parsing logic
│   ├── models
│   │   └── data.go       # Defines data structures for parsed data
│   └── utils
│       └── regex.go      # Utility functions for regex operations
├── go.mod                # Module definition file
├── go.sum                # Checksums for module dependencies
└── README.md             # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone goscraper
   cd goscraper
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   reflex -r '\.go' -s -- sh -c 'clear; go run src/main.go'
   ```

## Docker file
```
docker compose up --build
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.