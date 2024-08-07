{ pkgs ? import <nixpkgs> { } }:
pkgs.mkShell {
  buildInputs = [
    pkgs.golangci-lint
    pkgs.yq-go
    pkgs.kubebuilder
    pkgs.jq
    pkgs.go
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

  shellHook = ''
    export PATH=$PATH:${pkgs.go}/bin
    export PATH=$PATH:${pkgs.cosign}/bin
    export PATH=$PATH:${pkgs.kubebuilder}/bin
    export PATH=$PATH:${pkgs.go-licenses}/bin
    export PATH=$PATH:${pkgs.docker}/bin
    export PATH=$PATH:${pkgs.yq-go}/bin
    export PATH=$PATH:${pkgs.wget}/bin
    export PATH=$PATH:${pkgs.envsubst}/bin
    export PATH=$PATH:${pkgs.git}/bin
    export PATH=$PATH:${pkgs.jq}/bin
    export PATH=$PATH:${pkgs.shellcheck}/bin
    export PATH=$PATH:${pkgs.kubernetes-controller-tools}/bin
    export PATH=$PATH:${pkgs.kubectl}/bin
    export PATH=$PATH:${pkgs.kubernetes-helm}/bin
    export PATH=$PATH:${pkgs.act}/bin
    export PATH=$PATH:${pkgs.kustomize}/bin
    export PATH=$PATH:${pkgs.govulncheck}/bin
    export PATH=$PATH:${pkgs.gotools}/bin
    export PATH=$PATH:${pkgs.ginkgo}/bin
    export PATH=$PATH:${pkgs.operator-sdk}/bin
  '';
}
