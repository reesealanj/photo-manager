# Photo Manager

Photo Manager is a command-line utility written in Go that helps you back up files from a source directory to a destination directory. It organizes your files into folders based on their modification dates, making it easier to manage and access your backups. The tool supports various options like dry runs, debug mode, and can create destination directories if they don't exist.

## Features

- **Date-based Organization**: Files are organized into folders by year, month, and day based on their modification dates.
- **Dry Run Mode**: Simulate the backup process without making any actual changes.
- **Debug Mode**: Provides detailed logs for troubleshooting.
- **Automatic Directory Creation**: Creates destination directories if they don't exist.
- **Skip Existing Files**: Option to skip files that already exist in the destination.

## Getting Started

### Prerequisites

- **Go**: You need to have Go installed on your system to build the project. You can download it from [golang.org/dl](https://golang.org/dl/).

### Installation

Clone the repository and navigate to the project directory:

```bash
git clone { this repo } 
cd photo-manager
```

Build the binary using the provided Makefile:

```bash
make build-dev
```

This will create a binary named `photo-manager` in the `bin` directory.

## Usage

### Syntax

```bash
./bin/photo-manager [options]
```

### Options

- `--source=PATH`: Specify the source directory to back up. Default is `./source`.
- `--dest=PATH`: Specify the destination directory where backups will be stored. Default is `./dest`.
- `--dryRun`: Perform a dry run without copying any files.
- `--debug` or `-d`: Enable debug mode to get detailed logs.
- `--create-dest` or `-c`: Create the destination directory if it doesn't exist.
- `--allow-skip` or `--skip`: Skip files that already exist in the destination.

### Examples

#### Basic Backup

```bash
./bin/photo-manager --source="/path/to/source" --dest="/path/to/destination"
```

#### Dry Run

Simulate the backup process:

```bash
./bin/photo-manager --source="/path/to/source" --dest="/path/to/destination" --dryRun
```

#### Enable Debug Mode

Get detailed logs for troubleshooting:

```bash
./bin/photo-manager --source="/path/to/source" --dest="/path/to/destination" --debug
```

#### Create Destination Directory

Automatically create the destination directory if it doesn't exist:

```bash
./bin/photo-manager --source="/path/to/source" --dest="/path/to/destination" --create-dest
```

#### Skip Existing Files

Skip files that already exist in the destination:

```bash
./bin/photo-manager --source="/path/to/source" --dest="/path/to/destination" --allow-skip
```

## How It Works

1. **Argument Parsing**: The tool parses command-line arguments to configure its behavior.
2. **Validation**: Validates the source and destination directories.
   - If `--create-dest` is provided, it will create the destination directory if it doesn't exist.
3. **File Gathering**: Recursively scans the source directory and gathers all files.
4. **Date Parsing**: For each file, it extracts the modification date to determine the folder structure.
5. **Directory Preparation**: Prepares the destination directories based on the year, month, and day.
6. **File Copying**: Copies files from the source to the destination, maintaining the directory structure.
   - If `--dryRun` is enabled, it will only simulate this process.
   - If `--allow-skip` is enabled, it will skip copying files that already exist in the destination.

## Development

### Building the Project

You can build the project for development or production using the Makefile.

#### Build for Development

```bash
make build-dev
```

#### Build for Production (All Platforms)

```bash
make build-prod
```

This will create binaries for Linux, macOS, and Windows in the `bin` directory.

### Running the Application

#### Run in Development Mode

```bash
make run MODE=dev
```

#### Run in Production Mode

Simply run:

```bash
make run MODE=prod
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.