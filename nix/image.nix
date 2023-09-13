{ pkgs ? (
    let
      sources = import ./sources.nix;
    in
    import sources.nixpkgs {
      overlays = [
        (import "${sources.gomod2nix}/overlay.nix")
      ];
    }
  )
, name ? "coco"
, tag ? null
}:
let
  coco = import ../. { inherit pkgs; };
in
pkgs.dockerTools.buildImage {
  inherit name tag;
  contents = [
    coco
    pkgs.poppler_utils
  ];
  config = {
    Cmd = [ "/bin/coco-backend" ];
  };
}
