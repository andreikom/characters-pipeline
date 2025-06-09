# Multiverse Character Aggregator

This project develops a program that retrieves all available characters from three different fictional universes via
their public APIs:

- **PokéAPI**
- **Star Wars API (SWAPI)**
- **Rick and Morty API**

It **normalizes the data into a consistent format**, **stores it**, and **returns a single aggregated alphabetically (by name) sorted list** of
characters through a unified API.

## Normalized Character Schema

Each character in the aggregated response includes the following fields:

- **`name`**: Character’s name
- **`origin`**: Universe of origin (e.g., "Pokémon", "Star Wars", "Rick and Morty")
- **`species`**: Character’s species or type
- **`additional_attribute`**: A unique attribute from each universe:
    - **PokéAPI**: `base_experience` (or a similar numeric/stat field)
    - **SWAPI**: `birth_year` (e.g., "19BBY")
    - **Rick and Morty**: `status` (e.g., "Alive", "Dead", "Unknown")

## Requirements

- Golang must be installed

Install Go on macOS via Homebrew

```bash
brew install go
```

## Build & Test Commands

This project uses a shell script (build.sh) instead of a Makefile.
Available Commands

- `./build.sh build`        # Build the project (default target: macOS) and places it under a generated `./build` dir
- `./build.sh test`         # Run all Go tests in the project
- `./build.sh clean`        # Remove build artifacts

## RUN

Once the service binary is built - execute via: 

```bash
./build/characters-pipeline-darwin-amd64
```

The service will print a temporarily filestore path for the local storage of each persisted list, e.g.:

```
2025/06/09 08:37:13 Temp filestore dir: /var/folders/47/jhfh4qws1t79bspw6vg6r5mm0000gn/T/charactersync_store_820753468
```

## INTERACT

Once the service is running , fetch the unified list of characters with:

```bash
curl localhost:8080/characters
```

Notes:
- Currently port 8080 is hardcoded
- A complete run may take several minutes to complete

### Optional Query Parameters

    limitPerProvider=N: Limits the number of characters fetched per provider. The total number of characters returned will be:

N * (number of registered providers)

Example: Get only 5 characters per API provider:

```bash
curl "http://localhost:8080/characters?limitPerProvider=5"
```

