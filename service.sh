# Copyright (c) 2017-2026 Onur Yaşar
# Licensed under AGPL v3 + Commercial Exception
# See LICENSE.txt

# https://github.com/rymory/rymory-core
# rymory.org 
# onuryasar.org
# onxorg@proton.me 

echo ; 

export JWT_ISSUER="security.lemoras.com"
export ROOT_ACCOUNT="root@lemoras.com"
export TOKEN_SECRET_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export TOKEN_ROOT_SECRET_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export DATABASE_URL="XXXXXXXXXXXXXXXXXXXXXXXX"
export SPACES_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export SPACES_SECRET="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export SPACES_NAME="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export SPACES_ENDPOINT="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export VALID_API_URL="https://worker.lemoras.com/security/validation"
export TICKET_SECRET_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export RATE_SECRET_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export X_API_KEY="XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
export COOKIE_ALLOWED_DOMAINS="dev.local,lemoras.com"



echo set environment variables

echo deploy...

doctl serverless connect lemoras-core

doctl serverless deploy . --remote-build 

sh clear.sh

echo done
