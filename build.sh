#!/bin/bash

git pull
go build main.go
systemctl restart cs411
systemctl status cs411
