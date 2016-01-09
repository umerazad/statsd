import statsd

sc = statsd.StatsClient(host='localhost', port=1119, prefix=None, maxudpsize=8*1024)


for i in range(0, 15):
    sc.incr('counter', count=1, rate=1)
    sc.gauge('guage', i)
    sc.timing('timer', i, rate=1)

