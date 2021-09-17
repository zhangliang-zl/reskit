1 增加自动检测机制，参考 https://github.com/pangudashu/memcache

case1 GET 时当前server节点不可用时，使用下一个节点 , 直到所有的节点都遍历？
client 对每个server进行检测
case2 动态维护所有节点是否可用 心跳检测， 每隔5s检测一下， 不可用剔除，可用加入进来
减低锁开销