#!/bin/bash
for language in java python; do
    docker build -t ipv6tester-$language -f "Dockerfile.$language" .
done
