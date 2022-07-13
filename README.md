<!-- template:define:options
{
  "nodescription": true
}
-->
<img title="Logo" src="./examples/_images/logo.png" width="961">

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/bubblezone/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/bubblezone?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/bubblezone/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/bubblezone?style=flat-square">
  </a>


  <a href="https://github.com/lrstanley/bubblezone/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/workflow/status/lrstanley/bubblezone/test/master?label=test&style=flat-square&event=push">
  </a>

  <a href="https://codecov.io/gh/lrstanley/bubblezone">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/bubblezone/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/bubblezone">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/bubblezone?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/bubblezone">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/bubblezone?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/bubblezone/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/bubblezone/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/bubblezone/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/bubblezone/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/bubblezone/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/bubblezone?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/bubblezone/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [Problem](#x-problem)
  - [Solution](#heavy_check_mark-solution)
  - [Features](#sparkles-features)
  - [Usage](#gear-usage)
  - [Examples](#clap-examples)
    - [List example](#list-example)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :x: Problem

[BubbleTea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss)
allow you to build extremely fast terminal interfaces, in a semantic and scalable
way. Through abstracting layout, colors, events, and more, it's very easy to build
a user-friendly application. BubbleTea also supports mouse events, either through
the "basic" mouse events, like `MouseLeft`, `MouseRight`, `MouseWheelUp` and
`MouseWheelDown` ([and more](https://github.com/charmbracelet/bubbletea/blob/0a0182e55a30e85640a53b8e01dc9ef06824cce5/mouse.go#L38-L48)),
or through full motion tracking, allowing hover and mouse movement tracking.

This works great for a single-component application, where state is managed in one
location. However, when you start expanding your application, where components have
various children, and those children have children, calculating mouse events like
`MouseLeft` and `MouseRight` and determining which component was actually clicked
becomes complicated, and rather tedious.

## :heavy_check_mark: Solution

**BubbleZone** is one solution to this problem. BubbleZone allows you to wrap your
components in **zero-printable-width** (to not impact `lipgloss.Width()` calculations)
identifiers. Additionally, there is a scan method that wraps the entire application,
that stores the offsets of those identifiers as `zones`, then removing them from
the resulting output.

Any time there is a mouse event, pass it down to all children, thus allowing you
to easily check if the event is within the bounds of the components `zone`. This
makes it very simple to do things like focusing various components, clicking
"buttons", and more. Take a look at this example, where I didn't have to calculate
where the mouse was being clicked, and which component was under the mouse:

![bubblezone example](https://ls-screen.s3.us-west-004.backblazeb2.com/2022/07/WindowsTerminal_XxiuWQ2hVL.gif)

## :sparkles: Features

- :heavy_check_mark: It's **_fast_** -- given it has to process this information for every render, I
  tried to focus on performance where possible. If you see where improvements can
  be made, let me know!
- :heavy_check_mark: It doesn't impact width calculations when using `lipgloss.Width()` (if you're
  using `len()` it will).
- :heavy_check_mark: It's simple -- easily determine offset or if an event was within the bounds of
  a zone.
- :heavy_check_mark: Want the mouse event position relative to the component? Easy!
- :heavy_check_mark: Provides an _optional_ global manager, when you have full access to all components,
  so you don't have to inject it as a dependency to all components.

---

## :gear: Usage

<!-- template:begin:goget -->
<!-- do not edit anything in this "template" block, its auto-generated -->
```console
$ go get -u github.com/lrstanley/bubblezone@latest
```
<!-- template:end:goget -->

TODO

## :clap: Examples

### List example

- All titles are marked as a unique zone, and upon left click, that item is focused.
- [Example source](./examples/list-default/main.go).

![list-default example](https://ls-screen.s3.us-west-004.backblazeb2.com/2022/07/WindowsTerminal_SelC1Vzdas.gif)

---

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

   * :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
   * :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
   * :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/bubblezone/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

   * :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
	   community.
   * :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
   * :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/bubblezone/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2022 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
