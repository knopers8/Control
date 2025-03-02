## coconut repository default-revision

set default global and per-repository revision

### Synopsis

The repository default-revision command sets the global default repository revision.

To set a per repository default revision, the default revision specified needs to be preceded by the repository index (not its name), as is reported by `coconut repo list`.

Examples:
* `coconut repo default-revision basic-tasks` Sets `basic-tasks`as the global default-revision
* `coconut repo default-revision 0 master` Sets `master`as the default-revision for repo with index 0
* `coconut repo default-revision 2 vs-sftb` Sets `vs-sftb`as the default-revision for repo with index 2

```
coconut repository default-revision [flags]
```

### Options

```
  -h, --help   help for default-revision
```

### Options inherited from parent commands

```
      --config string            optional configuration file for coconut (default $HOME/.config/coconut/settings.yaml)
      --config_endpoint string   configuration endpoint used by AliECS core as PROTO://HOST:PORT (default "consul://127.0.0.1:8500")
      --endpoint string          AliECS core endpoint as HOST:PORT (default "127.0.0.1:32102")
      --nospinner                disable animations in output
  -v, --verbose                  show verbose output for debug purposes
```

### SEE ALSO

* [coconut repository](coconut_repository.md)	 - manage git repositories for task and workflow configuration

###### Auto generated by spf13/cobra on 14-Jun-2021
