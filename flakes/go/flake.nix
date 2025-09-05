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
          version = "1.24.6";
          src = pkgs.fetchurl {
            url = "https://golang.org/dl/go1.24.6.linux-amd64.tar.gz";
            sha256 = "sha256-u8o3zDlcl0/6SJPuNYGa0j67J0Jt+Hr5LpOp7GbvhxI=";
          };
        });
      });
}
