import sys

# input = sys.stdin.read
# data = input().split()

data = [
    "8",
    "1",
    "0",
    "5",
    "2",
    "1",
    "3",
    "1",
    "1",
    "0",
    "6",
    "1",
    "0",
    "7",
    "2",
    "2",
    "1",
    "1",
    "5",
    "2",
    "2",
]

q = int(data[0])
index = 1

# Изначально пустой список
lst: list = []

for i in range(q):
    typ = int(data[index])
    index += 1

    if typ == 1:
        # Запрос 1: добавить y после x-го элемента (x=0 -> в начало)
        x = int(data[index])
        y = int(data[index + 1])
        index += 2

        if x == 0:
            lst.insert(0, y)
        else:
            lst.insert(x + 1, y)

    elif typ == 2:
        # Запрос 2: вывести элемент на позиции x
        x = int(data[index])
        index += 1
        print(lst[x - 1])  # Индексация с 0, позиции с 1

    elif typ == 3:
        # Запрос 3: удалить элемент на позиции x
        x = int(data[index])
        index += 1
        del lst[x - 1]
