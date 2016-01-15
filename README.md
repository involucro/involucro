
# Involucro - Build and Deliver Software with Containers

[![Build Status](https://travis-ci.com/thriqon/involucro.svg?token=43LaxsfrMzmC1LqNeQwa)](https://travis-ci.com/thriqon/involucro)

## Introduction

Building and delivering software is a complex task, that can be made easier
with proper encapsulation: Containers. They are able to pack all required
dependencies (bitwise version controlled), so that you can release code with the
confidence that everything works the same way it worked on the developers machine.

The current default for building and packing software into containers is using
*Dockerfiles*<sup>[1](#trademarks)</sup> offered by the
[Docker®](https://www.docker.com) software. This, however, has several
structural disadvantages:

* Build Container is the Deliverable Container
* Only one-dimensional caching of build steps
* No support for slimmed containers
* Unnecessary layers

Involucro [from Latin: envelope] detaches the build process from the
deliverable container and re-establishes proper encapsulation in containers:
One Process, One Container.

## Installation

From source:

    $ go get github.com/thriqon/involucro
    $ $GOPATH/bin/involucro -v

As binary:

    $ wget https://storage.googleapis.com/involucro-1149.appspot.com/involucro

Or, for windows:

    https://storage.googleapis.com/involucro-1149.appspot.com/involucro.exe
    https://storage.googleapis.com/involucro-1149.appspot.com/involucro32.exe

## Usage

Involucro is configured by a Lua script file. By default, it is looking for
`invfile.lua` in the current directory, but this can be overridden (see below).

A configuration file contains a set of tasks, identified by a unique name.
These names can be specified when invoking `involucro`, and are executed in the
order they are given.  For example, `$ involucro build package` will run the
`build` and afterwards the `package` task.  A task can be created by invoking
`inv.task('<ID>')` in the configuration file.

For easy readability, the configuration file uses a fluent syntax to build the
tasks. The available methods are either modifying the next registered step, or
are registering a step. This type distinction is documented below for each
method.

**inv.task**`('<ID>')` (*modifier*) sets the task of the next registered step
to `<ID>`. It makes the methods `using`, `runTask`, and `wrap` available.

Each task consists of a list of steps that are run in the order they are given
in the file. There are different types of steps. Each step has one
*introductory* method made available from the task, a set of *modifying*
methods setting different properties of the step, and a final *registration*
method that registers these settings for execution. The current status can be
stored at any point into a variable, and reused later. However, the steps are
strictly run in the order their registration method was called.

### Run Step

A run step executes a Docker container. By default, the current working
directory is mounted as `/source` into the container, which is also configured
to be the working directory of the process running in the container. It is
mainly used to transform source code using external processes such as compilers
into a different form. 

**task.using**`('<IMAGE_ID>')` (*introductory*) starts off a run step by
specifying the repository name (optionally with tag) or the image ID of the
image to be run. Example: `task.using('gcc:4.9')`.

**runstep.withConfig**`(<TABLE>)` (*modifying*) sets the values in the Lua
table as configuration values for the Docker container. The values that can be
set here are only affecting the container itself, not how is connected with the
host. See *withHostConfig* for this. The options available are
[Config](https://godoc.org/github.com/fsouza/go-dockerclient#Config). The keys
are interpreted case insensitive. Example: `runstep.withConfig({Cmd =
{"/bin/echo", "Hello, World!"}})`.

**runstep.withHostConfig**`(<TABLE>)` (*modifying*) sets the values in the Lua
table as host configuration values. These values control the exact execution
semantics of the container from the hosts point of view. The available options
are documented here:
[HostConfig](https://godoc.org/github.com/fsouza/go-dockerclient#HostConfig).
Example: `runStep.withConfig({links = {"redis"}})`.

NOTE: By default, `involucro` binds the current directory as `/source`. If the
`Binds` key is set in the given table, it overwrites this binding. `Involucro`
however interprets the given bindings, and changes all relative source bindings
to absolute paths. This enables bindings such as `{binds = {"./dist:/data",
"/tmp:/tmp"}}`.

**runstep.withExpectation**`(<TABLE)` (*modifying*) registers expectations
towards the output and exit code of the process. By default, `involucro`
expects the process to exit cleanly with exit code `0`. Tests of executables
however may require expecting a process to fail. This can be set with the key
`code`: `runstep.withExpectation({code = 1})`.  Similarly, an expectation
towards the output of the process on `stdout` and/or `stderr` can be registered
with regular expressions conforming to [Re2
syntax](https://github.com/google/re2/wiki/Syntax).  Example:
`runstep.withExpectation({stdout = "Hello, World!\n"})`.

**runstep.run**(`'<CMD>'...`) (*registration*) registers the run step. The
arguments are used as the command-line arguments of the process being run. It
directly follows Docker semantics regarding process execution. Each argument is
used as a single argument for the process. Example: `runstep.run('/bin/echo',
'Hello, World!')`. Note that there is no wildcard expansion or variable
replacement if the arguments are not given to a shell, such as `/bin/sh`.
Example: `runstep.run('/bin/sh', '-c', 'echo *')`.

### Wrap Step

A wrap step takes the contents of a directory and creates an image layer out of
it, optionally with a parent image layer and meta data. The resulting image can
be tagged into a repository with a tag name (or `latest`, if none is set).

**task.wrap**`('<SOURCE_DIR>')` (*introductory*) starts off a wrap step by
specifying the directory containing the files that are to be wrapped into an
image. It is also possible to use the current directory (`.`).  Example:
`task.wrap('dist')`.

**wrapstep.at**`('<TARGET_DIR>')` (*modifying*) sets the directory in the resulting
image into which the files are copied.  This can be used to put HTML files into
the location the web server in the parent image expects them to be. This
directory doesn't need to exist yet. Example: `wrapstep.at('/data')`.

**wrapstep.inImage**`('<PARENT_IMAGE>')` (*modifying*) causes the resulting image
to be a child of the image identified by the parameter. If this modification is
omitted, the resulting image is parent-less. Example:
`wrapstep.inImage('nginx')`.

**wrapstep.withConfig**`(<TABLE>)` (*modifying*) sets configuration values similar
to the *withConfig* method of the run step above. This can be used to pre-set
an entrypoint or exposed ports. Example: `wrapstep.withConfig({exposedports =
{"80/tcp"}})`.

**wrapstep.as**`('<IMAGE_NAME>')` (*registration*) registers the step for
execution. The image constructed by the previous modifications is built and
tagged with the given name, which may include a registry designation. Example:
`wrapstep.as('app:latest')`

### Runtask Step

As a convenience, it is possible to run another task as part of a task. This emulates the conventional `all` task
from `Makefile`s. Exceptionally, the introductory method for this step is also the registration method.

**task.runTask**(`<ID>`) (*introductory registration*) registers a step that executes the task with the given ID
as part of the steps in this task. Example: `inv.task('all').runTask('compile').runTask('package')`.

### Tag Step

Sometimes, there should be two versions of the same image sharing the same image ID, for example to have the `latest`
tag equivalent to version `v2`. The tag step helps in this case.

**task.tag**(`<NAME>`) (*introductory*) starts a tagging by setting the name of the original image. This can be anything
Docker accepts, including `test/asd:v2`, but also actual image IDs. Example: `task.tag('test/asd')`.

**tagstep.as**(`<NAME>`) (*registration*) registers a step that tags the image named in introductory method to the name given
as parameter. Example: `tagstep.as('test/asd')`.

## Trademarks

Docker® is a registered trademark of Docker, Inc.


