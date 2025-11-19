{
  description = "A simple tool to sync Readwise highlights to org files";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs@{ self, nixpkgs, flake-parts }:
    let
      buildPackage = { buildGoModule, lib }:
        buildGoModule {
          pname = "go-org-readwise";
          version = self.rev or "dirty";

          src = ./.;

          vendorHash = null;

          meta = with lib; {
            description = "Sync Readwise highlights to org files";
            homepage = "https://github.com/vdemeester/go-org-readwise";
            license = licenses.mit;
            maintainers = [ ];
          };
        };
    in
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, pkgs, system, ... }: {
        packages = {
          default = pkgs.callPackage buildPackage { };
          go-org-readwise = config.packages.default;
        };

        apps = {
          default = {
            type = "app";
            program = "${config.packages.default}/bin/go-org-readwise";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ];
        };
      };

      flake = {
        overlays.default = final: prev: {
          go-org-readwise = final.callPackage buildPackage { };
        };
      };
    };
}
