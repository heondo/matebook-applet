# matebook-applet
System tray applet for Huawei Matebook

This simple applet is designed to make some of the proptietary Huawei PC Manager's functionality available via GUI on Linux. It currently only works for Huawei MateBook 13 and depends on [batpro script](https://github.com/nekr0z/linux-on-huawei-matebook-13-2019/blob/master/batpro) and [fnlock script](https://github.com/nekr0z/linux-on-huawei-matebook-13-2019/blob/master/fnlock).

## Installation and setup
The applet requires no installation as such (although you may want to add it to autorun so that it gets started automatically when you load your DE). However, it relies on at least one of the scripts to be properly installed:

* To get battery protection functionality:

    1. Download the [batpro script](https://github.com/nekr0z/linux-on-huawei-matebook-13-2019/blob/master/batpro), make it executable and copy to root executables path, i.e.:
        ```
        # mv batpro /usr/sbin/
        ```
    2. Allow your user to execute the script without providing sudo credentials. One way to do this would be to explicitly add sudoers permission (please be aware of possible security implications):
        ```
        # echo "your_username ALL = NOPASSWD: /usr/sbin/batpro" >> /etc/sudoers.d/batpro
        # chmod 0440 /etc/sudoers.d/batpro
        ```
  3. Double check that you can successfully run the script without providing additional authentication (i.e. password):
        ```
        $ sudo batpro status
        ```
* To get Fn-Lock functionality:

    1. Download the [fnlock script](https://github.com/nekr0z/linux-on-huawei-matebook-13-2019/blob/master/fnlock), make it executable and copy to root executables path, i.e.:
        ```
        # mv fnlock /usr/sbin/
        ```
    2. Allow your user to execute the script without providing sudo credentials. One way to do this would be to explicitly add sudoers permission (please be aware of possible security implications):

            # echo "your_username ALL = NOPASSWD: /usr/sbin/batpro" >> /etc/sudoers.d/fnlock
            # chmod 0440 /etc/sudoers.d/fnlock

    3. Double check that you can successfully run the script without providing additional authentication (i.e. password):

            $ sudo fnlock status

Now you may run the matebook-applet. Either download precompiled amd64 binary from [releases page](https://github.com/nekr0z/matebook-applet/releases) or compile it yourself:

        $ git clone https://github.com/nekr0z/matebook-applet.git
        $ cd matebook-applet
        $ go generate
        $ go build

## Usage
The user interface is intentionally as simple as they get. You get an icon in system tray that you can click and get a menu. The menu consists of current status, options to change it, and an option to quit the applet. Please be aware that the applet does not probe for current status on its own (this is intentional), so if you change your battery protection settings by other means it will not reflect the change. Clicking on the status line (top of the menu) updates it.

## Further development
I have plans to make it possible to utilize [a more general script](https://github.com/aymanbagabas/huawei_ec) to support MateBooks other than MateBook 13.

PRs are most welcome!

## Credits
This software includes the following sowtware or parts thereof:
* [getlantern/systray](https://github.com/getlantern/systray)
* [slurcooL/vfsgen](https://github.com/shurcooL/vfsgen)
* [The Go Programming Language](https://golang.org) Copyright © 2009 The Go Authors
* [Icons Mind icons](https://iconsmind.com)