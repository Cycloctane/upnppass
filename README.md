# upnpPass

UPnP Pass: Pass UPnP/DLNA discovery messages through subnet/vpn/port-forwarding.

This program acts like a [SSDP Server Proxy](https://datatracker.ietf.org/doc/html/draft-cai-ssdp-v1-01#section-7.2) for UPnP devices. It retrieves description from UPnP device and announces SSDP messages in local network to make the target UPnP device visible in local subnet.

It only advertises for the target. Location in SSDP packet remains the same as UPnP device's orginal address. UPnP clients (control points) should be able to connect to the target UPnP device directly.

It can also be used for accessing remote UPnP/DLNA service through port forwarding or vpn that do not route multicast traffic (like ipsec and openvpn).

Currently supports UPnP root devices and services. Proxy for embedded devices is not implemented yet.

## Usage

```bash
./upnppass -i $interface -u $description_url -t 1800
```

- `-i`: Network interface used for sending SSDP multicast messages.
- `-u`: URL of target UPnP device's root device description. (`http://host:8200/rootDesc.xml` for minidlna)
- `-t`: Advertisement duration (max-age).

For example, to make a minidlna server in remote network visible to localhost after forwarding remote minidlna server's 8200 port to local 127.0.0.1:

```bash
./upnppass -i loop -u http://127.0.0.1:8200/rootDesc.xml -t 1800
```
