from PyPDF2 import PdfFileWriter, PdfFileReader
from tkinter import *
from tkinter import font,ttk,messagebox,filedialog
from reportlab.lib.pagesizes import A4,A6
from reportlab.pdfgen import canvas
import random
from PIL import Image,ImageDraw,ImageFont
import db
import os.path
import pyqrcode

class Application(Frame):
    def __init__(self, master=None):
        Frame.__init__(self, master)
        self.pack()
        self.createWidgets()

    def createWidgets(self):
        self.grupos={1:0,2:0,3:0,4:0,5:0,6:0,7:0,8:0}      #dicionario para controlar os grupos

        self.font1=font.Font(family='DejaVu Sans',size=22)      #Letra utilizada

        self.buttonstyle = ttk.Style()                                                                              #Estilo dos butões
        self.buttonstyle.configure('TButton',background='#575757',foreground='white',font=self.font1,relief=FLAT)   #Estilo dos butões
        self.buttonstyle.map('TButton',background=[('active', '#646464')])                                          #Estilo dos butões
        self.grouppbuttonstyle = ttk.Style()                                                                        #Estilo dos butões de grupos pressionados
        self.grouppbuttonstyle.configure('pressed.TButton',background='#000080',foreground='white',font=self.font1,relief=SUNKEN)   #Estilo dos butões de grupos pressionados
        self.grouppbuttonstyle.map('pressed.TButton',background=[('active', '#0000ff')])                                           #Estilo dos butões de grupos pressionados
        self.groupubuttonstyle = ttk.Style()                                                                        #Estilo dos butões de grupos nao pressionados
        self.groupubuttonstyle.configure('unpressed.TButton',background='#575757',foreground='white',font=self.font1,relief=FLAT)   #Estilo dos butões de grupos nao pressionados
        self.groupubuttonstyle.map('unpressed.TButton',background=[('active', '#646464')])                                           #Estilo dos butões de grupos nao pressionados


        self.entrystyle=ttk.Style().configure("TEntry", bg='#2B2B2B', foreground="#ffffff",font=self.font1,borderwidth=0,relief=FLAT, width = 50)   #Estilo das Entries(não surtiu efeito quando exprimentei)

        self.listusr=StringVar()        #variavel que vai guardar o tipo de utilizador
        self.tipousers=['AAUE','Núcleos','Serviços','Catering','Conceções','Apoio Médico','Artista','Viatura','Live Act','UÉvora','Funcionários UÉ']
        self.nometext=StringVar()
        self.idtext=StringVar()

        #design#
        self.aaue=PhotoImage(file='logo.png')
        Label(self.master, image=self.aaue,bd=0).place(x=200,y=10)

        Label(self.master,text="Nome",font=('Helvetica',14),bg='#2B2B2B',fg='#ffffff').place(x=50,y=140)
        self.nome = Entry(self.master,textvariable=self.nometext,bg='#2B2B2B',foreground="#ffffff",font=self.font1,borderwidth=0,relief=FLAT)
        self.nome.place(x=50,y=160,height=40,width=400)

        Label(self.master, text="BI/CC/Matrícula",font=('Helvetica',14),bg='#2B2B2B',fg='#ffffff').place(x=50,y=230)
        self.bi = Entry(self.master,textvariable=self.idtext,bg='#2B2B2B',foreground="#ffffff",font=self.font1,borderwidth=0,relief=FLAT)
        self.bi.place(x=50,y=250,height=40,width=400)

        Label(self.master, text="Tipo de utilizador",font=('Helvetica',14),bg='#2B2B2B',fg='#ffffff').place(x=50,y=320)
        self.tipo = OptionMenu(self.master,self.listusr,*self.tipousers)
        self.tipo.config(bg='#2B2B2B',bd=0,relief=FLAT,fg='white',font=self.font1)
        self.tipo['menu'].config(bg='#2B2B2B',bd=0,relief=FLAT,fg='white',font=self.font1)
        self.tipo.place(x=50,y=340,height=40,width=400)

        Label(self.master, text="",font=('Helvetica',14),bg='#2B2B2B',fg='#ffffff').place(x=50,y=410)
        self.img = ttk.Button(self.master,style='TButton',text='Escolher Fotografia',command=self.getfoto)
        self.img.place(x=50,y=410,height=40,width=400)

        self.group1 = ttk.Button(self.master,style='unpressed.TButton',text="1",command=self.switch(1))
        self.group1.place(x=150,y=495,width=50,height=50)
        self.group2 = ttk.Button(self.master,style='unpressed.TButton',text="2",command=self.switch(2))
        self.group2.place(x=200,y=495,width=50,height=50)
        self.group3 = ttk.Button(self.master,style='unpressed.TButton',text="3",command=self.switch(3))
        self.group3.place(x=250,y=495,width=50,height=50)
        self.group4 = ttk.Button(self.master,style='unpressed.TButton',text="4",command=self.switch(4))
        self.group4.place(x=300,y=495,width=50,height=50)
        self.group5 = ttk.Button(self.master,style='unpressed.TButton',text="5",command=self.switch(5))
        self.group5.place(x=150,y=545,width=50,height=50)
        self.group6 = ttk.Button(self.master,style='unpressed.TButton',text="6",command=self.switch(6))
        self.group6.place(x=200,y=545,width=50,height=50)
        self.group7 = ttk.Button(self.master,style='unpressed.TButton',text="7",command=self.switch(7))
        self.group7.place(x=250,y=545,width=50,height=50)
        self.group8 = ttk.Button(self.master,style='unpressed.TButton',text="8",command=self.switch(8))
        self.group8.place(x=300,y=545,width=50,height=50)

        self.help=ttk.Button(self.master,style='unpressed.TButton',text='?',command=self.help).place(x=370,y=520,width=50,height=50)

        self.addtogroup = ttk.Button(self.master,text="Adicionar ao grupo",style='TButton',command=self.gerar).place(x=550,y=200,height=50,width=400)   #adicionar o command para a função gerar()
        self.adicionar = ttk.Button(self.master,text="Adicionar",style='TButton',command=self.add).place(x=275,y=620,width=175)
        sep=ttk.Separator(self.master,orient=VERTICAL).place(x=500,y=50,height=580)
        self.procuranome = ttk.Button(self.master,text="Procurar",style='TButton',command=self.searchbyname).place(x=550,y=100,height=50,width=175)    #adicionar o command para a função preview()
        self.procuraid = ttk.Button(self.master,text="Limpar",style='TButton',command=self.limpar).place(x=775,y=100,height=50,width=175)
        self.gerar=ttk.Button(self.master,text='Gerar',style='TButton',command=self.mergeall).place(x=50,y=620,width=175)
        countmerged=int(open('cnt1').read())
        countfiles=int(open('cnt2').read())
        nr=countfiles*4+countmerged
        self.nrgrupo=Label(self.master, text="Credenciais no grupo: "+str(nr),font=('Helvetica',14),bg='#2B2B2B',fg='#ffffff')
        self.nrgrupo.place(x=650,y=250)
        self.dellast= ttk.Button(self.master,text='Apagar Última',style='TButton',command=self.apagar).place(x=550,y=300,height=50,width=400)

    def mergeall(self):
        countmerged=int(open('cnt1').read())
        countfiles=int(open('cnt2').read())
        output=PdfFileWriter()
        for each in range(countfiles+1):
            inpu=PdfFileReader('credsA4/Credencial'+str(each+1)+'.pdf','rb')
            pagecount=inpu.getNumPages()
            for ipage in range(0,pagecount):
                output.addPage(inpu.getPage(ipage))
        outputStream=open('printready.pdf','wb')
        output.write(outputStream)
        outputStream.close
        for fil in os.scandir('credsA4/'):
            os.unlink(fil.path)
        open('cnt1','w').write('0')
        open('cnt2','w').write('0')
        nr=0
        self.nrgrupo['text']="Credenciais no grupo: "+str(nr)
        messagebox.showinfo("Sucesso","PDF pronto para impressão criado com sucesso")

    def getfoto(self):
        self.foto=filedialog.askopenfilename(filetypes=[('Image Files',("*.jpg","*.jpeg","*.png"))])

    def help(self):
        messagebox.showinfo("Ajuda para Zonas","Zona 1: Recinto\nZona 2: Bilheteira\nZona 3: SAFA\nZona 4: Central de Abastecimento\nZona 5: VIP\nZona 6: Frente de Palco\nZona 7: Backstage\nZona 8: Estacionamento")

    def add(self):
        alfa=db.adicionar(self.nometext.get(),self.idtext.get(),self.listusr.get())
        user=db.getbyall(alfa,self.nometext.get(),self.idtext.get(),self.listusr.get())
        if user!=[]:
            try:
                fotogra = Image.open(self.foto)
            except:
                pass
            else:
                fotogra=fotogra.resize((87,105))
                fotogra.save('fotos/'+self.idtext.get()+'.jpeg','jpeg')
            self.codalfa=alfa
            messagebox.showinfo("Sucesso","Utilizador "+self.nometext.get()+" criado/alterado com sucesso")

    def searchbyname(self):
        if self.nometext.get()!='' and self.idtext.get()!='':
            user=db.getbyboth(self.nometext.get(),self.idtext.get())
        elif self.idtext.get()=='' and self.nometext.get()!='':
            user=db.getbyname(self.nometext.get())
        elif self.idtext.get()!='' and self.nometext.get()=='':
            user=db.getbyid(self.idtext.get())

        if(user != []):
            self.nometext.set(user[0][1])
            self.idtext.set(user[0][2])
            self.listusr.set(user[0][3])
            self.codalfa=user[0][0]
            '''
            if user[0][3]=='AAUE':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
            elif user[0][3]=='Núcleos':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
                self.group4.configure(style='pressed.TButton')
                self.grupos[4]=1
            elif user[0][3]=='Serviços':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
                self.group3.configure(style='pressed.TButton')
                self.grupos[3]=1
                self.group4.configure(style='pressed.TButton')
                self.grupos[4]=1
                self.group5.configure(style='pressed.TButton')
                self.grupos[5]=1
                self.group6.configure(style='pressed.TButton')
                self.grupos[6]=1
                self.group7.configure(style='pressed.TButton')
                self.grupos[7]=1
                self.group8.configure(style='pressed.TButton')
                self.grupos[8]=1
            elif user[0][3]=='Catering':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
                self.group4.configure(style='pressed.TButton')
                self.grupos[4]=1
                self.group6.configure(style='pressed.TButton')
                self.grupos[6]=1
                self.group7.configure(style='pressed.TButton')
                self.grupos[7]=1
                self.group8.configure(style='pressed.TButton')
                self.grupos[8]=1
            elif user[0][3]=='Conceções' or user[0][2]=='Apoio Médico':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
                self.group8.configure(style='pressed.TButton')
                self.grupos[8]=1
            elif user[0][3]=='Artista':
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
                self.group6.configure(style='pressed.TButton')
                self.grupos[6]=1
                self.group7.configure(style='pressed.TButton')
                self.grupos[7]=1
            elif user[0][3]=='Viatura':
                self.group8.configure(style='pressed.TButton')
                self.grupos[8]=1
                '''
        else:
            messagebox.showinfo("Não Encontrado","Utilizador não encontrado na Base de dados")

    def limpar(self):
        self.nometext.set('')
        self.idtext.set('')
        self.listusr.set('')
        self.codalfa=None
        self.foto=None

    def switch(self,n):
        def wrapper(x=n):
            if x==1 and self.grupos[1]==0:
                self.group1.configure(style='pressed.TButton')
                self.grupos[1]=1
            elif x==1 and self.grupos[1]==1:
                self.group1.configure(style='unpressed.TButton')
                self.grupos[1]=0
            elif x==2 and self.grupos[2]==0:
                self.group2.configure(style='pressed.TButton')
                self.grupos[2]=1
            elif x==2 and self.grupos[2]==1:
                self.group2.configure(style='unpressed.TButton')
                self.grupos[2]=0
            elif x==3 and self.grupos[3]==0:
                self.group3.configure(style='pressed.TButton')
                self.grupos[3]=1
            elif x==3 and self.grupos[3]==1:
                self.group3.configure(style='unpressed.TButton')
                self.grupos[3]=0
            elif x==4 and self.grupos[4]==0:
                self.group4.configure(style='pressed.TButton')
                self.grupos[4]=1
            elif x==4 and self.grupos[4]==1:
                self.group4.configure(style='unpressed.TButton')
                self.grupos[4]=0
            elif x==5 and self.grupos[5]==0:
                self.group5.configure(style='pressed.TButton')
                self.grupos[5]=1
            elif x==5 and self.grupos[5]==1:
                self.group5.configure(style='unpressed.TButton')
                self.grupos[5]=0
            elif x==6 and self.grupos[6]==0:
                self.group6.configure(style='pressed.TButton')
                self.grupos[6]=1
            elif x==6 and self.grupos[6]==1:
                self.group6.configure(style='unpressed.TButton')
                self.grupos[6]=0
            elif x==7 and self.grupos[7]==0:
                self.group7.configure(style='pressed.TButton')
                self.grupos[7]=1
            elif x==7 and self.grupos[7]==1:
                self.group7.configure(style='unpressed.TButton')
                self.grupos[7]=0
            elif x==8 and self.grupos[8]==1:
                self.group8.configure(style='unpressed.TButton')
                self.grupos[8]=0
            elif x==8 and self.grupos[8]==0:
                self.group8.configure(style='pressed.TButton')
                self.grupos[8]=1
        return wrapper

    def gerar(self):
        if os.path.isfile('fotos/'+self.idtext.get()+'.jpeg'):
            self.foto='fotos/'+self.idtext.get()+'.jpeg'
        else:
            self.foto=None
        self.gerarA4()
        messagebox.showinfo("Sucesso", "Credencial "+str(self.codalfa)+" criada com sucesso.")

    def gerarA4(self):
        from reportlab.lib.units import mm
        countmerged=int(open('cnt1').read())
        countfiles=int(open('cnt2').read())
        fich='credsA4/Credencialtemp.pdf'
        fundo='fundo1.png' #escolher a imagem de fundo do pdf (sempre 296x420)
        self.c=canvas.Canvas(fich)
        if countmerged==4:
            countfiles+=1
            countmerged=0
        fname='credsA4/Credencial'+str(countfiles+1)+'.pdf'
        if countmerged==0:
            can=canvas.Canvas(fname)
            can.setLineWidth(1)
            can.line(297,0,297,842)
            can.line(0,421,595,421)
            can.showPage()
            can.line(297,0,297,842)
            can.line(0,421,595,421)
            can.showPage()
            can.save()
            x=0
            y=421.5
        elif countmerged==1:
            x=298
            y=421.5
        elif countmerged==2:
            x=0
            y=0
        elif countmerged==3:
            x=298
            y=0
        posx=150
        posy=255
        printready=PdfFileReader(fname,'rb')
        page1=printready.getPage(0)
        self.c.setStrokeColorRGB(1,1,1)
        self.c.drawImage(fundo,x,y)
        self.c.setFillColorRGB(1,1,1)
        self.c.rect(x+30,y+144,237,20, fill=1)             #caixa para o nome
        self.c.rect(x+30,y+106,237,20, fill=1)             #caixa para o bi/cc/Matricula
        self.c.rect(x+140,y+195,140,20, fill=1)            #caixa para tipo de utilizador
        self.c.rect(x+30,y+18,100,18, fill=1)             #caixa para o codigo alfanumerico
        self.c.setFont("Helvetica-Bold", 12)
        for i in self.grupos:
            if self.grupos[i]==1:
                if i==1:
                    self.c.setFillColorRGB(0.14,0.48,0.63)
                    self.c.setStrokeColorRGB(0.14,0.48,0.63)
                elif i==2:
                    self.c.setFillColorRGB(1,0,0)
                    self.c.setStrokeColorRGB(1,0,0)
                elif i==3:
                    self.c.setFillColorRGB(0.6,0.77,0.24)
                    self.c.setStrokeColorRGB(0.6,0.77,0.24)
                elif i==4:
                    self.c.setFillColorRGB(1,0.88,0.4)
                    self.c.setStrokeColorRGB(1,0.88,0.4)
                elif i==5:
                    self.c.setFillColorRGB(0.98,0.47,0.13)
                    self.c.setStrokeColorRGB(0.98,0.47,0.13)
                elif i==6:
                    self.c.setFillColorRGB(0.95,0.37,0.36)
                    self.c.setStrokeColorRGB(0.95,0.37,0.36)
                elif i==7:
                    self.c.setFillColorRGB(0.44,0.76,0.7)
                    self.c.setStrokeColorRGB(0.44,0.76,0.7)
                elif i==8:
                    self.c.setFillColorRGB(0.46,0.31,0.27)
                    self.c.setStrokeColorRGB(0.46,0.31,0.27)
                self.c.rect(x+posx,y+posy,30,30,fill=1)
                self.c.setFillColorRGB(0,0,0)
                self.c.drawString(x+posx+11,y+posy+11,str(i))
            elif self.grupos[i]==0:
                self.c.setStrokeColorRGB(0.11,0.21,0.34)
                self.c.setFillColorRGB(0.11,0.21,0.34)
                self.c.rect(x+posx,y+posy,30,30,fill=1)
                self.c.setFillColorRGB(1,1,1)
                self.c.drawString(x+posx+11,y+posy+11,'X')
            if i!=4:
                posx+=30
            else:
                posx-=90
                posy-=30
        self.c.setStrokeColorRGB(0,0,0)
        self.c.setFillColorRGB(0,0,0)
        usrfoto = ['AAUE','Apoio Médico','Catering','Núcleos','Live Act','UÉvora','Funcionários UÉ']
        if self.listusr.get()=='Artista':
            self.c.drawString(x+33,y+128, 'Válido para os dias:')
        elif self.listusr.get()=='Viatura':
            self.c.drawString(x+33,y+128,'Matricula:')
        elif self.listusr.get() in usrfoto:
            self.c.drawString(x+33,y+128,'BI/CC:')
            try:
                Image.open(self.foto)
            except:
                pass
            else:
                self.c.setStrokeColorRGB(1,1,1)
                self.c.setFillColorRGB(1,1,1)           #caixa para a foto
                self.c.drawImage(self.foto,x+30,y+180)
        else:
            self.c.drawString(x+33,y+128,'BI/CC:')

        self.c.setStrokeColorRGB(0,0,0)
        self.c.setFillColorRGB(0,0,0)

        qr=pyqrcode.create(self.codalfa.upper())
        qr.png('qrtemp.png',scale=2,module_color=(0,0,0),background=(255,255,255,0))

        self.c.drawImage('qrtemp.png',x+47,y+31)
        self.c.drawCentredString(x+80,y+22,self.codalfa.upper())

        self.c.drawString(x+33,y+167,'Nome:')
        self.c.setFont("Helvetica-Bold", 14)

        self.c.drawString(x+35,y+148,self.nome.get().upper())
        self.c.drawString(x+35,y+110,self.bi.get())
        self.c.drawCentredString(x+210,y+199,self.listusr.get().upper())
        self.c.showPage()
        self.c.save()
        temp=PdfFileReader(open('credsA4/Credencialtemp.pdf','rb'))
        page2=temp.getPage(0)
        page1.mergePage(page2)
        output=PdfFileWriter()
        output.addPage(page1)
        os.unlink('qrtemp.png')
        os.unlink('credsA4/Credencialtemp.pdf')
        countmerged+=1
        output.write(open('credsA4/Credencial'+str(countfiles+1)+'.pdf','wb'))
        open('cnt1','w').write(str(countmerged))
        open('cnt2','w').write(str(countfiles))
        nr=countfiles*4+countmerged
        self.nrgrupo['text']="Credenciais no grupo: "+str(nr)

    def apagar(self):
        countmerged=int(open('cnt1').read())
        countfiles=int(open('cnt2').read())
        nr=countfiles*4+countmerged-1
        self.nrgrupo['text']="Credenciais no grupo: "+str(nr)
        if countmerged==0:
            if countfiles>0:
                countfiles-=1
                countmerged=3
                open('cnt2','w').write(str(countfiles))
                open('cnt1','w').write(str(countmerged))
        else:
            countmerged-=1
            open('cnt1','w').write(str(countmerged))
        if countmerged==0:
            os.unlink('credsA4/Credencial'+str(countfiles+1)+'.pdf')
        elif countmerged==1:
            x=298
            y=421.5
        elif countmerged==2:
            x=0
            y=0
        elif countmerged==3:
            x=298
            y=0
        fich='credsA4/Credencialtemp.pdf'
        self.c=canvas.Canvas(fich)
        fname='credsA4/Credencial'+str(countfiles+1)+'.pdf'
        printready=PdfFileReader(fname,'rb')
        page1=printready.getPage(0)
        self.c.setStrokeColorRGB(0,0,0)
        self.c.setFillColorRGB(1,1,1)
        self.c.rect(x,y,297,421.5,fill=1)
        self.c.showPage()
        self.c.save()
        temp=PdfFileReader(open('credsA4/Credencialtemp.pdf','rb'))
        page2=temp.getPage(0)
        page1.mergePage(page2)
        output=PdfFileWriter()
        output.addPage(page1)
        os.unlink('credsA4/Credencialtemp.pdf')
        output.write(open('credsA4/Credencial'+str(countfiles+1)+'.pdf','wb'))


def main():
    root = Tk()
    root.configure(background='#2B2B2B')
    root.title('Gestor de Credenciais')
    root.geometry('1000x680')
    root.resizable(width=False,height=False)
    app = Application(master=root)
    root.mainloop()

if __name__ == '__main__':
    main()
