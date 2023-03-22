# End to End testing guide

## Table of Contents

- [Installing a specific VMClarity build on AWS](#installing-a-specific-vmclarity-build-on-aws)
  - [1. Build the containers and publish them to your docker registry](#1-build-the-containers-and-publish-them-to-your-docker-registry)
  - [2. Update installation/aws/VMClarity.cfn](#2-update-installationawsvmclaritycfn)
  - [3. Install VMClarity cloudformation](#3-install-vmclarity-cloudformation)
  - [4. Ensure that VMClarity backend is working correctly](#4-ensure-that-vmclarity-backend-is-working-correctly)
- [Performing an end to end test](#performing-an-end-to-end-test)

## Installing a specific VMClarity build on AWS

### 1. Build the containers and publish them to your docker registry

```
DOCKER_REGISTRY=<your docker registry> make push-docker
```

### 2. Update installation/aws/VMClarity.cfn

Update the cloud formation with the pushed docker images, for example:

```
@@ -123,7 +123,7 @@ Resources:
                     DATABASE_DRIVER=LOCAL
                     BACKEND_REST_HOST=__BACKEND_REST_HOST__
                     BACKEND_REST_PORT=8888
-                    SCANNER_CONTAINER_IMAGE=ghcr.io/openclarity/vmclarity-cli:latest
+                    SCANNER_CONTAINER_IMAGE=tehsmash/vmclarity-cli:9bba94334c1de1aeed63ed12de3784d561fc4f1b
                   - JobImageID: !FindInMap
                       - AWSRegionArch2AMI
                       - !Ref "AWS::Region"
@@ -145,13 +145,13 @@ Resources:
                 ExecStartPre=-/usr/bin/docker stop %n
                 ExecStartPre=-/usr/bin/docker rm %n
                 ExecStartPre=/usr/bin/mkdir -p /opt/vmclarity
-                ExecStartPre=/usr/bin/docker pull ghcr.io/openclarity/vmclarity-backend:latest
+                ExecStartPre=/usr/bin/docker pull tehsmash/vmclarity-backend:9bba94334c1de1aeed63ed12de3784d561fc4f1b
                 ExecStart=/usr/bin/docker run \
                   --rm --name %n \
                   -p 0.0.0.0:8888:8888/tcp \
                   -v /opt/vmclarity:/data \
                   --env-file /etc/vmclarity/config.env \
-                  ghcr.io/openclarity/vmclarity-backend:latest run --log-level info
+                  tehsmash/vmclarity-backend:9bba94334c1de1aeed63ed12de3784d561fc4f1b run --log-level info

                 [Install]
                 WantedBy=multi-user.target
```

### 3. Install VMClarity cloudformation

1. Ensure you have an SSH key pair uploaded to AWS Ec2
2. Go to CloudFormation -> Create Stack -> Upload template.
3. Upload the modified VMClarity.cfn
4. Follow the wizard through to the end
5. Wait for install to complete

### 4. Ensure that VMClarity backend is working correctly

1. Get the IP address from the CloudFormation stack's Output Tab
2. `ssh ubuntu@<ip address>`
3. Check the VMClarity Logs

   ```
   sudo journalctl -u vmclarity
   ```

## Performing an end to end test

1. Copy the example [scanConfig.json](/docs/scanConfig.json) into the ubuntu user's home directory

   ```
   scp scanConfig.json ubuntu@<ip address>:~/scanConfig.json
   ```

2. Edit the scanConfig.json

   a. Give the scan config a unique name

   b. Enable the different scan families you want:

    ```
    "scanFamiliesConfig": {
      "sbom": {
        "enabled": true
      },
      "vulnerabilities": {
        "enabled": true
      },
      "exploits": {
        "enabled": true
      }
    },
    ```

   c. Configure the scope of the test

      * By Region, VPC or Security group:

        ```
        "scope" {
          "objectType": "AwsScanScope",
          "regions": [
            {
             "name": "eu-west-1",
             "vpcs": [
               {
                 "name": "<name of vpc>",
                 "securityGroups": [
                   {
                     "name": "<name of sec group>"
                   }
                 ]
               }
             ]
            }
          ]
        }
        ```

      * By tag:

        ```
        "scope": {
          "instanceTagSelector": [
            {
              "key": "<key>",
              "value": "<value>"
            }
          ]
        }
        ```

      * All:

        ```
        "scope": {
          "all": true
        }
        ```

   d. Set operationTime to the time you want the scan to run. As long as the time
      is in the future it can be within seconds.

3. While ssh'd into the VMClarity server run

   ```
   curl -X POST http://localhost:8888/api/scanConfigs -H 'Content-Type: application/json' -d @scanConfig.json
   ```

4. Check VMClarity logs to ensure that everything is performing as expected

   ```
   sudo journalctl -u vmclarity
   ```

5. Monitor the scan results

   * Get scans:

     ```
     curl -X GET http://localhost:8888/api/scans
     ```

     After the operationTime in the scan config created above there should be a new
     scan object created in Pending.

     Once discovery has been performed, the scan's "targets" list should be
     populated will all the targets to be scanned by this scan.

     The scan will then create all the "scanResults" for tracking the scan
     process for each target. When that is completed the scan will move to
     "InProgress".

   * Get Scan Results:

     ```
     curl -X GET http://localhost:8888/api/scanResults
     ```
