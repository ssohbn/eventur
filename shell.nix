{ pkgs ? import <nixpkgs> { system = "x86_64-linux"; }}:
pkgs.mkShell {
	packages = with pkgs; [
		go
		screen
	];
}
