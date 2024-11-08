# Cryptid

Cryptid is a versatile CLI tool designed for secure and efficient cryptographic operations. It supports:

- Generating complex passwords
- Creating JWT secrets
- AES file encryption and decryption

This tool is ideal for developers and security professionals who need reliable cryptographic utilities in their workflow.

## Usage

### Password Generation

```bash
cryptid password complex -length 16
# Output: O(lIJj+Zp|=<0{5-
```

```bash
cryptid password complex -l 20 -no-numbers -no-symbols
# Output: YCuPnnmrLWcOQzJjYgcB
```

### JWT Secret

```bash
cryptid jwt hs256
# Output: m^2KW?2P%NF6ci3Ns8AV5E)^2*!2?(?S
```

```bash
cryptid jwt hs384
# Output: Cw)WkT?3>R>vN![;ZxZagfPr[X,vCSpN_42N.1XJlT0OnmVu
```

```bash
cryptid jwt hs512
# Output: x$;M6QH806^T<_#PH7,t#FEyLcO:@zfu+D2)@3C*W5MOkw3P?s0<9}AZ84EgR,uh
```

### AES Encryption

```bash
cryptid aes encrypt -f secret.txt -o secret.enc -passphrase mypass
```

```bash
cryptid aes decrypt -f secret.enc -o secret.txt -passphrase mypass
```

## Installation

Will be published to public repositories soon. For now, you can download the binary for your platform from the links below.

## Downloads

### macOS (Apple Silicon)

```bash
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/darwin/arm64/cryptid
```

### macOS (Intel)

```bash
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/darwin/amd64/cryptid
```

### Linux (x86_64)

```bash
 curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/linux/amd64/cryptid
```

### Linux (ARM64)

```bash
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/linux/arm64/cryptid
```

### Windows (x86_64)

```bash
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/windows/amd64/cryptid.exe
```

### Windows (ARM64)

```bash
curl -LJO https://raw.githubusercontent.com/edsonjaramillo/cryptid/main/dist/windows/arm64/cryptid.exe
```
