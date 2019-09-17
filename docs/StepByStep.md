# 项目开发记录

##项目搭建过程
1.使用go modules进行包管理
具体参考https://juejin.im/post/5d63c1b76fb9a06b1e7f48e1
```
go mod init GoRestServer
```
```
go mod tidy  // 拉取缺少的模块
go mod vendor // 将依赖复制到vendor下
go mod download // 下载依赖包
go mod verify // 检验依赖
go mod graph // 打印模块依赖图
```

2.使用gin作为基础web框架
```
govendor fetch github.com/gin-gonic/gin
```

3.配置文件
```
govendor fetch github.com/jinzhu/configor
```
CONFIGOR_DEBUG_MODE: 设置配置文件加载调试模式环境变量
通过设置环境变量可修改相关配置,前缀默认为TSS_, 比如
```
TSS_APP_RUNMODE=release
==
App.RunMode="release"
```

4. 使用gorm进行数据库管理

5. json查询条件与SQL where对应关系
```
a: b => a = b
a: {'gt': b} => a > b
a: {'lt': b} => a < b
a: {'gte': b} => a >= b
a: {'lte': b} => a <= b
a: {'not': b} => a <> b
a: {'between': [b,c]} => a between b and c
a: {'like': b} => a like b
a: {'in': [b,c,d] => a in [b,c,d]
```

6. 使用JOIN以支持多表关联查询和排序
 
 参考UserModel。
 json字段名称和table字段名称要一致，以方便进行不定条件查询和排序等
 
 7. 使用Swagger实现API文档
 
 8. 单表自我关联的实现问题，见Function表(TODO)
 
 9. 认证与授权
 - JWT
 - 如何根据角色功能自动授权
 
 10. Redis缓存

 11. File Upload/Download
