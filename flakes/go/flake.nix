{
  description = "A multi-platform Go development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-darwin" ];

      perSystem = nixpkgs.lib.genAttrs supportedSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };

          go-official-darwin = pkgs.stdenv.mkDerivation {
            pname = "go-official-darwin";
            version = "1.25.0";
            src = pkgs.fetchurl {
              url = "https://go.dev/dl/go1.25.0.darwin-arm64.tar.gz";
              outputHashAlgo = "sha256";
              outputHash = "sha256-VEkyhEFW2Bcveij3fyrJwVojBGaYtiQ/YzsKCwDAdJw=";
            };
            installPhase = ''
              mkdir -p $out
              cp -r * $out/
            '';
          };

          go-toolchain = if pkgs.stdenv.isDarwin then go-official-darwin else pkgs.go;
          cgo-flags = if pkgs.stdenv.isDarwin then "-framework CoreFoundation -framework Security" else "";

        in
        {
          # A minimal "dummy" package that just provides the go toolchain.
          # This satisfies Devbox's requirement for a package to exist.
          pkg = pkgs.runCommand "go-toolchain-env" {
            nativeBuildInputs = [ go-toolchain ];
          } ''
            mkdir -p $out/bin
            ln -s ${go-toolchain}/bin/go $out/bin/go
          '';

          # The actual development shell you want to use.
          shell = pkgs.mkShell {
            nativeBuildInputs = [ go-toolchain ];
            CGO_LDFLAGS = cgo-flags;
            shellHook = ''
              echo "--- Using Go toolchain for ${system} ---"
              go version
            '';
          };
        });
    in
    {
      # This structure correctly provides both outputs
      devShells = nixpkgs.lib.mapAttrs (system: attrs: attrs.shell) perSystem;
      packages = nixpkgs.lib.mapAttrs (system: attrs: { default = attrs.pkg; }) perSystem;
    };
}