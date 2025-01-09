{ pkgs ? import (fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/nixos-24.11.tar.gz";
  }) {
    config = {
      allowUnfree = true;
    };
  }
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    libGL
    xorg.libXi
    xorg.libXcursor
    xorg.libXrandr
    xorg.libXinerama
    wayland
    libxkbcommon
    jetbrains.goland
    go  # nixpkgs 24.11 ships go 1.23.3
    delve
    gdlv
  ];

  # this is needed for delve to work with cgo
  # see: https://wiki.nixos.org/wiki/Go#Using_cgo_on_NixOS
  hardeningDisable = [ "fortify" ];

  shellHook = '''';
}

