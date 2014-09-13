# Castbox

Castbox is a cross-platform, cross-arch, full-featured Chromecast
simulator and development environment.

You can build Chromecast apps and have a majorly faster turn-around time
for developments since you don't depend on hardware. 

Castbox should make a great development workflow via flexible configuration and detailed internal logging.

It was built to [make developers happy](http://makedevelopershappy.com).

There is also an introductory [blog post](http://blog.paracode.com/2014/09/11/open-sourcing-castbox/).

**NOTE**: Around July, Google changed its Chromecast protocols to include
private-key encryption. Due to that change, Castbox only works on the
legacy protocol and will continue to work with DIAL based apps (Youtube,
Google music, and so on).

## Usage

Make sure you have Go and [Gom](https://github.com/mattn/gom) installed.

Configure your Castfile (see below).

```
$ gom build
$ ./castbox
```

Castbox will pick up your default, local Castfile.


## Configuration

You probably can use the default Castfile, but if you're on OSX, you want to force set your current LAN IP before starting:

```
      "force_host":"10.0.0.11",
```

* `uuid`, and `name` are free-form fields for your usage. You can have
  as many castboxes as you like in your LAN as long as they have unique
IDs.

* `remote_chrome` if you are running Castbox "headless" on a
  Raspberry Pi you can
  specify where the Chrome API lives on the network manually.

* `force_chromebin` - if your Chrome binary is located somehere
  non-standard, you can set it manually.

* `idle_time_min` - apps can idle for this time, in minutes.

* `applications` - take one of the examples, but make sure to point
  your own URL. You can override existing "production" apps (e.g. Youtube) if you
specify its `app_name` exactly.



## Inject your own apps with a Castfile

The `Castfile` can _override_ and/or add over any of the currently fetched
apps. This means you can:

* Override existing apps, for example - provide your own interface to
the existing Youtube app.
* Add your own homebrew or under development apps - by giving it a
unique name or ID.
* Provide your own configuration, which will override the fetched
configuration. This is useful for setting up the default idle app; which
is done out of the box by us.

The Castfile will be picked up automatically at the current working
directory or through the `--castfile` flag.



## Chrome

You can safely ignore this section unless you like to tinker
around.


We support 2 ways to command and control a chrome instance in order to direct apps.


#### 1. We'll spawn it for you

Without specifying any flag, we'll spawn and control a special instance
of Chrome on the same machine.


#### 2. Or your own instance

The benefit of this is to be able to control chrome that sits on
another, more capable machine in terms of UI handling (perhaps you are running this on a budget,
no-ui machine).

You'll have to launch chrome in a special "remote debugging" mode where
it can take commands remotely.

For now, see https://src.chromium.org/svn/trunk/src/chrome/test/chromedriver/chrome_launcher.cc




# Contributing

Fork, implement, add tests, pull request, get my everlasting thanks and a respectable place here :).

# Copyright

Copyright (c) 2014 [Dotan Nahum](http://gplus.to/dotan) [@jondot](http://twitter.com/jondot). See MIT-LICENSE for further details.




