# Compression Tool

## Overview

This is a Huffman compression and decompression tool implemented in Go. It provides a simple command-line interface to compress and decompress files using Huffman encoding.

## Features

- Compress text files using Huffman encoding.
- Decompress files back to their original text format.
- Handles UTF-8 encoded input.
- Uses cli flags for input/output files and compress/decompress. 

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/nobletk/compression_tool.git
    cd compression_tool
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

## Usage

### Build

To build the project, use the Makefile:
```bash
make build
```
### Compress file 
To compress a text file:
```bash
go run ./cmd/app -i= filepath/input_file.txt -o=output_file.txt -c
```
### Decompress file 
To decompress a text file:
```bash
go run ./cmd/app -i= filepath/input_file.txt -o=output_file.txt -d
```

### Testing
```bash
make test
```
