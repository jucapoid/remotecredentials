"""
Dividir em funçoes
"""

from PIL import Image, ImageFont, ImageDraw, PSDraw
from tkinter import *
import pyqrcode

img = Image.open('cred.png')
draw = ImageDraw.Draw(img)

width, height = img.size

font = ImageFont.truetype("ariblk.ttf", 15)


w, h = font.getsize("1")
acesso(1, 2, 3, 4, 5, 6, 7)

# foto
draw.rectangle(((30, 90), (140, 220)), fill="white", outline="white")

foto = Image.open('foto.jpg')  # <---- input FOTO
foto = foto.resize((115, 135))
img.paste(foto, (30, 90, 145, 225))

draw.text((30, 250), "NOME:", (255, 255, 255), font=font)
draw.rectangle(((30, 275), (255, 300)), fill="white", outline="white")
draw.text((35, 276), "Major Faggot", (0, 0, 0), font=font)  # input <---- NAME


draw.text((30, 320), "B.I/C.C/N.M :", (255, 255, 255), font=font)
draw.rectangle(((30, 345), (255, 370)), fill="white", outline="white")
draw.text((35, 346), "14021467", (0, 0, 0), font=font)  # input <------ CC


draw.text((30, height - 30), "QUEIMA.AAUE.PT", (255, 255, 255), font=font)


qr = pyqrcode.create('A232T3QF2017')
qr.png('qrcode.png', scale=3)

qr = Image.open('qrcode.png')
x, y = qr.size

img.paste(qr, (0 + 49, 0 + 30 + 370, x + 49, y + 30 + 370))


draw.rectangle(((33, height - 55), (35 + 130, height - 30)), fill="white", outline="white")
draw.text((37, height - 54), 'A232T3QF2017', (0, 0, 0), font=font)  # nao se esta parte e fixa ou nao

img.save("cred_final.png")  # por o nome sendo a conjunção das strings do nome ou algo de identefique


def acesso(a1, a2, a3, a4, a5, a6, a7):
		draw.rectangle(((200, 170), (235, 205)), fill="red", outline="white")  # testar com test modificado com estes acessos
		draw.text((200 + 13, 170 + 6), a1, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((235, 170), (270, 205)), fill="red", outline="white")
		draw.text((235 + 13, 170 + 6), a2, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((270, 170), (305, 205)), fill="red", outline="white")
		draw.text((270 + 13, 170 + 6), a3, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((305, 170), (340, 205)), fill="red", outline="white")
		draw.text((305 + 13, 170 + 6), a4, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((200 + 35 / 2, 205), (235 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((200 + 35 / 2 + 13, 205 + 6), a5, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((235 + 35 / 2, 205), (270 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((235 + 35 / 2 + 13, 205 + 6), a6, (0, 0, 0), font=font, fill="black")
		draw.rectangle(((270 + 35 / 2, 205), (305 + 35 / 2, 240)), fill="red", outline="white")
		draw.text((270 + 35 / 2 + 13, 205 + 6), a7, (0, 0, 0), font=font, fill="black")
