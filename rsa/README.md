Create a REST API that has an endpoint to generate an RSA key pair. The endpoint must save the name of the key, the public key and the private key in a database. Encrypt the field where the private key is saved with AES-256

Create an endpoint that allows you to list the keys and search by name
Create an endpoint that receives a plain text and the ID of the key, and returns the encrypted text with the key
Create an endpoint that receives the encrypted text and key ID, and returns the plain text

## Generate the swagger Documents


install

go get -u github.com/go-swagger/go-swagger/cmd/swagger

## execute the command

swagger generate spec -o ./swagger.json --scan-models

the file swagger.json upload in
 
https://editor.swagger.io/