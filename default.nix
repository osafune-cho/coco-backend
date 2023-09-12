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

pkgs.buildGoApplication {
  pname = "coco-backend";
  version = "0.1.0";
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
}
