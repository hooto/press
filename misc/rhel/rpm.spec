%define app_home /opt/hooto/hootopress
%define app_user htap
%define app_grp  htap

Name: hootopress
Version: x.y.z
Release: 1%{?dist}
Vendor: hooto.com
Summary: Productivity Tools for Enterprise
License: Apache 2
Group: Applications
Source0: hootopress-x.y.z.tar.gz
BuildRoot:  %{_tmppath}/%{name}-%{version}-%{release}

Requires:       redhat-lsb-core,gd
Requires(pre):  shadow-utils
Requires(post): chkconfig

%description
%prep
%setup  -q -n %{name}-%{version}
%build

%install
rm -rf %{buildroot}
install -d %{buildroot}%{app_home}/bin
install -d %{buildroot}%{app_home}/etc
install -d %{buildroot}%{app_home}/var


cp -rp ./webui %{buildroot}%{app_home}/
cp -rp ./websrv %{buildroot}%{app_home}/
cp -rp ./modules %{buildroot}%{app_home}/

install -m 0755 -p bin/keeper %{buildroot}%{app_home}/bin/keeper
install -m 0755 -p bin/hootopress %{buildroot}%{app_home}/bin/hootopress
install -m 0640 -p etc/main.json.tpl %{buildroot}%{app_home}/etc/main.json


%clean
rm -rf %{buildroot}

%pre
# Add the "hootopress" user
getent group %{app_grp} >/dev/null || groupadd -r %{app_grp}
getent passwd %{app_user} >/dev/null || \
    useradd -r -g %{app_grp} -s /sbin/nologin \
    -d %{app_home} -c "hootopress user"  %{app_user}

%post

%preun

%postun

    
%files
%defattr(-,htap,htap,-)
%dir %{app_home}
%config(noreplace) %{app_home}/etc/main.json

%{app_home}/

