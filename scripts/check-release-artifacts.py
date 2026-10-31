"""Verify the release's supported archives, contents, and SHA-256 checksums."""

import hashlib
from pathlib import Path
import tarfile
import zipfile


def main():
    dist = Path("dist")
    targets = {
        "darwin_arm64", "darwin_x86_64", "freebsd_i386", "freebsd_x86_64",
        "linux_arm64", "linux_armv6", "linux_armv7", "linux_i386", "linux_x86_64",
        "windows_arm64", "windows_i386", "windows_x86_64",
    }
    seen = set()
    checksums = (dist / "checksums.txt").read_text().splitlines()
    for line in checksums:
        digest, name = line.split()
        path = dist / name
        if hashlib.sha256(path.read_bytes()).hexdigest() != digest:
            raise ValueError(f"Checksum mismatch: {name}")
        target = next((t for t in targets if name.endswith(
            "_" + t + (".zip" if t.startswith("windows_") else ".tar.gz")
        )), None)
        if target is None or target in seen or not name.startswith("killgrave_"):
            raise ValueError(f"Unexpected archive: {name}")
        seen.add(target)
        binary = "killgrave.exe" if target.startswith("windows_") else "killgrave"
        if path.suffix == ".zip":
            with zipfile.ZipFile(path) as archive:
                entries = set(archive.namelist())
        else:
            with tarfile.open(path) as archive:
                entries = set(archive.getnames())
        if entries != {binary, "LICENSE", "README.md", "CHANGELOG.md"}:
            raise ValueError(f"Unexpected contents in {name}: {entries}")
    if seen != targets:
        raise ValueError(f"Missing targets: {targets - seen}")
    print(f"Verified {len(seen)} release archives and checksums")


if __name__ == "__main__":
    main()
