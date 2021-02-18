#!bash

extra=

if [ "$(git rev-parse `cat VERSION`)" != "$(git rev-parse HEAD)" ]
then
  extra="-$(git rev-parse HEAD)"
fi

cat >version.go <<EOF
package ztcentral

// Version is the version of this library.
const Version = "$(cat VERSION)${extra}"
EOF

