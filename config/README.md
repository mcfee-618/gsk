# config 配置解析

## TODO:需要支持的特性

- 常用的配置通常分为两种:类json(yaml,ini,hcl)格式,csv表格式
- 数据来源通常分为本地磁盘,etcd,k8s等
- 需要支持监控配置变化
- 需要支持更丰富的字段类型,比如Slice,Map,Time,Week等类型
- 配置通常允许动态加载,而配置的发布通常又是全量发布,因此配置解析前需要判断是否发生过变化
- 对于csv格式的表格,外部使用时有可能会持有某个条目的指针,从而避免每次查询,因此表格需要支持版本校验
- 服务器运行时，可能需要查询当前配置以及md5等信息,用于校验配置是否正确,或者用于gm展示
- 在微服务架构中,不同的服务可能需要不同的表,但导表工具一般会放到一个文件夹中,可以手动管理表格加载
- 不同环境的区分
  - 一种情况下我们希望线上环境能和测试环境保持一致的配置,比如我们新添加了一个配置项,希望在发版的时候能自动同步到线上环境
  - 另外一种情况我们希望不同环境使用不同的配置,比如log日志级别,不同环境配置不一样,再比如数据配置

## 其他

- [taobao config-center](http://jm.taobao.org/2016/09/28/an-article-about-config-center/)