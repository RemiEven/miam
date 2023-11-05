# Miam Backend

# Build

`go build`

Cross-compile to raspberry pi 0 :

`CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 go build -v`

Then scp eg.

`scp miam pi@192.168.1.21:/tmp/miam`

# Launch

`./miam`

# See what's going on in the database

sqlitebrowser and boltBrowser can be used
