# Go Variables
CRYPTO_CMD = go run cmd/cryptid.go

# Commands
AES_CMD = $(CRYPTO_CMD) aes
PASSWORD_CMD = $(CRYPTO_CMD) password
JWT_CMD = $(CRYPTO_CMD) jwt

# Encryption variables
FILE_TO_ENCRYPT = .env
FILE_TO_DECRYPT = env.enc
PASSPHRASE = 123

VERSION = "v1.1.1"

# Targets
all:
	$(CRYPTO_CMD)

build:
	sh build.sh $(VERSION)

format:
	find . -type f -name "*.go" -exec go fmt {} \;

# Password targets
complex:
	$(PASSWORD_CMD) complex

length:
	$(PASSWORD_CMD) complex -length 20

no-numbers:
	$(PASSWORD_CMD) complex -no-numbers

no-symbols:
	$(PASSWORD_CMD) complex -no-symbols
	
only-alphabets:
	$(PASSWORD_CMD) complex -no-numbers -no-symbols

no-console:
	$(PASSWORD_CMD) complex -no-console

quiet:
	$(PASSWORD_CMD) complex -quiet

passphrase:
	$(PASSWORD_CMD) passphrase
# JWT targets

256:
	$(JWT_CMD) hs256

384:
	$(JWT_CMD) hs384

512:
	$(JWT_CMD) hs512

# AES targets

encrypt:
	$(AES_CMD) encrypt -f $(FILE_TO_ENCRYPT) -o $(FILE_TO_DECRYPT) -p "$(PASSPHRASE)"
	@rm $(FILE_TO_ENCRYPT)

decrypt:
	$(AES_CMD) decrypt -f $(FILE_TO_DECRYPT) -o $(FILE_TO_ENCRYPT) -p "$(PASSPHRASE)"
	@rm $(FILE_TO_DECRYPT)