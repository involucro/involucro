# Contributing

Involucro welcomes new development!
This document briefly describes how to contribute to the [ivolucro project](https://github.com/involucro/involucro).

## Before you Begin

If you have an idea for a feature to add or an approach for a bugfix,
it is best to communicate with us early on. The most common venues for this are
[GitHub issues](https://github.com/involucro/involucro/issues).
Browse through existing GitHub issues and if one seems related,
comment on it. We are generally available via [gitter](TODO).

## Reporting a new issue

If no existing involucro issue seems appropriate, a new issue can be
opened using [this form](https://github.com/involucro/involucro/issues/new).

## How to Contribute

* All changes to the [involucro](https://github.com/involucro/involucro)
  should be made through pull requests to this repository (with just two
  exceptions outlined below).

* If you are new to Git, the [Try Git](http://try.github.com/) tutorial is a good places to start.
  More learning resources are listed at https://help.github.com/articles/good-resources-for-learning-git-and-github/ .

* Make sure you have a free [GitHub](https://github.com/) account.

* Fork the [involucro repository](https://github.com/involucro/involucro) on
  GitHub to make your changes.
  To keep your copy up to date with respect to the main repository, you need to
  frequently [sync your fork](https://help.github.com/articles/syncing-a-fork/):
  ```
    $ git remote add upstream https://github.com/involucro/involucro
    $ git fetch upstream
    $ git checkout dev
    $ git merge upstream/dev
  ```

* Additions of new features to the code base should be pushed to the `master` branch (`git
  checkout master`).

* If your changes modify code - please ensure the resulting files
  conform to the Go [guidelines](https://gobyexample.com/).

* Commit and push your changes to your
  [fork](https://help.github.com/articles/pushing-to-a-remote/).

* Open a [pull
  request](https://help.github.com/articles/creating-a-pull-request/)
  with these changes. You pull request message ideally should include:

   * A description of why the changes should be made.

   * A description of the implementation of the changes.

   * A description of how to test the changes.

* The pull request should pass all the continuous integration tests which are
  automatically run by GitHub using e.g. Travis CI.

* After reviewing your pull request and optional discussions we will merge it into the main branch.
