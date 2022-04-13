# CronJob operator

## API  design
通常来说，CronJob 由以下几部分组成：

    * 一个时间表（ CronJob 中的 cron ）
    * 要运行的 Job 模板（ CronJob 中的 Job ）

当然 CronJob 还需要一些额外的东西，使得它更加易用

    * 一个已经启动的 Job 的超时时间（如果该 Job 执行超时，那么我们会将在下次调度的时候重新执行该 Job）。
    * 如果多个 Job 同时运行，该怎么办（我们要等待吗？还是停止旧的 Job ？）
    * 暂停 CronJob 运行的方法，以防出现问题。
    * 对旧 Job 历史的限制
