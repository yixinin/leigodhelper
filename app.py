# import math
# 极限思想求三角形面积

d = 10
h = 5
step = 10000000

delta = 5 / step

s = 0
for i in range(step):
    s = s + (d/step) * (i*delta)

print(s, d*h / 2)
