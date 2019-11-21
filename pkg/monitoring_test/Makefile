.SILENT :

DAYS:=1460

SUBECT:=example.com

CERTIFICATE:=ca.pem

ALL: cert.pem

## Delete all generated certificates
clean:
	@echo "Deleting all generated files..."
	-rm *.pem *.csr *.srl
.PHONY: clean

## Generate CA key and certificate
ca.pem:
	@echo "Generating CA key and certificate..."
	echo 01 > ca.srl
	openssl genrsa -out ca-key.pem 2048
	openssl req -subj '/CN=APIMon/' -new -x509 -days $(DAYS) -key ca-key.pem -out ca.pem

## Generate client key
key.pem:
	@echo "Generating key..."
	openssl genrsa -out key.pem 2048

## Generate client certificate
cert.pem: key.pem ca.pem
	@echo "Generating client certificate..."
	openssl req -subj '/CN=test' -new -key key.pem -out test.csr
	openssl x509 -req -days $(DAYS) -in test.csr -CA ca.pem -CAkey ca-key.pem -out cert.pem

## Show certificate content
debug:
	openssl x509 -in $(CERTIFICATE) -text
.PHONY: clean