# 🔒 Hyde: Ephemeral & Secure File Encryption

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) 

**Hyde is a simple, privacy-focused tool for encrypting and decrypting files using strong AES-256 encryption. It offers both a command-line interface (CLI) and a web application, designed with security, ephemerality, and ease-of-use in mind.**

## Why I created Hyde

I created Hyde for my personal use so I could encrypt files quickly and easily. I wanted a tool that was simple, secure, and didn't require me to remember complex commands or configurations. I also wanted to ensure that my files were not stored on any server after the operation was complete.

## Core Principles:

*   🛡️ **Secure:** Employs robust **AES-256** encryption, an industry standard for protecting sensitive data.
*   💨 **Ephemeral & Stateless:** Designed *never* to store your files or passwords server-side post-operation.
*   🔒 **Privacy-Focused:** No tracking, no logs (by default design), no unnecessary data collection. Your data remains yours. Intended primarily for self-hosting.
*   ✅ **Accessible & Simple:** Provides both a user-friendly web interface for quick tasks and a powerful CLI for scripting and local operations. No complex configurations needed.
*  0️⃣ **Zero Configuration:** No setup or configuration required. Just run the container and start encrypting/decrypting files.

## ✨ Features

*   **🖥️ Web UI:** Intuitive browser-based interface for easy encryption/decryption.
*   **⌨️ CLI Tool:** Versatile command-line tool for local file operations, ideal for scripting and automation. Includes bash completion.
*   **🔗 HTTP API:** Simple API used by the web app, potentially usable for custom integrations (though primarily internal).
*   **🔑 Strong Encryption:** Files protected with AES-256.
*   **📦 All-in-One Container:** Web UI, API, and CLI bundled together in a single Docker image for convenience.
*   **🌐 Cross-Platform:**
    *   **Docker (Recommended):** Run consistently across Linux, macOS, and Windows.
    *   **Binaries:** Native CLI binaries available for [Linux, macOS, and Windows](https://github.com/edsonjaramillo/hyde/releases).
*   **🚫 Stateless by Design:** Ensures no sensitive data remains on the server after processing.

## 🚀 Installation

The recommended way to run Hyde is via Docker. This ensures a consistent, isolated, and ephemeral environment.

**Prerequisites:**

*   [Docker](https://docs.docker.com/get-docker/) installed.
*   [Docker Compose](https://docs.docker.com/compose/install/) installed (optional, simplifies deployment).

**1. Pull the Docker Image:**

```bash
docker pull ghcr.io/edsonjaramillo/hyde:0.1.0
```

2. Run Hyde:

You can run Hyde using docker run or docker-compose.

**Option A:** docker run (Quick Start)
```bash
# Runs Hyde in the background, exposing Web UI (3000) and API (8080)
docker run -d --name hyde-playground -p 3000:3000 -p 8080:8080 ghcr.io/edsonjaramillo/hyde:0.1.0
```

**Option B:** docker run (With Local File Access for CLI)

Use a volume mount (-v) to map a local directory into the container. This allows the CLI inside the container to access files outside the container.
```bash
# Maps the current directory's 'hyde-files' subfolder to '/data' inside the container
# Create the 'hyde-files' folder first if it doesn't exist: mkdir hyde-files
docker run -d --name hyde-playground -p 3000:3000 -p 8080:8080 \
  -v "$(pwd)/hyde-files:/data" \
  ghcr.io/edsonjaramillo/hyde:0.1.0
```

**Option C:** docker-compose (Recommended for Configuration)
```yaml
version: '3.8'
services:
  hyde:
    image: ghcr.io/edsonjaramillo/hyde:0.1.0
    container_name: hyde-playground
    ports:
      - "3000:3000" # Web UI Port (Host:Container)
      - "8080:8080" # API Port (Host:Container)
    volumes:
      # Mount local './hyde-files' directory to '/data' inside the container
      - ./hyde-files:/data
    # Optional: Override default environment variables if needed
    # environment:
    #   - WEB_PORT=3000
    #   - API_PORT=8080
    #   - ALLOWED_ORIGINS=http://localhost:3000 # Adjust if accessing from different origin
    restart: unless-stopped # Optional: Keep the container running
```

## 🖥️ Playground (Web UI)

**Open your web browser:** 

Navigate to `http://localhost:3000` (or the port you specified in the Docker run command).

**Go to `/encrypt` page**
- Upload a file using the file input.
- Enter a password for encryption.
- Click the "Encrypt" button.
- Pop-up will appear with the encrypted file. You can download it directly from there.

**Go to `/decrypt` page**
- Upload the encrypted file using the file input.
- Enter the password used for encryption.
- Click the "Decrypt" button.
- Pop-up will appear with the decrypted file. You can download it directly from there.

## ⌨️ Guided Playground (CLI)

**Open a shell session**:
```bash
docker exec -it -u hyde hyde-playground /bin/bash
```

**Go to Home Directory**:
```bash
cd /home/hyde
```

**Tips before starting commands**:
- There is a test file named `test-file.txt` with `Hello, World!` written into it in the home directory. You can use this file to test the CLI commands.
- Copy and paste the hyde commands below to test the CLI functionality. Bash Completion is enabled, so you can also use the `tab` key to auto-complete commands and options to make it easier to use.

**Encrypt a file**:
```bash
hyde encrypt test-file.txt --password abc123 --delete
```
`--password` or `-p`: The password used for encryption. This is required.

`--delete`: This optional flag tells the tool to delete the original file after it has been successfully encrypted.  **Use with caution.**

Output file will be named `test-file.txt.enc` and will be stored in the same directory as the original file.

**Decrypt a file**:
```bash
hyde decrypt test-file.txt.enc --password abc123 --delete
```
`--password` or `-p`: The password used for decryption. This is required.

`--delete`: This optional flag tells the tool to delete the encrypted file after it has been successfully decrypted. **Use with caution.**

Output file will be named `test-file.txt` with the original content restored. It will be stored in the same directory as the encrypted file.


## 🛠️ Development & Future

Hyde aims to remain a simple, focused tool. The priority is maintaining security and reliability.

**Maintenance:** Security patches and dependency updates will be applied as needed.

**Simplicity:** Major feature additions are unlikely unless they strongly align with the core principles without adding significant complexity. Consistency and ease-of-use are paramount.