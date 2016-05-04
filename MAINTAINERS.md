# Helm Classic Maintainers

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/156px-Warning.svg.png) | Helm Classic (the `helm/helm` repository) is **no longer actively developed** but will remain available until `kubernetes/helm` has stabilized.
|---|---|

Helm and [Deployment Manager](https://github.com/kubernetes/deployment-manager)
have recently joined forces to make deploying and managing software on
Kubernetes as simple as possible. The combined effort now lives in the Kubernetes GitHub organization at
[kubernetes/helm](https://github.com/kubernetes/helm).

This document serves to describe the leadership structure of the Helm Classic project, and to list the current
project maintainers. _Maintainers of Helm Classic currently perform bug fixes and critical maintenance only._

# What is a maintainer?

(Unabashedly stolen from the [Docker](https://github.com/docker/docker/blob/master/MAINTAINERS) project)

There are different types of maintainers, with different responsibilities, but
all maintainers have 3 things in common:

1. They share responsibility in the project's success.
2. They have made a long-term, recurring time investment to improve the project.
3. They spend that time doing whatever needs to be done, not necessarily what
is the most interesting or fun.

Maintainers are often under-appreciated, because their work is harder to appreciate.
It's easy to appreciate a really cool and technically advanced feature. It's harder
to appreciate the absence of bugs, the slow but steady improvement in stability,
or the reliability of a release process. But those things distinguish a good
project from a great one.

# Helm Classic maintainers

Helm Classic has two groups of maintainers: core and contributing.

## Core maintainers

Core maintainers are knowledgeable about all areas of Helm Classic. Some maintainers work on Helm Classic (and Helm)
full-time, although this is not a requirement.

The duties of a core maintainer include:
* Classify and respond to GitHub issues and review pull requests
* ~~Help to shape the Helm Classic roadmap and lead efforts to accomplish roadmap milestones~~
* Participate actively in ~~feature development~~ and bug fixing
* Answer questions and help users in Slack

The current core maintainers of Helm Classic:
* Matt Butcher - <mbutcher@deis.com> ([@technosophos](https://github.com/technosophos))
* Gabe Monroy - <gabe@deis.com> ([@gabrtv](https://github.com/gabrtv))
* Kent Rancourt - <kent@deis.com> ([@krancour](https://github.com/krancour))
* Keerthan Reddy Mala - <kmala@deis.com> ([@kmala](https://github.com/kmala))


### Pull requests

No pull requests can be merged until at least one core maintainer signs off with an "LGTM" label.
The other LGTM can come from either a core maintainer or contributing maintainer. A maintainer who
creates a pull request should also be the one to merge it, after two LGTMs.

## Contributing maintainers

Contributing maintainers may have deep knowledge about some but not necessarily all areas of Helm Classic.
Core maintainers will enlist contributing maintainers to weigh in on issues, review pull
requests, or join conversations as needed in their areas of expertise.

The duties of a contributing maintainer are similar to those of a core maintainer, but may be
scoped to relevant areas of the Helm Classic project.

Contributing maintainers are defined in practice as those who have write access to the Helm Classic
repository. All maintainers can review pull requests and add LGTM labels as appropriate.

## Becoming a maintainer

The Helm Classic project will succeed exactly as its community thrives. It is essential that the breadth
of potential Kubernetes users find Helm Classic useful enough to help it grow. If you use Helm Classic every day,
we want you to help determine where the ship is steered.

Generally, potential contributing maintainers are selected by the Helm Classic core maintainers based in
part on the following criteria:
* Sustained contributions to the project over a period of time
* A willingness to help Helm Classic users on GitHub and in Slack
* A friendly attitude!

The Helm Classic core maintainers must agree in unison before inviting a community member to join as a
contributing maintainer.
