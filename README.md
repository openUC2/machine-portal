# machine-portal
A web browser landing page for easy user access to deployed apps

## Introduction

Users connect to [openUC2 OS](https://github.com/openUC2/os-RPi) (which runs on a Raspberry Pi
computer) from a web browser on a client device by opening a URL like <http://home.uc2>.
This tool provides a web page with a list of links to the various network services running on the OS
(e.g. the ImSwitch GUI, the Cockpit system administration panel, or the system file manager) for
easy access. It is meant to be served from a reverse-proxy on port 80 along with all other network
services, configured as in [openUC2/pallet](https://github.com/openUC2/pallet).

## Usage

### Local Deployment

First, you will need to download machine-portal, which is available as a single self-contained
executable file. You should visit this repository's
[releases page](https://github.com/openUC2/machine-portal/releases/latest) and download an archive
file for your platform and CPU architecture; for example, on a Raspberry Pi 5, you should download
the archive named `machine-portal_{version number}_linux_arm.tar.gz` (where the version number should
be substituted). You can extract the machine-portal binary from the archive using a command like:
```bash
tar -xzf machine-portal_{version number}_{os}_{cpu architecture}.tar.gz machine-portal
```

Then you may need to move the machine-portal binary into a directory in your system path, or you can just run the machine-portal binary in your current directory (in which case you should replace `machine-portal` with `./machine-portal` in the commands listed below).

Once you have machine-portal, you can run it as follows on a Raspberry Pi:
```bash
./machine-portal
```

Then you can view the landing page at <http://localhost:3000> . Note that if you are running it on a
computer other than the Raspberry Pi with openUC2 OS, then you will need to set some environment
variables (see below) to non-default values.

### Development

To install various backend development tools, run `make install`. You will need to have installed Go first.

Before you start the server for the first time, you'll need to generate the webapp build artifacts by running `make buildweb` (which requires you to have first installed [Node.js](https://nodejs.org/en/) and [Corepack](https://github.com/nodejs/corepack)). Then you can start the server by running `make run` with the appropriate environment variables (see below); or you can run `make runlive` so that your edits to template files will be reflected after you refresh the corresponding pages in your web browser. You will need to have installed golang first. Any time you modify the webapp files (in the web/app directory), you'll need to run `make buildweb` again to rebuild the bundled CSS and JS.

### Building

Because the build pipeline builds Docker images, you will need to either have Docker Desktop or (on Ubuntu) to have installed QEMU (either with qemu-user-static from apt or by running [tonistiigi/binfmt](https://hub.docker.com/r/tonistiigi/binfmt)). You will need a version of Docker with buildx support.

To execute the full build pipeline, run `make`; to build the docker images, run `make build` (make sure you've already run `make install`). Note that `make build` will also automatically regenerate the webapp build artifacts, which means you also need to have first installed Node.js as described in the "Development" section. The resulting built binaries can be found in directories within the dist directory corresponding to OS and CPU architecture (e.g. `./dist/machine-portal_window_amd64/machine-portal.exe` or `./dist/machine-portal_linux_amd64/machine-portal`)

### Environment Variables

#### Machine Name

If you are running machine-portal on a computer which is not a Raspberry Pi with the standard openUC2 OS, then you'll need to set some environment variables. Specifically, you'll need to set:

- Either `MACHINENAME_NAME`, which should be a string representing the name of the machine to be displayed on the landing page, or `MACHINENAME_NAMEFILE`, which should be the path to a file containing the name of the machine to be displayed on the landing page.

For example, you could run machine-portal with the machine name `metal-slope-23501` with one of the following commands:
```bash
# If you downloaded a machine-portal binary:
MACHINENAME_NAME=metal-slope-23501 ./machine-portal
# If you are developing the project:
MACHINENAME_NAME=metal-slope-23501 make run
```

#### Custom Templates

You can override the default webpage templates embedded in the machine-portal binary by providing a path to the templates directory with the `TEMPLATES_PATH` variable, relative to the current working directory in which you start the machine-portal program. For example, you could provide a more-minimal "hello world" landing page by creating a new file named `index.page.tmpl` with following contents in a new `custom-templates/home` subdirectory in the directory from which you will launch machine-portal:
```html
{{template "shared/base.layout.tmpl" .}}

{{define "title" -}}
  {{- $machineName := .Data.MachineName -}}
  Machine {{$machineName}}
{{- end}}
{{define "description"}}Machine portal{{end}}

{{define "content"}}
  <main>
    <section class="section content">
      <div class="container">
        <h1>Hello, world!</h1>
        <p>
          Greetings from a custom template!
        </p>
    </section>
  </main>
{{end}}
```

and then running the following command:
```bash
# If you downloaded a machine-portal binary:
TEMPLATES_PATH=custom-templates MACHINENAME_NAME=template-test ./machine-portal
# If you are developing the project:
TEMPLATES_PATH=custom-templates MACHINENAME_NAME=template-test make run
```

Note that running `make runlive` will cause your `TEMPLATES_PATH` environment variable to be ignored, so that the templates directory at `web/templates` (relative to the root of this repository) is always used.

## Licensing

Except where otherwise indicated, source code provided here is covered by the following information:

Copyright Ethan Li and openUC2 project contributors

SPDX-License-Identifier: `Apache-2.0 OR BlueOak-1.0.0`

You can use the source code provided here either under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0) or under the [Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0); you get to decide. We are making the software available under the Apache license because it's [OSI-approved](https://writing.kemitchell.com/2019/05/05/Rely-on-OSI.html), but we like the Blue Oak Model License more because it's easier to read and understand.

### Origins

The [github.com/openUC2/machine-portal](https://github.com/openUC2/machine-portal) repository was
initialized by Ethan Li as a hard fork of
[github.com/PlanktoScope/device-portal](https://github.com/PlanktoScope/device-portal), from a
commit of the PlanktoScope/device-portal repository which only included contributions made by Ethan
Li (but not any other contributors to the PlanktoScope project).
