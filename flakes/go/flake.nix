{
  description = "A flake for Go 1.25 versions";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.go_1_25.overrideAttrs (old: {
          version = "1.25.1";
          src = pkgs.fetchurl {
            url = "https://golang.org/dl/go1.25.1.linux-amd64.tar.gz";
            sha256 = "sha256-dxag2UCg9q6OHzs/TzYpncU+MbFoQNvRcSVDEsQcoS4=";
          };
        });
      });
}
