# Setups routing for Internet on a SDN network (Proxmox)

SDN_IP_RANGE='10.0.3.0/24'

echo "Removing old up/down scripts..."
rm /etc/network/if-up.d/enable_sdn_internet.sh
rm /etc/network/if-post-down.d/disable_sdn_internet.sh

echo "Creating new up/down scripts..."

echo '#!/bin/sh
# Adds routing when the SDN network goes up
if [ "$IFACE" = "myvnet1" ]; then
        echo 1 > /proc/sys/net/ipv4/ip_forward
        iptables -t nat -A POSTROUTING -s '\'"$SDN_IP_RANGE"\'' -o vmbr0 -j MASQUERADE
        iptables -t raw -I PREROUTING -i fwbr+ -j CT --zone 1
fi' >/etc/network/if-up.d/enable_sdn_internet.sh

echo '#!/bin/sh
# Removes routing when the SDN network goes down
if [ "$IFACE" = "myvnet1" ]; then
        iptables -t nat -D POSTROUTING -s '\'"$SDN_IP_RANGE"\'' -o vmbr0 -j MASQUERADE
        iptables -t raw -D PREROUTING -i fwbr+ -j CT --zone 1
fi
' >/etc/network/if-post-down.d/disable_sdn_internet.sh

echo "Making the scripts executable..."

chmod -R 777 /etc/network/if-up.d/enable_sdn_internet.sh /etc/network/if-post-down.d/disable_sdn_internet.sh

echo "Done!"
