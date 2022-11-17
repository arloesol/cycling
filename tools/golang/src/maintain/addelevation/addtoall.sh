#!/usr/bin/env bash
for src in ../../../../../static/gpxfiles/*/*.gpx; do go run addelevation.go $src; done