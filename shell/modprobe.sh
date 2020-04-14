modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
if [[ $(uname -r |cut -d . -f1) -ge 4 && $(uname -r |cut -d . -f2) -ge 19 ]]; then
  modprobe -- nf_conntrack
else
  modprobe -- nf_conntrack_ipv4
fi
