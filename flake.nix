{
  description = "go-bitcoinkernel";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        nativeDeps = with pkgs; [
          cmake
          go_1_23
          golangci-lint
        ];
        runtimeDeps = with pkgs; [
          boost
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = nativeDeps;
          buildInputs = runtimeDeps;
        };
      });
}
