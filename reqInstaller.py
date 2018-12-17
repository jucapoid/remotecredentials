#!/bin/python3

# This is a requirement installer for remotecredencials
# Author: cotonetearabe
# For Debian and Arch based

import subprocess


def getDistroNInstall():
	with open('/etc/os-release', 'r') as fRelease:
		distName = fRelease.readlines()[0].strip('NAME=')
		if distName == 'Arch Linux' or distName == 'Arch' or distName == 'ArchLinux':
			print(subprocess.getoutput('sudo pacman -Sy --noconfirm python-pip'))
			print(subprocess.getoutput('sudo pacman -Sy --noconfirm go go-tools'))
		elif distName == '"Ubuntu"' or distName == '"Debian"' or distName == '"Linux Mint':
			print(subprocess.getoutput('sudo apt install python3-pip -yy'))
			print(subprocess.getoutput('sudo apt install golang -yy'))


def main():
	getDistroNInstall()
	print(subprocess.getoutput('sudo pip3 install -r old/requirements.txt'))


if __name__ == '__main__':
	main()
