# ns-process

Code to accompany the "Namespaces in Go" series of articles.

[Part 1: Linux Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf)
[Part 2: Namespaces in Go - Basics](https://medium.com/@teddyking/namespaces-in-go-basics-e3f0fc1ff69a)
[Part 3: Namespaces in Go - User](https://medium.com/@teddyking/namespaces-in-go-user-a54ef9476f2a)
[Part 4: Namespaces in Go - reexec](https://medium.com/@teddyking/namespaces-in-go-reexec-3d1295b91af8)
[Part 5: Namespaces in Go - Mount](https://medium.com/@teddyking/namespaces-in-go-mount-e4c04fe9fb29)
[Part 6: Namespaces in Go - Network](https://medium.com/@teddyking/namespaces-in-go-network-fdcf63e76100)
[Part 7: Namespaces in Go - UTS](https://medium.com/@teddyking/namespaces-in-go-uts-d47aebcdf00e)

## Usage

Each of the code extracts in the articles reference a git tag, which can be
checked out from this repo. The code is buildable and runnable at each tag, but
note that it will only run successfully on Linux machines.

## Testing

The test suite isn't explicitly mentioned in the articles, but if you'd like to
run the tests you'll need to install [ginkgo](https://github.com/onsi/ginkgo)
and [gomega](https://github.com/onsi/gomega).  Note that some of the tests may
require root privileges.
