#!/usr/bin/env bash

dir="release"
bin="license-gen"

go build -o $bin

[[ -d $dir ]] && rm -rf $dir
mkdir $dir

mv $bin $dir
cp -rf .env $dir
cp -rf config.yml $dir
cp -rf assets $dir
