f = open("file.txt", "r")
positions=[]
x=0
for line in f:
    data=list(line)
    for y in range(0,len(data)):
        if(data[y]=='#'):
            positions.append([x,y])
    x+=1

#print(positions)
#[]int{1,2,3},[]int{1,2,3}

res=""
l='{'
r='}'

#a.FillTile(astar.Point{1, 1}, -1)
for element in positions:
    res+="a.FillTile(astar.Point{"+str(element[0])+","+str(element[1])+ "}, "+"-1)"+"\n"
print(res)