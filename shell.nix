{ pkgs ? (
    let
      sources = import ./nix/sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
}:
let
  goEnv = pkgs.mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    goEnv
    niv
    gomod2nix
    go
    gotools
    go-tools
    delve
    gopls
    nixpkgs-fmt
  ];

  GOROOT = "${pkgs.go}/share/go";
}
