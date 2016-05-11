#!/usr/bin/env bash
token=$(cf oauth-token)
http --verify=no --verbose http://localhost:8080/userinfo Authorization:"$token"