

# Build

`go build`

Cross-compile to raspberry pi 0 :

`GOOS=linux GOARCH=arm GOARM=5 go build`

Then scp eg.

`scp miam pi@192.168.1.21:/tmp/miam`

# Launch

`./miam`

