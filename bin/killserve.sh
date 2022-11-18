#!/usr/bin/env bash
kill $(ps aux | grep '[h]ugo server --disableFastRender -e production' | awk '{print $2}')