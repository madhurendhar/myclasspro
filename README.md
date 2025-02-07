# goscraper

This project is a simple Golang API that fetches HTML content from a specified URL and parses it dynamically using regex. The API is designed for performance, speed, and accuracy, accommodating changes in the HTML structure.

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone https://github.com/rahuletto/goscraper
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