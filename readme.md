sqrlGateKeeper
===============

A reverse proxy that provides SQRL authentication.

This is only proof of concept. DO NOT USE FOR PRODUCTION SERVERS

If you have an internal server where you want to restrict access you can
put sqrlGateKeeper in front.

## Requirements
GO 1.5

## Usage

./sqrlGateKeeper -s 127.0.0.1:8080 -l 80 -p 12345

-s Sets the remote server  
-l Sets the listening port  
-p Sets the admin password. This is used for adding new users.
