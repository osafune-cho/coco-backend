{ pkgs ? import ./nix/pkgs.nix }:
pkgs.buildGoModule {
  pname = "coco-backend";
  version = "0.1.0";

  src = ./.;

  vendorHash = null;
}
