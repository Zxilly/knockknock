# knock knock

Based on [Post](https://www.v2ex.com/t/858294)

## Download

[Release](https://github.com/Zxilly/knockknock/releases/tag/latest)

If your CPU is modern architecture, you can use `amd64-v3`, or you can choose `amd64-v1`.

## Usage

```bash
kk ip:port
```

For example:

```bash
PS C:\knockknock> ./kk 9.9.9.9:443
Testing 9.9.9.9:443
Seems OK.
PS C:\knockknock> ./kk 8.8.8.8:443
Testing 8.8.8.8:443
Seems blocked.
```

> ### Warning
> 
> It may not stable if you test multi times on the same IP.