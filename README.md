# rank

排行榜系用采用CQRS命令和结果分离的架构，简单分为两个服务。一个为计算逻辑服务（job），一个为结果服务（service）。
计算逻辑服务主要用来计算榜单的分数属于命令端，基于一些用户可以累加分数的行为来加分/减分。结果服务主要用来展示活动数据，排行榜、用户名次、得分。两端逻辑互不影响，公用底层存储。

计算逻辑服务主要作为consumer来消费消息，通过上游的消息来计算分数。结果服务只是单纯的展示。

app(mq/grpc)->job->(mysql/redis)<-service