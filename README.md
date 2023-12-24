# fnoios
fnoios (_Frida_ _No_ _ios_-deploy) is a Frida based tool that redirects ios applications stdout and stderr to pseudo terminals and reads from it. In situations where you would like to see what 
the app logs but you cannot spawn the application with `ios-deploy` you can utilize this. You need to have `frida-devkit` installed previously.

# Installation

_Using go install_

```bash
$ go install github.com/NSEcho/fnoios@latest
```

_Manually_

```bash
$ git clone https://github.com/NSEcho/fnoios && cd fnoios
$ go build
$ ./fnoios --help
iOS stdout/stderr => pty

Usage:
  fnoios [flags]

Flags:
  -a, --app string   Application name to attach to
  -h, --help         help for fnoios
  -p, --pid int      PID of process to attach to (default -1)
  -s, --spawn        Spawn the app/file
```

You can attach to the application by providing its name as `-a` flag or you can spawn the application by passing `-a` with 
bundle identifier along with the `-s` flag.