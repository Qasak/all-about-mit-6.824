## 可扩展性

以twitter为例

实现扇出（fan-out）——每个用户关注了很多人，也被很多人关注。

两种方式：

+ 发布推文时，只需将新推文插入全局推文集合即可。当一个用户请求自己的主页时间线 时，首先查找他关注的所有人，查询这些被关注用户发布的推文并按时间顺序合并。在 如图1-2所示的关系型数据库中，可以编写这样的查询：

  ```sql
  SELECT tweets.*, users.*
      FROM tweets
      JOIN users ON tweets.sender_id = users.id
      JOIN follows ON follows.followee_id = users.id
      WHERE follows.follower_id = current_user
  
  ```

  