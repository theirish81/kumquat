import sys
import json

data = ''

for line in sys.stdin:
    data = data+line

parsed_data = json.loads(data)
num = parsed_data["params"]["number"]
flag = False
if num > 1:
    for i in range(2, num):
        if (num % i) == 0:
            flag = True
            break
if flag:
    print(num, "is not a prime number")
else:
    print(num, "is a prime number")