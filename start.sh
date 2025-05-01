#!/bin/sh
/app/server -dev &
exec nginx -g 'daemon off;'