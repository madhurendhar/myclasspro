# goscraper

This project is a simple Golang API that fetches HTML content from a specified URL and parses it dynamically using regex. The API is designed for performance, speed, and accuracy, accommodating changes in the HTML structure.

> [!TIP]
> GoScraper is now open-sourced and self-deployable. Run your own goscraper instance
>
>  - `SUPABASE_URL` and `SUPABASE_KEY` should be the same creds as the front-end env variables used to access Supabase
>  - `ENCRYPTION_KEY` is an unique hash that is used to encrypt your data while saving the data to supabase.
> `openssl rand --hex 32` run this on shell to get a random hex key used for encryption
> - `VALIDATION_KEY` is used to validate the request from front-end, so use the same key as front-end has. a different key will reject requests from your front-end
>  - `URL` are the urls which the backend should allow to request. **(CORS)**

### `.env`
```
SUPABASE_URL=""
SUPABASE_KEY=""
ENCRYPTION_KEY=""
VALIDATION_KEY=""
URL=""
```

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

3. **Development Run the application:** (DEV SERVER)
   ```
   go run main.go
   ```

3. **Build and Run the application:** (BUILD SERVER)
   ```
   go build main.go
   ./main
   ```

## Docker file
```
docker compose up --build
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.