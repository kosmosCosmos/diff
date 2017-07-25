# hots
热点数据库的帖子在一个小时内增加了多少写的代码

对 [hot-title](https://github.com/kosmosCosmos/hot-title)的补充，[hot-title](https://github.com/kosmosCosmos/hot-title)该项目只是抓取了每个帖子的信息，
但并没有给出每个帖子每段时间内的增长数，我通过[xorm](https://github.com/xormplus/xorm)读取出数据然后与下一个时间段的数据相互比较最终得出结果。
  - 优点：
  - 对每个帖子都进行时间判定，避免坟贴
  - 可配置的文件
  
