## Installation

One command, on **both** servers:

```bash
sudo bash -c "$(curl -sL https://raw.githubusercontent.com/trojan504/fstunnel/main/install.sh)"
```

The installer always installs the latest [release](https://github.com/trojan504/fstunnel/releases/latest). It downloads the binaries for the server's architecture, the kernel helper and the systemd units, checks each file against the release's `SHA256SUMS`, and installs only if every hash matches. A corrupt or modified download stops the install and nothing is changed.

```bash
sudo fstunnel
```
