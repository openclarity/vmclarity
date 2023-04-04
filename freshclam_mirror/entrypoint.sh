#!/bin/sh

# Run freshclam and cron in the background
freshclam --datadir=/usr/share/nginx/html/clamav --quiet &
cron -f &

# Start Nginx
exec nginx -g 'daemon off;'
