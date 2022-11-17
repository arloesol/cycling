#!/usr/bin/env bash
for src in ../../../../../content/routes/*/*.md; do go run gpxdata.go $src; done