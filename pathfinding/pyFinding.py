from pprint import pprint
from os import system


'''
Consts for map
'''
#Edit this if you got a diferent layout
DOOR = "|"
WALL = "#"
FLOOR =  "."


def read_map(path: str) -> []:
    data = []
    with open(path,"r") as f:
        l = f.readline()
        while l:
            #Read line and 
            tmp = list(l)
            tmp = tmp[:-1] if tmp[-1] == '\n' else tmp
            data.append(tmp)
            l = f.readline()
    return data

path = "game/maps/copy.map"

m = read_map(path)


def renderMatrix(mapData, texture,visited):
    for i,row in enumerate(mapData):
        for j, c in enumerate(row):
            #Print for path
            #Print for everything else
            print(
                texture[c],
                end = ' ',
            )
        print()


def clear():
    system("clear")

def valid(x,y, m, visited):
    if x < 0 or x >= len(m) or y < 0 or y >= len(m[0]):
        return False
    return (x,y) not in visited and m[x][y] != "#"

def dfs(x,y, m, texture_map):

    #Visited
    visited = {(x,y)}

    #Hash table for parents
    parents = {

    }
    
    stack = [(x,y)]

    goal = -1,-1


    while stack:
        #Extract the coords
        i,j = stack[0]
        
        if m[i][j] == "|":
            goal = (i,j)
            break

        #neighbours
        up = (i+1,j)
        down = (i-1,j)
        left = (i,j-1)
        right = (i,j+1)
        upright = (i+1, j+1)
        upleft = (i+1, j-1)
        downleft = (i-1, j-1)
        downright = (i-1, j+1)

        #Test for valid up,down, left and right blocks
        #diagonals
        d = [up,upright,right, downright, down,downleft,left, upleft]
        for _, n in enumerate([up, right, down, left]):
            if (valid(*n, m, visited)):
                #Draw visited on matrix
                m[n[0]][n[1]] = "+"
                visited.add(n)
                stack += [n]
                #print(n, [i,j], _)
                parents[n] = (i,j)
        

        #print(stack)
        stack = stack[1:]

        #Render each iteration to visualize BFS
        renderMatrix(m, texture_map, visited)
        clear()

    #End point
    end_point = goal


    for e in visited:
        m[e[0]][e[1]] = "+"


    #Draw path
    while goal != (x,y):
        print(goal)
        goal = parents[goal]
        m[goal[0]][goal[1]] = "0"

    #Draw starting point

    m[x][y] = "A"

    m[end_point[0]][end_point[1]] = "B"

    renderMatrix(m, texture_map,visited)




def main():
    texture_map = {
        "#": "▫️",
        ".": " ",
        "|": "🚪",
        "A": "A",
        "B": "B",
        "+": "+",
        "0": "0",
    }
    while 1:
        try:
            coords  = map(int,input("Coords (x y): ").split())
        except:
            print("Invalid coords \nex: 1 2")
            continue
        
        dfs(*coords, m, texture_map)


if __name__ == "__main__" :
    main()