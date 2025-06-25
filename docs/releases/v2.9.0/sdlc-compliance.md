SSDLC Compliance Report: Atlas Kubernetes Operator Manager v2.9.0
=================================================================

- Release Creators: jose.vazquez@mongodb.com
- Created On:       2025-06-25

Overview:

- **Product and Release Name**

    - Atlas Kubernetes Operator v2.9.0, 2025-06-25.

- **Process Document**
  - http://go/how-we-develop-software-doc

- **Tool used to track third party vulnerabilities**
  - [Kondukto](https://arcticglow.kondukto.io/)

- **Dependency Information**
  - See SBOMS Lite manifests (CycloneDX in JSON format) for `Intel` and `ARM` are to be found [here](.)
  - See [instructions on how the SBOMs are generated or how to generate them manually](../../dev/image-sboms.md)

- **Static Analysis Report**
  - No SAST findings. Our CI system blocks merges on any SAST findings.
  - No vulnerabilities were ignored for this release.

- **Release Signature Report**
  - Image signatures enforced by CI pipeline.
  - See [Signature verification instructions here](../../dev/signed-images.md)
  - Self-verification shortcut:
    ```shell
    make verify IMG=mongodb/mongodb-atlas-kubernetes-operator:2.9.0 SIGNATURE_REPO=mongodb/signatures
    ```

- **Security Testing Report**
  - Available as needed from Cloud Security.

- **Security Assessment Report**
  - Available as needed from Cloud Security.

Assumptions and attestations:

- Internal processes are used to ensure CVEs are identified and mitigated within SLAs.

- All Operator images are signed by MongoDB, with signatures stored at `docker.io/mongodb/signatures`.
