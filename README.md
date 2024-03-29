# go-ztcentral

Golang client library for interacting with the [ZeroTier Central Network Management Portal](https://my.zerotier.com)

NOTE: This does not work with self-hosted controllers.

Example:

```go
package main

import (
	"context"
	"log"
	"os"

	ztcentral "github.com/zerotier/go-ztcentral"
)

func main() {
	c, err := ztcentral.NewClient(os.Getenv("ZEROTIER_CENTRAL_TOKEN"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// get list of networks
	networks, err := c.GetNetworks(ctx)
	if err != nil {
		log.Println("error:", err.Error())
		os.Exit(1)
	}

	// print networks and members
	for _, n := range networks {
		log.Printf("%s\t%s", *n.Id, *n.Config.Name)
		members, err := c.GetMembers(ctx, *n.Id)
		if err != nil {
			log.Println("error:", err.Error())
			os.Exit(1)
		}

		for _, m := range members {
			log.Printf("\t%s\t %s", *m.Id, *m.Name)
		}
	}
}
```

# Development

Some useful make tasks:

- `make reflex-lint` and `make reflex-test` run the linters/testers with file watchers.
- `VERSION=x.y.z make release` - make a release with version x.y.z. Edits files and pushes tags.

# License

Copyright (c) 2021, ZeroTier, Inc.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
