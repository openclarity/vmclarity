# VMClarity Freshclam Mirror

This Dockerfile provides definition for a server that serves ClamAV signature files and updates them every 3 hours.
This is achieved using nginx, freshclam and cron.

# Usage

## Directing freshclam to use this mirror
To direct your freshclam instance to this mirror, configure your freshclam.conf file as such:
```
PrivateMirror http://<ip>:1000
ScriptedUpdates no
```

## Building
```
docker build -t <registry-name>/vmclarity-freshclam-mirror .
```

## Running
```
docker run -d -p 1000:80 --name vmclarity-freshclam-mirror <registry-name>/vmclarity-freshclam-mirror
```

## Manual download
```
curl -X GET http://<ip>:1000/clamav/main.cvd --output main.cvd
```
