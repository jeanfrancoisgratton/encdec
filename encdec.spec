%ifarch aarch64
%global _arch aarch64
%global BuildArchitectures aarch64
%endif

%ifarch x86_64
%global _arch x86_64
%global BuildArchitectures x86_64
%endif

%define debug_package   %{nil}
%define _build_id_links none
%define _name   encdec
%define _prefix /opt
%define _version 1.30.00
%define _rel 0
#%define _arch aarch64
%define _binaryname encdec

Name:       encdec
Version:    %{_version}
Release:    %{_rel}
Summary:    encdec

Group:      SSL
License:    GPL2.0
URL:        https://github.com/jeanfrancoisgratton/encdec

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: aarch64
BuildRequires: gcc

%description
Encode and decode AES-256 strings and files

%prep
%autosetup

%build
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o %{_sourcedir}/%{_binaryname} .

%clean
rm -rf $RPM_BUILD_ROOT

%pre


%install
install -Dpm 0755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post
BIN="%{_prefix}/bin/%{_binaryname}"

%preun

%postun

%files
%defattr(-,root,root,-)
%attr(2755,root,root) %{_prefix}/bin/%{_binaryname}



%changelog
* Mon Nov 17 2025 Binary package builder <builder@famillegratton.net> 1.30.00-0
- Completed proper error handling (jean-francois@famillegratton.net)
- Fixed package name in abuild script (builder@famillegratton.net)
- apk fix (jean-francois@famillegratton.net)
- interim commit (jean-francois@famillegratton.net)
- Version bump (jean-francois@famillegratton.net)
- removed github actions file, un-needed (jean-francois@famillegratton.net)

* Tue Sep 09 2025 Binary package builder <builder@famillegratton.net> 1.21.03-0
- new package built with tito

* Thu Dec 19 2024 APK Builder <builder@famillegratton.net> 1.21.03-0
- GO version bump (jean-francois@famillegratton.net)

* Sat Oct 19 2024 RPM Builder <builder@famillegratton.net> 1.21.02-0
- re-initialized tito

* Tue Aug 13 2024 RPM Builder <builder@famillegratton.net> 1.21.02-0
- global variable re-scoping (jean-francois@famillegratton.net)

* Mon Aug 12 2024 RPM Builder <builder@famillegratton.net> 1.21.01-0
- version bump, and tag fix in GHA (jean-francois@famillegratton.net)

* Mon Aug 12 2024 RPM Builder <builder@famillegratton.net> 1.21.00-0
- inverted the -q switch (jean-francois@famillegratton.net)

* Mon Aug 12 2024 RPM Builder <builder@famillegratton.net> 1.20.01-2
  francois@famillegratton.net)
- Go and package version bump (jean-francois@famillegratton.net)

* Sun Aug 11 2024 RPM Builder <builder@famillegratton.net> 1.20.01-1
- Make the specfile more arch-independant (jean-francois@famillegratton.net)

* Sun Aug 11 2024 RPM Builder <builder@famillegratton.net> 1.20.01-0
- Changed arch moniker on aarch64 (jean-francois@famillegratton.net)
- Fixed -q issue (jean-francois@famillegratton.net)

* Fri Aug 09 2024 RPM Builder <builder@famillegratton.net> 1.20.00-0
- Better file handling, gha added, go version bump (jean-
  francois@famillegratton.net)

* Tue Jun 25 2024 RPM Builder <builder@famillegratton.net> 1.10.00-0
- Added a -q switch (jean-francois@famillegratton.net)
- Updated to GO 1.22.4 (jean-francois@famillegratton.net)

* Mon Nov 06 2023 RPM Builder <builder@famillegratton.net> 1.02.00-1
- Go version bump (jean-francois@famillegratton.net)

* Mon Nov 06 2023 RPM Builder <builder@famillegratton.net> 1.02.00-0
- Fixed argcount error, version numbering scheme change (jean-
  francois@famillegratton.net)
- Fixed chown issue in packaging (jean-francois@famillegratton.net)

* Wed Aug 02 2023 builder <builder@famillegratton.net> 1.000-1
- Updated changelog and some forgotten relase numbers in packaging scripts
  (jean-francois@famillegratton.net)

* Wed Aug 02 2023 builder <builder@famillegratton.net> 1.000-0
- Prod-ready release: 1.000-0 (jean-francois@famillegratton.net)

* Mon Jul 31 2023 builder <builder@famillegratton.net> 0.200-1
- Doc and version updates (jean-francois@famillegratton.net)

* Mon Jul 31 2023 builder <builder@famillegratton.net> 0.200-0
- Go version bump, also forgotten in previous commits (builder@famillegratton.net)
- Version bump (forgotten in previous commit) (jean-
  francois@famillegratton.net)
- Completed 0.200 (jean-francois@famillegratton.net)
- File Dec-Enc facilities, stub 1 (jean-francois@famillegratton.net)
- New brandh to fully encode/decode files (jean-francois@famillegratton.net)

* Mon Jul 10 2023 builder <builder@famillegratton.net> 0.100-0
- new package built with tito

