// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Info\"
      },
      \"propertyName\": \"eventType\"
    }
  }

Value:
  {
    \"eventType\": \"Findings\",
    \"findings\": [
      {},
      {
        \"annotations\": [
          {
            \"key\": \"scanner/scanned-by\",
            \"value\": \"cisdocker\"
          }
        ],
        \"findingInfo\": {
          \"category\": \"best-practice\",
          \"description\": \"Avoid 'latest' tag\
\",
          \"id\": \"DKL-DI-0006\",
          \"location\": \"nginx\",
          \"message\": \"Avoid latest tag\",
          \"objectType\": \"Misconfiguration\",
          \"severity\": \"MisconfigurationMediumSeverity\"
        },
        \"input\": {
          \"path\": \"nginx\",
          \"type\": \"IMAGE\"
        },
        \"scanID\": \"ff34ee37-ba88-4cf7-a25e-ffddb928dd29\"
      },
      {
        \"annotations\": [
          {
            \"key\": \"scanner/scanned-by\",
            \"value\": \"cisdocker\"
          }
        ],
        \"findingInfo\": {
          \"category\": \"best-practice\",
          \"description\": \"setuid file: urwxr-xr-x /usr/bin/newgrp\
setuid file: urwxr-xr-x /usr/bin/passwd\
setuid file: urwxr-xr-x /usr/bin/chsh\
setuid file: urwxr-xr-x /usr/bin/chfn\
setuid file: urwxr-xr-x /usr/bin/su\
setgid file: grwxr-xr-x /usr/bin/chage\
setuid file: urwxr-xr-x /usr/bin/umount\
setuid file: urwxr-xr-x /usr/bin/mount\
setgid file: grwxr-xr-x /usr/sbin/unix_chkpwd\
setgid file: grwxr-xr-x /usr/bin/expiry\
setgid file: grwxr-xr-x /usr/bin/wall\
setuid file: urwxr-xr-x /usr/bin/gpasswd\
\",
          \"id\": \"CIS-DI-0008\",
          \"location\": \"nginx\",
          \"message\": \"Confirm safety of setuid/setgid files\",
          \"objectType\": \"Misconfiguration\",
          \"severity\": \"MisconfigurationLowSeverity\"
        },
        \"input\": {
          \"path\": \"nginx\",
          \"type\": \"IMAGE\"
        },
        \"scanID\": \"ff34ee37-ba88-4cf7-a25e-ffddb928dd29\"
      },
      {
        \"annotations\": [
          {
            \"key\": \"scanner/scanned-by\",
            \"value\": \"cisdocker\"
          }
        ],
        \"findingInfo\": {
          \"category\": \"best-practice\",
          \"description\": \"Last user should not be root\
\",
          \"id\": \"CIS-DI-0001\",
          \"location\": \"nginx\",
          \"message\": \"Create a user for the container\",
          \"objectType\": \"Misconfiguration\",
          \"severity\": \"MisconfigurationMediumSeverity\"
        },
        \"input\": {
          \"path\": \"nginx\",
          \"type\": \"IMAGE\"
        },
        \"scanID\": \"ff34ee37-ba88-4cf7-a25e-ffddb928dd29\"
      },
      {
        \"annotations\": [
          {
            \"key\": \"scanner/scanned-by\",
            \"value\": \"cisdocker\"
          }
        ],
        \"findingInfo\": {
          \"category\": \"best-practice\",
          \"description\": \"not found HEALTHCHECK statement\
\",
          \"id\": \"CIS-DI-0006\",
          \"location\": \"nginx\",
          \"message\": \"Add HEALTHCHECK instruction to the container image\",
          \"objectType\": \"Misconfiguration\",
          \"severity\": \"MisconfigurationLowSeverity\"
        },
        \"input\": {
          \"path\": \"nginx\",
          \"type\": \"IMAGE\"
        },
        \"scanID\": \"ff34ee37-ba88-4cf7-a25e-ffddb928dd29\"
      },
      {
        \"annotations\": [
          {
            \"key\": \"scanner/scanned-by\",
            \"value\": \"cisdocker\"
          }
        ],
        \"findingInfo\": {
          \"category\": \"best-practice\",
          \"description\": \"export DOCKER_CONTENT_TRUST=1 before docker pull/build\
\",
          \"id\": \"CIS-DI-0005\",
          \"location\": \"nginx\",
          \"message\": \"Enable Content trust for Docker\",
          \"objectType\": \"Misconfiguration\",
          \"severity\": \"MisconfigurationLowSeverity\"
        },
        \"input\": {
          \"path\": \"nginx\",
          \"type\": \"IMAGE\"
        },
        \"scanID\": \"ff34ee37-ba88-4cf7-a25e-ffddb928dd29\"
      }
    ]
  }
","latency":603500,"latency_human":"603.5µs","bytes_in":2546,"bytes_out":159}
{"time":"2024-03-19T09:54:43.091327+01:00","id":"","remote_ip":"127.0.0.1","host":"0.0.0.0:8765","method":"POST","uri":"/scan/ff34ee37-ba88-4cf7-a25e-ffddb928dd29/event","user_agent":"Go-http-client/1.1","status":400,"error":"code=400,
message=request body has an error: doesn't match schema #/components/schemas/ScanEvent:
Error at \"/eventInfo\": doesn't match any schema from \"anyOf\", internal=request body has an error:
doesn't match schema #/components/schemas/ScanEvent: Error at \"/eventInfo\": doesn't match any schema from \"anyOf\"
Schema:
  {
    \"anyOf\": [
      {
        \"$ref\": \"#/components/schemas/ScannerHandshakeEventInfo\"
      },
      {
        \"$ref\": \"#/components/schemas/ScannerHeartbeatEventInfo\"
      },
      {
        \"$ref\": \"#/components/schemas/ScannerFindingsEventInfo\"
      }
    ],
    \"discriminator\": {
      \"mapping\": {
        \"Findings\": \"#/components/schemas/ScannerFindingsEventInfo\",
        \"Handshake\": \"#/components/schemas/ScannerHandshakeEventInfo\",
        \"Heartbeat\": \"#/components/schemas/ScannerHeartbeatEventInfo\"
      },
      \"propertyName\": \"eventType\"
    }
  }

Value:
  {
    \"eventType\": \"Heartbeat\",
    \"message\": \"scan completed successfully\",
    \"state\": \"Completed\",
    \"summary\": {
      \"jobsDone\": 6,
      \"jobsFailed\": 0,
      \"jobsRemaining\": 0
    }
  }
","latency":207667,"latency_human":"207.667µs","bytes_in":157,"bytes_out":159}
*/

package golang

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=types/types.cfg.yaml ../openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=client/client.cfg.yaml ../openapi.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=server/server.cfg.yaml ../openapi.yaml
