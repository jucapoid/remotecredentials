#!/bin/python3
from PIL import Image, ImageFont, ImageDraw, PSDraw
from tkinter import *
import pyqrcode, sys
import PIL.Image as PImage


def acesso(acessoA):
		draw.rectangle(((200, 170), (235, 205)), fill="red", outline="white")
		draw.text((200 + 13, 170 + 6), acessoA[0], (0, 0, 0), font=font)
		draw.rectangle(((235, 170), (270, 205)), fill="red", outline="white")
		draw.text((235 + 13, 170 + 6), acessoA[1], (0, 0, 0), font=font)
		draw.rectangle(((270, 170), (305, 205)), fill="red", outline="white")
		draw.text((270 + 13, 170 + 6), acessoA[2], (0, 0, 0), font=font)
		draw.rectangle(((305, 170), (340, 205)), fill="red", outline="white")
		draw.text((305 + 13, 170 + 6), acessoA[3], (0, 0, 0), font=font)
		draw.rectangle(((200 + 35 / 2, 205), (235 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((200 + 35 / 2 + 13, 205 + 6), acessoA[4], (0, 0, 0), font=font)
		draw.rectangle(((235 + 35 / 2, 205), (270 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((235 + 35 / 2 + 13, 205 + 6), acessoA[5], (0, 0, 0), font=font)
		draw.rectangle(((270 + 35 / 2, 205), (305 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((270 + 35 / 2 + 13, 205 + 6), acessoA[6], (0, 0, 0), font=font)


photo = sys.argv[1]
name = sys.argv[2]
cc = sys.argv[3]
acessoA = sys.argv[4:]  # X or number
if len(sys.argv[4:]) > 8:
	print("Too much args")
	sys.exit(1)
credN = sys.argv[1] + sys.argv[2]
fundo = "cred.png"
"""
photo = "photo.jpg"
name = "zeTeste"
cc = "11227788"
credN = photo + name
"""
# credN = sys.argv[1:3].join()


img = PImage.open(fundo)
draw = ImageDraw.Draw(img)

width, height = img.size

# font = ImageFont.load("arial.pil")
# font = ImageFont.truetype("ariblk.ttf", 15)
font = ImageFont.load_default()

w, h = font.getsize("1")
# acesso(1, 2, 3, 4, 5, 6, 7)

# foto
draw.rectangle(((30, 90), (140, 220)), fill="white", outline="white")

# foto = Image.open(sys.argv[1])  # <---- input FOTO
foto = PImage.open(photo)  # <---- input FOTO
foto = foto.resize((115, 135))

# img.paste(foto, (30, 90, 145, 225))

draw.text((30, 250), "NOME:", (255, 255, 255), font=font)
draw.rectangle(((30, 275), (255, 300)), fill="white", outline="white")
# draw.text((35, 276), sys.argv[2], (0, 0, 0), font=font)  # input <---- NAME
draw.text((35, 276), name, (0, 0, 0), font=font)  # input <---- NAME


draw.text((30, 320), "B.I/C.C/N.M :", (255, 255, 255), font=font)
draw.rectangle(((30, 345), (255, 370)), fill="white", outline="white")
# draw.text((35, 346), sys.argv[3], (0, 0, 0), font=font)  # input <------ CC
draw.text((35, 346), cc, (0, 0, 0), font=font)  # input <------ CC


draw.text((30, height - 30), "QUEIMA.AAUE.PT", (255, 255, 255), font=font)


qr = pyqrcode.create('A232T3QF2017')
qr.png('qrcode.png', scale=3)

qr = PImage.open('qrcode.png')
x, y = qr.size

img.paste(qr, (0 + 49, 0 + 30 + 370, x + 49, y + 30 + 370))


draw.rectangle(((33, height - 55), (35 + 130, height - 30)), fill="white", outline="white")
draw.text((37, height - 54), 'A232T3QF2017', (0, 0, 0), font=font)  # is this a constant?

acesso(acessoA)

img.save(credN + ".png")

