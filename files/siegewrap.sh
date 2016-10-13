#!/bin/bash
sleep $((30-$(date +%s) % 30))
siege $@
