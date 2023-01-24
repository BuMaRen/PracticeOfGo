1. server解析不了body中的中文
    * data中一开始声明struct的时候用的是bson，但是unmarshal用的是json所以字段不一样的内容匹配不了
2. 回写数据给curl的时候不能写bson，curl解不了。收到客户端的消息全部用json解码