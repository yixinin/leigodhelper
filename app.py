# import math
# 极限思想求三角形面积

d = 10
h = 5
step = 100

delta = h / step
x = d / delta


s = 0
for i in range(step):
    s = s + (d/step) * (i*delta)
# s = (d / st) * 1*(h / st) +  (d / st) * 2*(h / st)  + ...
# s = (d/st) * (h/st) * (1 + 2 + ...)
# s = (d*h)/(st^2)* (st*(st-1)/2)
# s = (d*h)/st^2 * (st^2-st)/2
# s = (d*h) * 1/2 -(d*h)/st^2 *st/2
# s = d*h/2 - (d*h)/2st 极限(st -> +int.max, (d*h)/2st -> 0)

print(s, d*h / 2)

print("{}".format((d*h)/(2*step)))
