.SILENT: release

pkgver=0.0.2

hashes:
	sha256sum potatoe
	sha256sum quotes.txt

srcinfo:
	cd aur && makepkg --printsrcinfo > .SRCINFO
