#!/usr/bin/env bash
protoc -I proto/ proto/customerservice.proto --go_out=plugins=grpc:proto
