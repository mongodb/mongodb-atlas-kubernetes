{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  buildInputs = [
    pkgs.golangci-lint
    pkgs.yq-go
    pkgs.kubebuilder
    pkgs.jq
    pkgs.go
    pkgs.gotests
    pkgs.act
    pkgs.kubectl
    pkgs.docker
    pkgs.kubernetes-controller-tools
    pkgs.kustomize_4
    pkgs.git
    pkgs.envsubst
    pkgs.wget
    pkgs.cosign
    pkgs.kubernetes-helm
    pkgs.govulncheck
    pkgs.gotools
    pkgs.go-licenses
    pkgs.ginkgo
    pkgs.operator-sdk
    pkgs.shellcheck
  ];
}
