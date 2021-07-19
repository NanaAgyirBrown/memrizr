.PHONY: keypair migrate-create migrate-up migrate-down migrate-force

# PWD = $(Shell pwd)
PWD = $(CURDIR)
ACCPATH = $(PWD)/account
MPATH = $(ACCPATH)/migrations
PORT = 5432

# Default number of migrations to execute up or down
N = 1

# Create keypair should be in your file below

# creating migration files in sequential order -seq
migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

# applying db updates
migrate-up:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable up $(N)

# Reverting db updates
migrate-down:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

# force a version number (sequence number) when there is a need to fix an error in a previous migration. Fall-back will be to manually fix the errors
migrate-force:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)

keypair:
	@echo "creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCPATH)/rsa_public_$(ENV).pem