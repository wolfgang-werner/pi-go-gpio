
# Simple GPIO library for Raspberry Pi in Go

This module implements a Go module to manipulate the GPIO pins.
It has been created as a first try to understand the Raspberry Pi programming and GPIO.
If you look for a decent implementation that will work in furture linux versions have a look at [gpiod](https://github.com/warthog618/gpiod). 

## Features

- only uses the sysfs mappings
- with a proper setup it can be used without being root
- tested on Raspberry Pi 3 / Raspberry Pi 4
- tested with Ubuntu Server 20.10 (Groovy Gorilla)

## GPIO

This module uses the pins to identify the GPIO ports.

This image from [https://pinout.xyz](https://pinout.xyz) proved helpful to find the right associations.

![Raspberry Pi GPIO Pinout](./doc/img/raspberry-pi-pinout.png)

## Setup

I started with a fresh operating system image. 
To enable user access to the gpio via sysfs one has to add these configurations.

### add group to allow gpio access 

Add a group to be used as owner of the gpio sys file entries

````shell
~ addgroup gpio
~ adduser ubuntu gpio
````

### add configuration to udev rules

Add a file ´´/etc/udev/rules.d/99-gpio.rules´´ with the following content, reboot your pi

````shell
~ sudo su -
~ cat /etc/udev/rules.d/99-gpio.rules
#
# gpio access for group gpio
#
SUBSYSTEM=="gpio", GROUP="gpio", MODE="0660"
SUBSYSTEM=="gpio*", PROGRAM="/bin/sh -c '\
chown -R root:gpio /sys/class/gpio && chmod -R 770 /sys/class/gpio;\
chown -R root:gpio /sys/devices/virtual/gpio && chmod -R 770 /sys/devices/virtual/gpio;\
chown -R root:gpio /sys$devpath && chmod -R 770 /sys$devpath\
'"
~ reboot now
````
