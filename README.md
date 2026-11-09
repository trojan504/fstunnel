## Installation

One command, on **both** servers:

```bash
sudo bash -c "$(curl -sL https://raw.githubusercontent.com/CIAUB/fstunnel/main/install.sh)"
```

The installer always installs the latest [release](https://github.com/CIAUB/fstunnel/releases/latest). It downloads the binaries for the server's architecture, the kernel helper and the systemd units, checks each file against the release's `SHA256SUMS`, and installs only if every hash matches. A corrupt or modified download stops the install and nothing is changed.

```bash
sudo fstunnel
```

To update, choose **Update** in the menu or run the install command again. Running tunnels are kept. Without a version, the installer selects the latest stable GitHub release and excludes drafts and pre-releases. To install a specific version:

```bash
sudo bash -c "$(curl -sL https://raw.githubusercontent.com/CIAUB/fstunnel/main/install.sh)" v0.2.1
```

`bash install.sh v0.2.1` and `FSTUNNEL_VERSION=v0.2.1` are also supported. An explicit version argument takes precedence over the environment variable.
