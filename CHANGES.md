# List a potential improvements

## Colors from config.toml

Ability to change colors definition from the `config.toml`

## Refactor what's displayed

### Ability to choose what's displayed from config.toml

default or option not set = all items displayed.
if option set in config.toml, match the selection

```toml
[display]
[display.options]
    "system.os" = true
    "system.uptime" = true
    "system.cpu" = true
    "system.ram" = true
    "system.ip" = true
    "system.gateway" = true
    "system.dns" = true
    "storage.root" = true
```
