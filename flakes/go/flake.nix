{
  description = "A flake for Go 1.24.4";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.default = pkgs.go_1_24.overrideAttrs (old: {
          version = "1.24.4";
          src = pkgs.fetchurl {
            url = "https://golang.org/dl/go1.24.4.linux-amd64.tar.gz";
            sha256 = "sha256-d+XaM7tyrq7xukQYtv5RG8TQQYc8v4LlqmMYdA35hxc=";
          };
        });
      });
}
