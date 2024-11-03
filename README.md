# Cryptid

Cryptid is a versatile CLI tool designed for secure and efficient cryptographic operations. It supports:

- Generating complex passwords
- Creating JWT secrets
- AES file encryption and decryption

This tool is ideal for developers and security professionals who need reliable cryptographic utilities in their workflow.

## Usage

### Password Generation

```console
cryptid password complex -length 16
# Output: ljLoqT5BcfB27@BH
```

```console
cryptid password complex -l 20 -no-numbers -no-symbols
# Output: YCuPnnmrLWcOQzJjYgcB
```

### JWT Secret

```console
cryptid jwt hs256
# Output: m^2KW?2P%NF6ci3Ns8AV5E)^2*!2?(?S
```

### AES Encryption

```console
cryptid aes encrypt -f secret.txt -o secret.enc -passphrase mypass
```

```console
cryptid aes decrypt -f secret.enc -o secret.txt -passphrase mypass
```

## Downloads

### macOS (Apple Silicon)

```console
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/darwin/arm64/cryptid
```

### macOS (Intel)

```console
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/darwin/amd64/cryptid
```

### Linux (x86_64)

```console
 curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/linux/amd64/cryptid
```

### Linux (ARM64)

```console
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/linux/arm64/cryptid
```

### Windows (x86_64)

```console
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/windows/amd64/cryptid.exe
```

### Windows (ARM64)

```console
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/windows/arm64/cryptid.exe
```
