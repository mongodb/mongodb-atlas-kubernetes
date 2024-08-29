{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        flake-utils.url = "github:numtide/flake-utils";
    };

    outputs = { self, nixpkgs, flake-utils }:
        flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-darwin" ] (system:
            let
                pkgs = nixpkgs.legacyPackages.${system};

                pname = "k8s-controller-tools";
                version = "0.16.1";

                src = pkgs.fetchFromGitHub {
                    owner = "kubernetes-sigs";
                    repo = "controller-tools";
                    rev = "v${version}";
                    sha256 = "0czfaaaxs7v7a7hy3cy83ph10r2lk4hlfn75whzp2njms5krvxh4";
                };

                k8s-controller-tools = pkgs.buildGoModule rec {
                    inherit pname version src;

                    patches = [ ./version.patch ];

                    vendorHash = "sha256-3p9K08WMqDRHHa9116//3lFeaMtRaipD4LyisaKWV7I=";

                    ldflags = [
                        "-s"
                        "-w"
                        "-X sigs.k8s.io/controller-tools/pkg/version.version=v${version}"
                    ];

                    doCheck = false;

                    subPackages = [
                        "cmd/controller-gen"
                        "cmd/type-scaffold"
                        "cmd/helpgen"
                    ];

                    meta = with pkgs.lib; {
                        description = "Tools to use with the Kubernetes controller-runtime libraries";
                        homepage = "https://github.com/kubernetes-sigs/controller-tools";
                        changelog = "https://github.com/kubernetes-sigs/controller-tools/releases/tag/v${version}";
                        license = licenses.asl20;
                        maintainers = with maintainers; [ michojel ];
                    };
                };
            in
            {
                packages.default = k8s-controller-tools;
            }
        );
}
