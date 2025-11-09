# List a potential improvments

## Colors from config.toml

Ability to change colors definition from the `config.toml`

## System OS info

Get an OS info label from /etc/os-release in Sytem title as the first label
```
System
  > OS      {ID} {VERSION}

System
  > OS      debian 13 (trixie)
```

## Service display improvment

Change the `status` content and rendering.

Instead of displaying the whole line of the "Active:" (`status`) on new line:

- display a shorter string on same line (do not apply label_width as service name can be long)
- same color coded of `active (running)`, `inactive (dead)`, ... as the `marker`
- dimmed relative time. Content after `;`.

```
# From
● xe-linux-distribution
    active running since Sun 2025-11-09 09:38:17 CET; 6min ago
○ xo-server
    inactive dead
# To
● xe-linux-distribution: active (running) since 6min ago
○ xo-server: inactive (dead)
```

## Refactor what's displayed

### Compact System & Network

```
System
  > OS        {os_string}
  > Uptime    {uptime_string}
  > CPU       {cpu_cores} ({cpu_usage}%)
  > RAM       {current_ram}/{total_ram} ({ram_usage}%)
  > IP        {ip_addr}/{cidr_subnetmask} ({interface})
  > Gateway   {gateway}
  > DNS       {dns}
```

### Ability to choose what's displayed from config.toml

default or option not set = all items displayed.
if option set in config.toml, match the selection

```toml
[display]
    [options]
    "system.os": true
    "system.uptime": true
    "system.cpu": true
    "system.ram": true
    "system.ip": true
    "system.gateway" :true
    "system.dns" :true
    "storage.root": true
```
