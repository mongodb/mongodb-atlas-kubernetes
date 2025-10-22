{
  description = "A dev shell with a custom-fetched Go 1.25.3";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        goVersion = "1.25.3";

        go-src =
          if pkgs.stdenv.isLinux && pkgs.stdenv.hostPlatform.system == "x86_64-linux" then {
            url = "https://go.dev/dl/go${goVersion}.linux-amd64.tar.gz";
            sha256 = "sha256-AzXzFLbnv+CMPQz6p8GduWG3uZ+yC+YrCoJsmSrRTg8=";
          }
          else if pkgs.stdenv.isDarwin && pkgs.stdenv.hostPlatform.system == "aarch64-darwin" then {
            url = "https://go.dev/dl/go${goVersion}.darwin-arm64.tar.gz";
            sha256 = "sha256-fAg+PSwA3r/rL3fZpMAKGqyXETuJuczEKpBIevNDc4I=";
          }
          else throw "This flake does not support system: ${pkgs.stdenv.hostPlatform.system}";

        go_1_25_3 = pkgs.stdenv.mkDerivation {
          pname = "go-custom";
          version = goVersion;

          src = pkgs.fetchurl {
            inherit (go-src) url sha256;
          };

          dontBuild = true;

          installPhase = ''
            mkdir -p $out
            cp -a ./* $out/
          '';
        };

      in
      {
        packages.go_1_25_3 = go_1_25_3;
        packages.default = go_1_25_3;

        devShells.default = pkgs.mkShell {
          packages = [
            go_1_25_3
          ];
        };
      }
    );
}
