# 简单的api框架

## 安装相关包
使用go get 命令

```sh
github.com/aiwuTech/fileLogger
github.com/go-xorm/xorm
github.com/gorilla/mux
github.com/kylelemons/go-gypsy/yaml
github.com/go-sql-driver/mysql
github.com/garyburd/redigo/redis
```

### 修改配置文件
复制一份conf.yaml.example为conf.yaml

YAML 支持的数据结构有三种
1 对象：键值对的集合，又称为映射（mapping）/ 哈希（hashes） / 字典（dictionary）
2 数组：一组按次序排列的值，又称为序列（sequence） / 列表（list）
3 纯量（scalars）：单个的、不可再分的值

注：配置时空格问题
直接推荐使用对象形式：
```sh
host: 8080  

hash: { name: Steve, foo: bar } 

db:
 host: 115.29.5.235
 port: 3306
```


## Running

To run this exmaple, from the root of this project:

```sh
go run ./main.go
```
