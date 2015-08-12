%define app_home /opt/lesscms
%define app_user lesscms
%define app_grp  lesscms

Name: lesscms
Version: x.y.z
Release: 1%{?dist}
Vendor: lessOS.com
Summary: Productivity Tools for Enterprise
License: Apache 2
Group: Applications
Source0: lesscms-x.y.z.tar.gz
BuildRoot:  %{_tmppath}/%{name}-%{version}-%{release}

Requires:   redhat-lsb-core
Requires:   lesscms
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
install -d %{buildroot}%{app_home}/vendor/github.com/eryx/hcaptcha/var/fonts
install -d %{buildroot}%{app_home}/modules


cp -rp ./webui %{buildroot}%{app_home}/
cp -rp ./websrv %{buildroot}%{app_home}/


install -m 0755 -p bin/lesscms %{buildroot}%{app_home}/bin/lesscms
cp -rp etc/main.json.tpl %{buildroot}%{app_home}/etc/main.json
cp -rp modules %{buildroot}%{app_home}/
install -p -D -m 0755 misc/rhel/init.d-scripts %{buildroot}%{_initrddir}/lesscms
install -p -D -m 0755 vendor/github.com/eryx/hcaptcha/var/fonts/cmr10.ttf %{buildroot}%{app_home}/vendor/github.com/eryx/hcaptcha/var/fonts/cmr10.ttf


%clean
rm -rf %{buildroot}

%pre
# Add the "lesscms" user
getent group %{app_grp} >/dev/null || groupadd -r %{app_grp}
getent passwd %{app_user} >/dev/null || \
    useradd -r -g %{app_grp} -s /sbin/nologin \
    -d %{app_home} -c "lesscms user"  %{app_user}

if [ $1 == 2 ]; then
    service lesscms stop
fi


%post
# Register the lesscms service
if [ $1 -eq 1 ]; then
    /sbin/chkconfig --add lesscms
    /sbin/chkconfig --level 345 lesscms on
	/sbin/chkconfig lesscms on
fi


%preun
if [ $1 = 0 ]; then
    /sbin/service lesscms stop  > /dev/null 2>&1
    /sbin/chkconfig --del lesscms
fi


%postun

    
%files
%defattr(-,root,root,-)
%dir %{app_home}
%{_initrddir}/lesscms
%config(noreplace) %{app_home}/etc/main.json

%{app_home}/

