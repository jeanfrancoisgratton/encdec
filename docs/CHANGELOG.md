# Changelog

All notable changes to **encdec** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this
project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) as of 1.4.0.
Earlier releases used a `major.minor.patch` scheme with zero-padded fields (`1.32.00`,
`1.21.03`, ...), kept here verbatim for traceability.

Entries up to and including 1.21.03 were reconstructed from the RPM `%changelog` in
`__redhat/encdec.spec`; the current git repository only goes back to 2025-09-09.

## [1.5.0] — 2026-08-12

### Changed
- Encryption and decryption are now delegated to the `helperFunctions` package instead of
  being implemented in-tree; `src/executor/helpers.go` was dropped in favor of the shared
  implementation.
- `-s` / `--secret` is now an *optional* passphrase of any length (defaulting to the empty
  passphrase). It is no longer a mandatory 32-byte key.
- Naming a destination file no longer overwrites the source file with it.
- `-k` / `--keep` now applies to in-place runs only, i.e. when no destination file is given.
- Packaging is driven by `Makefile`s in each of `__alpine/`, `__archlinux/`, `__debian/`
  and `__redhat/`; the numbered helper shell scripts were retired.

### Added
- `-F` / `--force`, to overwrite an existing destination file.
- Expanded unit tests for both file and string encrypt/decrypt paths.

### Removed
- `__redhat/PACKAGE_RPM.md`, superseded by the RedHat `Makefile`.

## [1.4.1] — 2026-07-13

### Changed
- Secret key management moved from the interactive `-p` flag to `-s`; `-p` was removed.
- README rewritten and expanded; `LICENSE` and `TODO.md` moved under `docs/`.

### Added
- Test files: `src/executor/filecrypt_test.go` and `src/executor/stringcrypt_test.go`.
- Project banner image (`images/encdec_banner.png`).

## [1.4.0] — 2026-07-12

### Changed
- Version numbering fully aligned with SemVer (`1.32.00` → `1.4.0`).
- Go version bump to 1.26.5.
- rpmbuild enhancements in `__redhat/Makefile`.

### Removed
- Last remnants of tito (`.tito/`), left over from the previous build system.

## [1.32.00] — 2026-05-12

### Changed
- Binary packaging structure overhauled for RHEL and Arch Linux; the spec file moved to
  `__redhat/encdec.spec` and a `Makefile` plus `updateChangelog.sh` were added.
- Go version bump to 1.26.3.

### Added
- Arch Linux packaging: `__archlinux/PKGBUILD` and its build-dependency scripts.

## [1.31.00] — 2026-04-14

### Fixed
- Decoding a string actually returned a re-encoded string instead of the decoded one.
- Version output, including the git commit hash.
- Unnecessary use of `Printf()`.

### Changed
- Error handling made consistent across the executor functions.
- Go version bump to 1.26.2, then downgraded to 1.26.1 because Alpine did not yet support
  1.26.2.
- Moved to `helperfunctions` v5; `src/executor/structs.go` renamed to `types.go`.

## [1.30.00-1] — 2025-11-17

### Changed
- Build dependencies updated; release number bumped accordingly.

## [1.30.00] — 2025-11-17

### Added
- A forgotten error path in the file encryption routine.

### Changed
- Proper error handling completed throughout.
- Go version bump to 1.25.4; major package and build-dependency update.
- `src/executor/crypto.go` renamed to `helpers.go`.

### Fixed
- Package name and other issues in the Alpine `APKBUILD` script.

### Removed
- The GitHub Actions release workflow, no longer needed.

## [1.21.03] — 2024-12-19

### Changed
- Go version bump to 1.23.4.

## [1.21.02] — 2024-08-13

### Changed
- Global variable re-scoping / variables reshuffling.

## [1.21.01] — 2024-08-12

### Changed
- Version bump, and tag fix in the GitHub Actions workflow.

## [1.21.00] — 2024-08-12

### Changed
- Inverted the quiet/verbose `-q` switch.

## [1.20.01] — 2024-08-11

### Fixed
- `-q` handling.
- Arch moniker on aarch64; the spec file was made more arch-independent.

## [1.20.00] — 2024-08-09

### Changed
- Better file handling for the destination file.
- Go version bump.

### Added
- GitHub Actions workflow.

## [1.10.00] — 2024-06-25

### Added
- `-q` switch.

### Changed
- Moved to the GitHub `helperFunctions` package.
- Updated to Go 1.22.4.

## [1.02.00] — 2023-11-06

### Fixed
- Argument count error.
- chown issue in packaging.

### Changed
- Version numbering scheme changed.
- Go version bump.

## [1.000] — 2023-08-02

### Changed
- First production-ready release.
- Changelog and forgotten release numbers in the packaging scripts updated.

## [0.200] — 2023-07-31

### Added
- File encryption/decryption capabilities.

### Changed
- Documentation and version updates; Go version bump.

## [0.100] — 2023-07-10

### Added
- Initial stub release.

[1.5.0]: https://github.com/jeanfrancoisgratton/encdec/releases/tag/v1.5.0
[1.4.1]: https://github.com/jeanfrancoisgratton/encdec/releases/tag/v1.4.1
[1.4.0]: https://github.com/jeanfrancoisgratton/encdec/releases/tag/v1.4.0
