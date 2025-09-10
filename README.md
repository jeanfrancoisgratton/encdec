<H1>encdec</H1>
Encrypt and Decrypt strings<br><br>
___
(while technically different, I have a loose use of terms like encrypt-encode; I'm not pedantic, and aware that they are different; bear with me, please :-) )
<H2>Overview</H2>
<H3>How does it work</H3>
In its default mode, this tool will encrypt and decrypt any string you will pass as an argument.<br>
With the `-f` parameter, you pass a source and destination filename, and it will AES-256 encode/decode the file.
<H3>String mode</H3>
Simple: `encdec {encode|decode} $STRING`<br>
Where $STRING is the string to encode or decode. Note that the decoded string will be output to stdout.

<H3>File mode (requires the -f flag)</H3>
Also simple: `encdec {encode|decode} $FILENAME`
You need to provide the full source pathname.

<H3>Private key</H3>
AES-256 needs a private key to encrypt and decrypt a file or string. There is a built-in key in the software that you can easily find browsing the source code.<br>I left it there for testing purposes, but for obvious security reasons, you should avoid using it.<br><br>
If you want to use your own, the `-p` flag will prompt you eventually to provide your own key. You will need to safe-keep it somewhere, as there are no mechanisms to retreive it if lost. Also, note that the key _needs to be exactly 32bytes long_.

<H2>Build, Install</H2>
You can either build from source, or use the provided binary packages (Alpine, RedHat, Debian)

<H3>Source install</H3>
Once you've cloned the repo, go into the `src/` directory, and run `./upgrade_pkgs.sh` (about this, please note that package upgrades have not been tested beyond current versions, for obvious reasons).

After this step, you can run `./build.sh` (providing an optional `-o` param with a target dir is available). This will build and copy the binary into `/opt/bin`

You might wish to `strip` the produced binary, as it still contains debugging code.

<H3>Binary packages</H3>
For now, I provide all the necessary scripts and files to build the RPM, DEB and APK packages, but all of these depend on my own build containers, which, for now, I do not provide. Maybe some other time, but not now as they depend way too heavily on my home architecture.

So, as an alternative, you can find the binary packages at
https://github.com/jeanfrancoisgratton/encdec/releases
