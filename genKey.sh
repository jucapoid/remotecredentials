#!/bin/bash

# To generate keys from the shell

date +%s | sha256sum | base64 | head -c 32 ; echo

