# How to Contribute

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Helm Classic (the `helm/helm` repository) is **no longer actively developed** but will remain available until `kubernetes/helm` has stabilized.
|---|---|

Helm and [Deployment Manager](https://github.com/kubernetes/deployment-manager)
have recently joined forces to make deploying and managing software on
Kubernetes as simple as possible. The combined effort now lives in the Kubernetes GitHub organization at
[kubernetes/helm](https://github.com/kubernetes/helm).

We love getting Pull Requests (PRs) _for bug fixes only_. To make sure we keep code quality
and consistency high, we do have some process in place. In a nutshell:

- Code should follow Go coding standards and pass `go lint` and `go
  vet`. (These tools are automatically run on every pull request.)
- PRs should contain a single ~~feature or~~ fix, and follow the conventions
  linked below.
- Contributors must agree to the DCO
- Every patch must be signed off by two core contributors (and this is
  made easier when the community at large weighs in on PRs, too).

Helm Classic follows the contribution guidelines established by the Deis
project. Please take a look at them.

- [Deis Contribution Guidelines](https://github.com/deis/deis/blob/master/CONTRIBUTING.md)
- [Details on commit messages](http://docs.deis.io/en/latest/contributing/standards/#commit-style-guide)

Again, thanks for taking the time to contribute. We know the process
isn't trivial, and we really appreciate your helping ensure that Helm Classic
develops into a high quality tool.

## Interpreting the Labels in GitHub

We use labels on GitHub to indicate the state of something. Here are a
few of the more interesting labels.

- Awaiting review: The PR is ready to be considered for merging. If a PR
  does not have this, we assume that it's a work in progress (even if
  the issue does not say WIP)
- Proposal: something we're discussing, but haven't decided on
- Enhancement: a feature or chore
- Bug: a bug
- Needs manual testing: one of the core team must manually test before
  LGTM
- LGTM1: First LGTM (Looks Good To Me)
- LGTM2: Second LGTM; this means a core contributor can merge it

