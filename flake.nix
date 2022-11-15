{
  description = "sshign: simple command line tool to verify ssh signatures";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in {
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          sshign = pkgs.buildGoModule {
            pname = "sshign";
            version = "v0.0.1";
            src = ./.;

            vendorSha256 = "sha256-CY8uKFkxC9Xu3QaxYXbBuwzb2EjzFbAga3MS18KsmS4=";
            proxyVendor = true;
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.sshign);
      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            shellHook = ''
              PS1='\u@\h:\@; '
              echo "Go `${pkgs.go}/bin/go version`"
            '';
            nativeBuildInputs = with pkgs; [ git go gopls go-tools ];
          };
        });
    };
}

