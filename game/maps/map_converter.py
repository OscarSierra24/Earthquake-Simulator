"""
author: quirinoc 

Converts map in format
########
#......|
#...#..#
########
To
11111111
10000002
10001001
11111111
For easier matrix usage


0 Floor
1 Wall
2 Door
"""



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

def transform_map(map_data: [], to_map: {}) -> []:
    for i, row in enumerate(map_data):
        for j, col in enumerate(row):
            print(to_map, map_data[i][j])
            map_data[i][j] = to_map[map_data[i][j]]

def save_map(path: str, map_data: []):
    with open(path, "w") as f:
        for row in map_data:
            f.writelines(
                "".join(row) + "\n"
            )

if __name__ == "__main__" :
    import sys

    if len(sys.argv) < 3:
        print("Usage err: python map_converter.py inputfile outputfile")

    file_input = sys.argv[1]
    file_output = sys.argv[2]

    m = read_map(file_input)

    mask = {
        "#": "1",
        ".": "0",
        "|": "2"
    }
    transform_map(m, mask)

    save_map(file_output,m)


