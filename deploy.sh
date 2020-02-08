#!/bin/bash

set -e

docker build -t miam-deployer -f deploy/Dockerfile.rasp .
docker run -it --rm miam-deployer --ask-pass /miam/deploy/playbook.yml

# pour avoir les logs du services sur le rasp:

# systemctl status miam.service
