Summary:    Group deduplication.
Name:       groupcover
Version:    0.1.0
Release:    0
License:    GPL
BuildArch:  x86_64
BuildRoot:  %{_tmppath}/%{name}-build
Group:      System/Base
Vendor:     Leipzig University Library, https://www.ub.uni-leipzig.de
URL:        https://github.com/miku/groupcover

%description

Group deduplication.

%prep

%build

%pre

%install
mkdir -p $RPM_BUILD_ROOT/usr/local/sbin
install -m 755 groupcover $RPM_BUILD_ROOT/usr/local/sbin

mkdir -p $RPM_BUILD_ROOT/usr/local/share/man/man1
install -m 644 span.1 $RPM_BUILD_ROOT/usr/local/share/man/man1/groupcover.1

%post

%clean
rm -rf $RPM_BUILD_ROOT
rm -rf %{_tmppath}/%{name}
rm -rf %{_topdir}/BUILD/%{name}

%files
%defattr(-,root,root)

/usr/local/sbin/groupcover
/usr/local/share/man/man1/groupcover.1

%changelog
* Sat Dec 21 2016 Martin Czygan
- 0.0.1 initial preview
