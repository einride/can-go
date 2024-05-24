# Contributing

Thanks for showing interest in contributing to [CAN-go]!

## Suggesting Features

Before creating a PR, please consider that [CAN-go] is a tool that Einride uses
internally and as such, features that are not aligned with how we are using CAN
will probably not be accepted. In that case feel free to fork our project and
make the changes you want there instead.

## Reporting Issues

If you find a bug, please create an
[issue](https://github.com/einride/can-go/issues) for it, or create a PR
following the [Pull Request Guidelines](#pull-request-guidelines).

## Development

To start developing on [CAN-go] it is enough to clone the repo and run `make`.

## Pull Request Guidelines

[CAN-go] is using the
[Conventional Commits](https://www.conventionalcommits.org/) commit message
convention.

Keep your commits as small as possible, but still keep all changes related to a
logical change in the same commit. When receiving review feedback, fix up the
commits with the changes addressing the feedback and force-push, please don't
send fix commits.

Before opening a PR, please make sure you get no errors when running `make` and
that there is sufficient test coverage for added or changed functionality.

[can-go]: https://go.einride.tech/can
