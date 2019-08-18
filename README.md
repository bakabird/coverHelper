# 1. 创建数据库

保证安装了sqlite3，通过`sqlite3 url2Path.db`，创建一个sqlite数据库。



* 管理员可以给某个条目添加自定义的cover(仅本地上传)，服务器在接受到后将会保存在服务器中。



# A Socket Server



1. 使用sqlite，可配置。
2. 管理的bilibiliCover的保存文件夹；
3. 监听的端口。

## C服务(CHECK)

```mermaid
sequenceDiagram
    participant CH_database
    participant CH_socketServer
    participant thinkblogServer
    thinkblogServer->>CH_socketServer: "C" + 以空格分隔的各个URL
    loop 每个URL
      CH_socketServer->>CH_database: db中有以URL为主键的项么
      CH_database->>CH_socketServer: 返回结果
      alt 不存在
        CH_socketServer->>CH_socketServer: 将URL对应封面下载到图片仓库
        CH_socketServer->>CH_database: 保存对应的[url]path项到DB
      end
    end
    CH_socketServer->>thinkblogServer: OK 
```

## 不同的Conn都需要下载同一个URL

```mermaid
sequenceDiagram
    participant db
    participant 	dealSave[A]
    participant 	dealSave[B]
    participant 	dealSave[C]
    participant 	dealSave[D]
    participant 	dealSave[E]

	dealSave[A]->>db: 某个未下载的URL是否存在
	dealSave[B]->>db: 某个未下载的URL是否存在
  dealSave[E]->>db: 某个未下载的URL是否存在
	db->>dealSave[A]: 不存在
  db->>dealSave[B]: 不存在
    db->>dealSave[E]: 不存在
  dealSave[A]->>DownloadActionManager: 申请“下载任务锁”
DownloadActionManager->>dealSave[A]: 申请到锁
  dealSave[B]->>DownloadActionManager: 申请“下载任务锁”
  dealSave[A]->>dealSave[A]: 发现其中还没有这个任务：添加任务，并设为【执行中】
dealSave[A]->>DownloadActionManager: 归还锁
  dealSave[A]->>dealSave[A]: 去完成相应URL的下载
DownloadActionManager->>dealSave[B]: 申请到锁
  dealSave[B]->>dealSave[B]: 发现任务处于【执行中】
 dealSave[B]->>DownloadActionManager: 归还锁
 dealSave[E]->>DownloadActionManager: 申请“下载任务锁”
  dealSave[B]->>DownloadActionManager: Wait对应任务的<任务完成>
  dealSave[A]->>dealSave[A]: 完成下载
  dealSave[C]->>db: 某个未下载的URL是否存在
  db->>dealSave[C]: 不存在
  Note right of  dealSave[C]: 随后的行为<br/>与dealSave[B]<br/>一样...
  dealSave[A]->>dealSave[A]: 完成了往db中的录入
dealSave[A]->>DownloadActionManager: 申请“下载任务锁”
DownloadActionManager->>dealSave[A]: 申请到锁
  dealSave[A]->>dealSave[A]: 设置相应任务为【完成】
dealSave[A]->>DownloadActionManager: 归还锁
  dealSave[D]->>db: 某个未下载的URL是否存在
  db->>dealSave[D]: 存在了
  dealSave[D]->>dealSave[D]: 任务结束
DownloadActionManager->>dealSave[E]: 申请到锁
dealSave[E]->>dealSave[E]: 发现任务处于【完成】
dealSave[E]->>DownloadActionManager: 归还锁
dealSave[E]->>dealSave[E]: 任务结束
  dealSave[B]->>dealSave[B]: 任务结束
  dealSave[A]->>DownloadActionManager: Brodecast对应任务的<任务完成>
  dealSave[A]->>dealSave[A]: 任务结束
  DownloadActionManager->>dealSave[B]: <任务完成>
  dealSave[B]->>dealSave[B]: 任务结束
```

