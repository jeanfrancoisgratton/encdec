<img src="./images/encdec_banner.png" alt="encdec logo" height="768" width="1024" />

# encdec
___
A small command-line tool to **encrypt/decrypt strings and files** using AES-256 (CFB mode).

> A note on terminology: throughout this project the words *encode/encrypt* and *decode/decrypt*
> are used loosely and interchangeably. They are technically different things — this is a
> deliberate, non-pedantic shortcut, not a mistake.

---

## Overview

`encdec` operates in two modes:

- **String mode** (default): encrypts/decrypts a string passed as an argument, printing the
  result to stdout. Ciphertext is emitted as URL-safe Base64.
- **File mode** (`-f`): encrypts/decrypts a file on disk. Files are streamed in 64&nbsp;KB chunks,
  so arbitrarily large files can be processed without loading them into memory.

Both modes use AES-256 in CFB mode with a randomly generated IV that is stored alongside the
ciphertext (prepended to the file, or embedded in the Base64 output for strings).

## Usage

```
encdec {encode|decode} [flags] <string | file> [destfile]
```

Subcommands (with aliases):

| Command             | Alias | Description                          |
|---------------------|-------|--------------------------------------|
| `encode`            | `enc` | Encrypt a string or file             |
| `decode`            | `dec` | Decrypt a string or file             |
| `changelog`         | `cl`  | Show the version changelog           |

Flags:

| Flag              | Scope           | Description                                                        |
|-------------------|-----------------|-------------------------------------------------------------------|
| `-s, --secret`    | global          | 32-byte encryption/decryption key (**must be exactly 32 bytes**)  |
| `-q, --quiet`     | global          | Print only the resulting string, no decoration                    |
| `--debug`         | global          | Show extra debug output                                           |
| `-f, --file`      | encode/decode   | Operate on a file instead of a string                             |
| `-k, --keep`      | encode/decode   | Keep the original file instead of replacing it in place           |

### String mode

```sh
encdec encode "some secret text"
encdec decode AAECAwQF...     # the Base64 blob produced above
```

The decoded/encoded string is written to stdout. Use `-q` to print just the raw value (handy
for piping into other commands).

### File mode (requires `-f`)

```sh
encdec encode -f secrets.txt              # encrypts secrets.txt in place
encdec encode -f secrets.txt out.enc      # writes ciphertext to out.enc
encdec decode -f secrets.txt              # decrypts secrets.txt in place
```

Behaviour with respect to the destination file:

- If **no destination** is given, encdec writes to a temporary file (`<source>.enc` /
  `<source>.dec`) and, unless `-k` is passed, removes the original and renames the result back
  to the original name — i.e. the file is transformed **in place**.
- If a **destination** is given, the result is written there.
- Passing `-k, --keep` preserves the original source file.

### The secret key

AES-256 requires a 32-byte key. `encdec` ships with a **hard-coded default key** that is trivially
visible in the source — it exists purely for testing and demos. **Do not rely on it for anything
sensitive.**

To use your own key, pass it with `-s`:

```sh
encdec encode -s "0123456789abcdef0123456789abcdef" "top secret"
```

The key **must be exactly 32 bytes long**, otherwise the operation aborts. There is no key-recovery
mechanism: if you lose the key you used to encrypt something, the data is unrecoverable, so store
it safely.

## Build & Install

You can build from source or grab a pre-built package.

### From source

Requirements: a Go toolchain (see `src/go.mod` for the version) and Git.

```sh
git clone https://github.com/jeanfrancoisgratton/encdec.git
cd encdec/src

./updateBuildDeps.sh   # optional: refresh module dependencies
./build.sh             # builds and installs the binary
```

By default `build.sh` compiles a stripped, trimmed static binary
(`CGO_ENABLED=0`, `-ldflags="-s -w"`) into `/opt/bin`. You can override the destination:

```sh
./build.sh /usr/local/bin      # install elsewhere
./build.sh -b mytool /opt/bin  # override the binary name
```

### Binary packages

Packaging scripts and specs are provided for several distributions:

- **Alpine** — `__alpine/` (APKBUILD)
- **Arch Linux** — `__archlinux/` (PKGBUILD)
- **Debian** — `__debian/`
- **RedHat / RPM** — `__redhat/`

These build recipes depend on the author's own build containers and are not guaranteed to work
out of the box elsewhere. As an easier alternative, pre-built packages are published at:

<https://github.com/jeanfrancoisgratton/encdec/releases>

## License

Released under the **GNU General Public License v3** — see [docs/LICENSE](docs/LICENSE).

## Author

Jean-François Gratton — <jean-francois@famillegratton.net>
