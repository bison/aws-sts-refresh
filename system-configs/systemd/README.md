# Installation

## Create "User" Unit/Timer files

Put these files wherever they goes on your Linux distribution.
It should probably be `~/.config/systemd/user`

see: https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Unit%20File%20Load%20Path


## Modify env file

In `aws-sts-refresh.env`, set all the values to match your system.

In `aws-sts-refresh.service`, set the path to the binary (probably your GOPATH). This must be an absolute path.

## Validate that it works

systemctl --user start aws-sts-refresh.service


## Start and enable timer
```
systemctl --user start aws-sts-refresh.timer
systemctl --user enable aws-sts-refresh.timer
```

## Enable timer to start at boot
```
loginctl enable-linger <YOUR-USERNAME>
```

## Troubleshooting

Tail logs with (don't forget `--user`):

```
journalctl --user -u aws-sts-refresh.service -lf
```
