# sweb
golang + redis + mongo

1.定义结构体，并在主函数中注册，自动注册增、删、查、改的接口（get, post, put, delete），支持自定义选择哪些接口可用
2.Get支持单id，多id查询和多字段模糊查询
3.insert id默认自增1
4.用户注册生成token，特定的接口需要用户验证
