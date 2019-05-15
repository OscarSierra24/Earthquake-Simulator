from pprint import pprint

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

path = "../game/maps/copy.map"

m = read_map(path)


def valid(x,y, m, visited):
    if x < 0 or x >= len(m) or y < 0 or y >= len(m[0]):
        return False
    return (x,y) not in visited and m[x][y] != "#"

def dfs(x,y, m):

    #Visited
    visited = {(x,y)}

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

        #Neightbors
        up = (i+1,j)
        down = (i-1,j)
        left = (i,j-1)
        right = (i,j+1)

        #Test for valid up,down, left and right blocks
        for _, n in enumerate([up, down, left, right]):
            if (valid(*n, m, visited)):
                visited.add(n)
                stack += [n]
                #print(n, [i,j], _)
                parents[n] = (i,j)
        

        #print(stack)
        stack = stack[1:]

    while goal != (x,y):
        print(goal)
        goal = parents[goal]
        m[goal[0]][goal[1]] = "0"

    print(*["".join(r) for r in m], sep='\n')

dfs(1,25, m)