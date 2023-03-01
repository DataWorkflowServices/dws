%undefine _missing_build_ids_terminate_build
%global debug_package %{nil}

Name: dws-clientmount
Version: 1.0
Release: 1%{?dist}
Summary: Client mount daemon for data workflow service

Group: 1
License: Apache-2.0
URL: https://github.com/HewlettPackard/dws
Source0: %{name}-%{version}.tar.gz

BuildRequires:	golang
BuildRequires:	make

%description
This package provides clientmountd for performing mount operations for the
data workflow service

%prep
%setup -q

%build
COMMIT_HASH=$(cat .commit) make build-daemon

%install
mkdir -p %{buildroot}/usr/bin/
install -m 755 bin/clientmountd %{buildroot}/usr/bin/clientmountd

%files
/usr/bin/clientmountd
