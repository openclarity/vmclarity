<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./img/logos/VMClarity-logo-dark-bg-horizontal@4x.png">
  <source media="(prefers-color-scheme: light)" srcset="./img/logos/VMClarity-logo-light-bg-horizontal@4x.png">
  <img alt="VMClarity Logo" src="./img/logos/VMClarity-logo-light-bg-horizontal@4x.png">
</picture>

VMClarity is an open source tool for agentless detection and management of Virtual Machine
Software Bill Of Materials (SBOM) and security threats such as vulnerabilities, exploits, malware, rootkits, misconfigurations and leaked secrets.

<img src="./img/vmclarity_demo.gif" alt="VMClarity demo" />


# Table of Contents<!-- omit in toc -->

- [Why VMClarity?](#why-vmclarity)
- [Quick Start](#quick-start)
- [Overview](#overview)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Code of Conduct](#code-of-conduct)
- [License](#license)

# Why VMClarity?

Virtual machines (VMs) are the most used service across all hypescalers. AWS,
Azure, GCP, and others have virtual computing services that are used not only
as standalone VM services but also as the most popular method for hosting
containers (e.g., Docker, Kubernetes).

VMs are vulnerable to multiple threats:
- Software vulnerabilties
- Leaked Secrets/Passwords
- Malware
- System Misconfiguration
- Rootkits

There are many very good open source and commercial-based solutions for
providing threat detection for VMs, manifesting the different threat categories above.

However, there are challenges with assembling and managing these tools yourself:
- Complex installation, configuration, and reporting
- Integration with deployment automation
- Siloed reporting and visualization

The VMClarity project is focused on unifying detection and management of VM security threats in an agentless manner.

# Quick start
## Install VMClarity

<details><summary>On AWS</summary>
<p>

1. Start the CloudFormation [wizard](https://console.aws.amazon.com/cloudformation/home#/stacks/create/review?stackName=VMClarity&templateURL=https://s3.eu-west-2.amazonaws.com/vmclarity-v0.4.0/VmClarity.cfn), or upload the [latest](https://github.com/openclarity/vmclarity/releases/latest) CloudFormation template 
2. Specify the SSH key to be used to connect to VMClarity under 'KeyName'
3. Once deployed, copy VmClarity SSH Address from the "Outputs" tab

</p>
</details>

For a detailed installation guide, please see [AWS](installation/aws/README.md).
For a detailed UI tour, please see [tour](TOUR.md).

## Access VMClarity UI
1. Open an SSH tunnel to VMClarity server
    ```
    ssh -N -L 8888:localhost:8888 -i  "<Path to the SSH key specified in 1.1.i>" ubuntu@<VmClarity SSH Address specified in 1.1.ii>
    ```

2. Access VMClarity UI in the browser: http://localhost:8888/
3. Access the [API](/api/openapi.yaml) via http://localhost:8888/api



# Overview

VMClarity uses a pluggable scanning infrastructure to provide:
- SBOM analysis
- Package and OS vulnerability detection
- Exploit detection
- Leaked secret detection
- Malware detection
- Misconfiguration detection
- Rootkit detection

The pluggable scanning infrastructure uses several tools that can be
enabled/disabled on an individual basis. VMClarity normalizes, merges and
provides a robust visualization of the results from these various tools.

These tools include:
- SBOM Generation and Analysis
  - [Syft](https://github.com/anchore/syft)
  - [Trivy](https://github.com/aquasecurity/trivy)
  - [Cyclonedx-gomod](https://github.com/CycloneDX/cyclonedx-gomod)
- Vulnerability detection
  - [Grype](https://github.com/anchore/grype)
  - [Trivy](https://github.com/aquasecurity/trivy)
  - [Dependency-Track](https://github.com/DependencyTrack/dependency-track)
- Exploits
  - [Go exploit db](https://github.com/vulsio/go-exploitdb)
- Secrets
  - [gitleaks](https://github.com/gitleaks/gitleaks)
- Malware
  - [ClamAV](https://github.com/Cisco-Talos/clamav)
- Misconfiguration
  - [Lynis](https://github.com/CISOfy/lynis)
- Rootkits
  - [Chkrootkit](https://github.com/Magentron/chkrootkit)

A high-level architecture overview is available [here](ARCHITECTURE.md)

# Roadmap
VMClarity project roadmap is available [here](https://github.com/orgs/openclarity/projects/2/views/7).

# Contributing

If you are ready to jump in and test, add code, or help with documentation,
please follow the instructions on our [contributing guide](/CONTRIBUTING.md)
for details on how to open issues, setup VMClarity for development and test.

# Code of Conduct

You can view our code of conduct [here](/CODE_OF_CONDUCT.md).

# License

[Apache License, Version 2.0](/LICENSE)
