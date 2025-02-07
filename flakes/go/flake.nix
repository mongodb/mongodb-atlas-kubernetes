{
  description = "A flake for Go 1.23.6";
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
        packages.default = pkgs.go_1_23.overrideAttrs (old: {
          version = "1.23.6";
          src = pkgs.fetchurl {
            url = "https://golang.org/dl/go1.23.6.linux-amd64.tar.gz";
            sha256 = "sha256-k3lEHqMQ3gAPM6TcdnvZZucqsoJicOA454ssU8LngC0=";
          };
        });
      });
}
