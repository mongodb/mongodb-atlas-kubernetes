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
        version = "1.24.7";
      in
      {
        packages.default = pkgs.go_1_24.overrideAttrs (old: {
          pname = "go";
          inherit version;
          src = pkgs.fetchurl {
            url = "https://golang.org/dl/go${version}.linux-amd64.tar.gz";
            sha256 = "sha256-2hgZHdt9uKkzmBbz4rVL3e2AR83CpdZwWUePjRWVxD8=";
          };
        });
      });
}
