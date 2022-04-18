# Contributor: Hennik Hunsaker <hennikhunsaker@microbox.cloud>
# Maintainer: Hennik Hunsaker <hennikhunsaker@microbox.cloud>
pkgname=slurp
pkgver=0.2.3
pkgrel=0
pkgdesc="A simple, api-driven storage system for storing code builds and cached libraries for cloud-based deployment services."
url="https://github.com/mu-box/slurp"
arch="all"
license="MIT"
depends=""
makedepends="go git bash"
checkdepends=""
install=""
subpackages=""
source=""
srcdir="/tmp/abuild/slurp"
builddir=""

build() {
	go get -t -v .
	go install github.com/mitchellh/gox@latest
	export PATH="$(go env | grep GOPATH | sed -E 's/GOPATH="(.*)"/\1/')/bin:${PATH}"
	./scripts/build.sh
}

check() {
	# Replace with proper check command(s)
	:
}

package() {
	install -m 0755 -D ./build/slurp "$pkgdir"/sbin/slurp
}
