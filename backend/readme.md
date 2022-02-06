# Miam Backend

# Build

`go build`

Cross-compile to raspberry pi 0 :

`CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -v`

Then scp eg.

`scp miam pi@192.168.1.21:/tmp/miam`

# Launch

`./miam`

# TODO

- health
- config
- gestion d'erreur
    - mutualisation dans les handlers + utilisation de wrap pour avoir un semblant de stacktrace + messages pas statiques pour les erreurs dans common
- version via build and -ldflag
- cross-compiling dans image
- docker with bare exec
    - seems difficult since we use cgo. distroless is probably the best here
- lint

# See what's going on in the database

sqlitebrowser and boltBrowser can be used
