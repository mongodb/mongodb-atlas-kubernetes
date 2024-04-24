SSDLC Compliance Report: Atlas Kubernetes Operator Manager v${VERSION}
=================================================================

- Release Creators: ${AUTHORS}
- Created On:       ${DATE}

Overview:

- **Product and Release Name**

    - Atlas Kubernetes Operator v${VERSION}, ${DATE}.
    - Release Type: ${RELEASE_TYPE}

- **Process Document**
  - http://go/how-we-develop-software-doc

- **Tool used to track third party vulnerabilities**
  - Silk

- **Dependency Information**
  - See SBOMS Lite manifests (CycloneDX in JSON format) for [Intel](./linux-amd64.sbom.json) or [ARM](./linux-arm64.sbom.json)

- **Static Analysis Report**
  - No reports (filtered before release by CI tests)${IGNORED_VULNERABILITIES}

- **Release Signature Report**
  - Image signatures enforced by CI pipeline.
  - See [Signature verification instructions here](../../dev/signed-images.md)
  - Self-verification shortcut:
    ```shell
    make verify IMG=mongodb/mongodb-atlas-kubernetes-operator:${VERSION} SIGNATURE_REPO=mongodb/signatures
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