# Bible Notes Backend

A Go application designed to facilitate the management and querying of Bible data. This backend tool allows you to convert Bible texts from a JSON format into an efficient SQLite database and then query specific verses or ranges of verses from that database via the command line.

## Features

*   **JSON to SQLite Conversion:** Easily transform structured JSON Bible data into a portable SQLite database (`.db` file).
*   **Flexible Verse Querying:** Retrieve Bible verses using a simple notation for books, chapters, and verses (e.g., `John 3:16` or `Psalm 23:1-6`).

## Prerequisites

To build and run this application, you need:

*   **Go:** Version 1.25.0 or newer.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/ezek-iel/bible-notes-backend.git
    cd bible-notes-backend
    ```

2.  **Build the application:**
    ```bash
    go build -o bible-notes-backend .
    ```
    This command will create an executable named `bible-notes-backend` in your current directory.

## Usage

The application supports two primary subcommands: `convert` and `query`.

### 1. Convert Command

Use the `convert` command to transform a JSON file containing Bible data into an SQLite database.

**Syntax:**

```bash
./bible-notes-backend convert --source-file <input.json> --destination-file <output.db>
```

*   `--source-file`: Path to your input JSON file (e.g., `NLT.json`).
*   `--destination-file`: Path where the output SQLite database file will be created (e.g., `bible.db`).

**Example:**

To convert `NLT.json` into `bible.db`:

```bash
./bible-notes-backend convert --source-file NLT.json --destination-file bible.db
```

### 2. Query Command

Use the `query` command to retrieve specific Bible verses from an existing SQLite database.

**Syntax:**

```bash
./bible-notes-backend query --verse-notation "Book Chapter:Verse" --database-file <database.db>
```

*   `--verse-notation`: The Bible verse(s) you want to query. Supported formats:
    *   `Book Chapter:Verse` (e.g., `"John 1:14"`)
    *   `Book Chapter:StartVerse-EndVerse` (e.g., `"John 1:14-16"`)
*   `--database-file`: Path to your SQLite Bible database file (e.g., `bible.db`).

**Examples:**

To query a single verse (John 1:14):

```bash
./bible-notes-backend query --verse-notation "John 1:14" --database-file bible.db
```

To query a range of verses (John 1:14-16):

```bash
./bible-notes-backend query --verse-notation "John 1:14-16" --database-file bible.db
```

The output will display each verse in the format: `Book Chapter:Verse - Text`.

## Project Structure

*   `main.go`: The main entry point of the application, handling command-line parsing and dispatching to `convert` or `query` logic.
*   `converter/`: Contains the logic for converting JSON Bible data into an SQLite database.
*   `query/`: Contains the logic for querying verses from the SQLite database.
*   `utils/`: (If present) Utility functions used across the project.
*   `tests/`: (If present) Unit and integration tests.
*   `NLT.json`: An example JSON file containing Bible data.
*   `filename.db`: An example SQLite database file (might be generated or used as a placeholder).
*   `go.mod`, `go.sum`: Go module definition and dependency checksums.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

(Coming soon)