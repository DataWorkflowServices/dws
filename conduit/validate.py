#! /usr/bin/python3

import sys
import yaml
import io

from openapi_schema_validator import validate

args = sys.argv[1:]
crdFile = args[0]
example = sys.__stdin__.read()

crd = yaml.safe_load(io.open(crdFile))
example = yaml.safe_load(example)

schema = crd["spec"]["versions"][0]["schema"]["openAPIV3Schema"]

validate(example,schema)
