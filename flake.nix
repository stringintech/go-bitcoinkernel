{
  description = "Development shell for go-bitcoinkernel";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/b3d51a0365f6695e7dd5cdf3e180604530ed33b4";
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
          ninja
          pkg-config
          go
          python3
          git
        ];
        runtimeDeps = with pkgs; [
          boost
          libevent
          zlib
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = nativeDeps;
          buildInputs = runtimeDeps;
          shellHook = ''
            export PKG_CONFIG_PATH=${pkgs.lib.makeSearchPath "lib/pkgconfig" runtimeDeps}''${PKG_CONFIG_PATH:+:$PKG_CONFIG_PATH}
            export LD_LIBRARY_PATH=${pkgs.lib.makeLibraryPath runtimeDeps}''${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}
            export CGO_CFLAGS="-I${pkgs.lib.makeIncludePath runtimeDeps}''${CGO_CFLAGS:+ ''${CGO_CFLAGS}}"
            export CGO_LDFLAGS="-L${pkgs.lib.makeLibraryPath runtimeDeps}''${CGO_LDFLAGS:+ ''${CGO_LDFLAGS}}"
          '';
        };
      });
}
