# Helm Classic Sec: Work On Secrets

The `helm-sec` plugin provides a tool for working with Kubernetes
secrets.

It can:

- Handle encoding secrets for you.
- Generate or manage some kinds of secrets for you.
- Create or modify secrets files for you.

## Examples

The simplest invocation of `helmc sec` generates a secret file and sends
it to stdout:

```
$ helm-sec name value
kind: Secret
apiVersion: v1
metadata:
  name: name
data:
  name: dmFsdWU=
```

(Note that `dmFsdWU=` is `value` base64 encoded)

You can send the output to a file by specifying the file name with the
`--file` or `-f` flags:

```
$ helm-sec -f secret.yaml name value
```

And `helmc sec` can generate passwords or tokens for you:

```
$ helmc sec --password mysecret
---> Password: jb@OTr}k|dG<jc,m
kind: Secret
apiVersion: v1
metadata:
  name: mysecret
data:
  mysecret: amJAT1RyfWt8ZEc8amMsbQ==
```

Use the `--length,-l` flag to specify how long of a password or token
you'd like. You can also use `--alphanum` and `--alpha` to generate
alphanumeric or alphabetic phrases (no special characters).

`helmc sec` can also generate keypairs. To generate a NaCl Box pair, use
the `--box` flag:

```
$ helm-sec --box mysecret
kind: Secret
apiVersion: v1
metadata:
  name: mysecret
  data:
    mysecret: RdueTOMUjMjSnarkdlOR0Hq1Q/abTTNQ1xFp/Wn4dNI=
    mysecret.pub: 9olV4AbjU6QdgJcGfevT+pCLc/0NfpbD6z9OPkfYK28=
```

In the future, `helm-sec` will also generate:

- SSH key pairs
- SSL certs
- More
