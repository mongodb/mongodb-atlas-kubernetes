{
  description = "A flake for multi-platform Go toolchain selection";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        goArchives = {
          "x86_64-linux" = {
            goFileName = "go1.25.0.linux-amd64.tar.gz";
            sha256 = "sha256-KFKvDLIKExObNEiZLmm4aOUO0Pih5ZQO4d6eGaEjthM=";
          };
          "aarch64-darwin" = {
            goFileName = "go1.25.0.darwin-arm64.tar.gz";
            sha256 = "sha256-4n9F10sLz2x2QWk0kSj+8J/J/1t9+D/P/g/5J/J/t9+";
          };
        };
      in
      {
        #  packages.go-nixpkgs-overridden = pkgs.go_1_24.overrideAttrs (old: {
        #  version = "1.24.4";
        #  src = pkgs.fetchurl {
        #    url = "https://golang.org/dl/go1.24.4.linux-amd64.tar.gz";
        #    sha256 = "sha256-d+XaM7tyrq7xukQYtv5RG8TQQYc8v4LlqmMYdA35hxc=";
        #  };
        # });

        packages.go-binary = pkgs.runCommand "go-binary-1.25.0" {
          nativeBuildInputs = [ pkgs.unzip pkgs.gzip pkgs.gnutar ];
          src = pkgs.fetchurl {
            url = "https://go.dev/dl/${goArchives.${system}.goFileName}";
            sha256 = goArchives.${system}.sha256;
          };
        } ''
          tar -xf $src --one-top-level
          mv go $out
        '';

        packages.default = self.packages.${system}.go-binary;
      });
}