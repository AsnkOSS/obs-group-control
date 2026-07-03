{ pkgs, ... }:

{
  packages = with pkgs; [
    pkg-config
    libGL
    libx11
    libxcursor
    libxrandr
    libxinerama
    libxi
    libxxf86vm
  ];

  languages = {
    go = {
      enable = true;
      version = "1.26.0";
    };
    javascript = {
      enable = true;
      package = pkgs.nodejs_24;
      corepack = {
        enable = true;
      };
    };
  };
}
