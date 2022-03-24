Readme

```
// Just fork me a certificate
// Just generate a vanilla, self signed, non CA, x509 certificate, without extensions and all the openssl-isms
// ...you probably need the extensions for your use case though
// Do not use self signed CA=TRUE (ala OpenSSL) certificates for mTLS in the browser... will cause you trust any external certificate signed with your "malign" test cert, bad idea.
// You can do this with openSSL, but it's a massive hassle.
```

Running:

```
// Generate a key and certificate
// Modify consts: length, days, signature, if required
go run . -s "SubjectName"

// View certificate contents
openssl x509 -in <SubjectName>.crt -text

// Bundle into a p12
openssl pkcs12 -export -inkey <SubjectName>.pem -in <SubjectName>.crt -name <SubjectName> -out <SubjectName>.p12
```