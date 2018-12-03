#!/bin/python3

# This is a installer for remotecredencials
# Author: cotonetearabe
# For Debian and Arch based
#

import os


def getDistroNInstall():
	with open('/etc/os-release', 'r') as fRelease:
		distName = fRelease.readlines()[0].strip('NAME=')
		if distName == 'Arch Linux':
			os.system('sudo pacman -Sy --noconfirm python-pip')
			os.system('sudo pacman -Sy --noconfirm go go-tools')
		elif distName == '"Ubuntu"' or distName == '"Debian"' or distName == '"Linux Mint':
			os.system('sudo apt install python-pip -yy')
			os.system('sudo apt install go go-tools -yy')


def main():
	getDistroNInstall()
	os.system('sudo pip3 install -r old/requirements.txt ')


if __name__ == '__main__':
	main()
