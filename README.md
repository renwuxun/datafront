#### datafront
数据前端缓存服务，基于groupcache

#### 使用
启动两个实例：
```bash
./datafront --me=127.0.0.1:8080 --others=127.0.0.1:8081 &
./datafront --me=127.0.0.1:8081 --others=127.0.0.1:8080 &
```
向任意实例取数据，实例内部会自动分发请求
```bash
curl '127.0.0.1:8080/front/get?group=dummygroup&key=abc'
curl '127.0.0.1:8080/front/get?group=dummygroup&key=abc1'
```
输出：
```bash
key: [abc-0] from peer: 127.0.0.1:8081
key: [abc1-0] from peer: 127.0.0.1:8080
```
清除缓存：
```bash
curl '127.0.0.1:8080/front/purge?group=dummygroup'
curl '127.0.0.1:8081/front/purge?group=dummygroup'
```