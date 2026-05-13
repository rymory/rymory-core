# Copyright (c) 2017-2026 Onur Yaşar
# Licensed under AGPL v3 + Commercial Exception
# See LICENSE.txt

# https://github.com/rymory/rymory-core
# rymory.org 
# onuryasar.org
# onxorg@proton.me 

#!/bin/bash

set -e

virtualenv virtualenv
source virtualenv/bin/activate
pip install -r requirements.txt
deactivate