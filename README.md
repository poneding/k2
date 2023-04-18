# k2  

## 简介

k2（kubernetes-tools）是围绕 kubernetes 的运维工具，目前包含的功能有：

- 合并多个 kubeconfig 文件
- 从 serviceaccount 生成 kubeconfig 文件

## 安装

```bash
go install github.com/poneding/k2@latest
```

## 命令自动补全

按照以下命令方式配置，重启新的终端生效。

```bash
# bash:
# Linux:
$ k2 completion bash > /etc/bash_completion.d/k2
# macOS:
$ k2 completion bash > $(brew --prefix)/etc/bash_completion.d/k2

# zsh:
$ echo "autoload -U compinit; compinit" >> ~/.zshrc
$ k2 completion zsh > "${fpath[1]}/_k2"

# fish:
$ k2 completion fish > ~/.config/fish/completions/k2.fish

# PowerShell:
PS> k2 completion powershell > k2.ps1
```

## 使用

### config

`k2 config` 命令包含一些围绕 `kubeconfig` 文件的运维操作。

**gen**  
使用 `gen` 子命令，从 `ServiceAccount` 生成 `kubeconfig` 文件

```bash
k2 config gen --from <serviceaccount> -n <namespace>
k2 config gen --from <serviceaccount> -n <namespace> --to <file>
```

- 使用 `--from` 或 `-f` 传入 `ServiceAccount` 名称;  
- 使用 `--namespace` 或 `-n` 指定 `ServiceAccount` 所在的命名空间，否则默认指定 `default`；  
- 使用 `--kubeconfig` 指定访问集群配置文件，默认指定 `~/.kube/config`；
- 使用 `--to` 或 `-t` 将生成的 `kubeconfig` 写入到目标文件，目标文件不存在则创建。不使用该参数会直接将生成结果打印到控制台。

**merge**  
使用子命令 `merge` 合并多个 `kubeconfig` 文件

```bash
k2 config merge --from <config_a> --to <config_b>
k2 config merge --from <config_a>,<config_b> --to <config_c>
```

- 使用 `--from` 或 `-f` 指定待合并的 `kubeconfig` 文件，多个文件以 `,` 分隔；  
- 使用 `--to` 或 `-t` 指定目标 `kubeconfig` 文件，目标文件不存在则创建。不使用该参数会直接将合并结果打印到控制台。
