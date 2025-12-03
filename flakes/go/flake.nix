{
  description = "A dev shell with a custom-fetched Go 1.25.5";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        goVersion = "1.25.5";

        go-src =
          if pkgs.stdenv.isLinux && pkgs.stdenv.hostPlatform.system == "x86_64-linux" then {
            url = "https://go.dev/dl/go${goVersion}.linux-amd64.tar.gz";
            sha256 = "sha256-npt1XWOzas8wwSqaP8N5JDcUwcbT3XKGHaY38zbrs1s=";
          }
          else if pkgs.stdenv.isDarwin && pkgs.stdenv.hostPlatform.system == "aarch64-darwin" then {
            url = "https://go.dev/dl/go${goVersion}.darwin-arm64.tar.gz";
            sha256 = "sha256-vtjr6CTj07J+hHHRMH+AP8arjh0Ot6SuGWl5vZuAHdM=";
          }
          else throw "This flake does not support system: ${pkgs.stdenv.hostPlatform.system}";

        go_1_25_ako = pkgs.stdenv.mkDerivation {
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
        packages.go_1_25 = go_1_25_ako;
        packages.default = go_1_25_ako;

        devShells.default = pkgs.mkShell {
          packages = [
            go_1_25_ako
          ];
        };
      }
    );
}
