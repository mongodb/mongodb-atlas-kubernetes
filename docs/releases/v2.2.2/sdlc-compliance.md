SSDLC Compliance Report: Atlas Kubernetes Operator Manager v2.2.2
=================================================================

- Release Creators: s.urbaniak@mongodb.com
- Created On:       2024-04-25

Overview:

- **Product and Release Name**

    - Atlas Kubernetes Operator v2.2.2, 2024-04-25.
    - Release Type: Minor

- **Process Document**
  - http://go/how-we-develop-software-doc

- **Tool used to track third party vulnerabilities**
  - Silk

- **Dependency Information**
  - See SBOMS Lite manifests (CycloneDX in JSON format) for [Intel](./linux-amd64.sbom.json) or [ARM](./linux-arm64.sbom.json)

- **Static Analysis Report**
  - No reports (filtered before release by CI tests)
  - List of explicitly ignored vulnerabilities:
    - https://pkg.go.dev/vuln/GO-2024-2687 not bumped to go 1.22.2 yet

- **Release Signature Report**
  - Image signatures enforced by CI pipeline.
  - See [Signature verification instructions here](../../dev/signed-images.md)
  - Self-verification shortcut:
    ```shell
    make verify IMG=mongodb/mongodb-atlas-kubernetes-operator:2.2.2 SIGNATURE_REPO=mongodb/signatures
    ```

- **Security Testing Report**
  - Available as needed from Cloud Security.

- **Security Assessment Report**
  - Available as needed from Cloud Security.

Assumptions and attestations:

1. Internal processes are used to ensure CVEs are identified and mitigated within SLAs.

2. The Dependency document does not specify third party OSS CVEs fixed by the release and the date we discovered them.

3. There is no CycloneDX field for original/modified CVSS scor or discovery date. The `x-` prefix indicates this.

3. Assumption: We can include the SBOMs as links to read-only files on S3. The links can be included as metadata or text file links in release artifacts e.g. as labels on OCI containers.