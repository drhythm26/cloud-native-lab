# 错题本

本文件只收错题；实验复盘见各组件 `note/`。

## 01 Prometheus/prometheusrule：进程挂死时 up 的值与来源

- 日期：2026-07-19
- 我的答案：up=1，是 service 产生的
- 正确要点：up=0。up 是 Prometheus 抓取器自报的合成指标——每次抓取后由 Prometheus 自己写入，成功记 1、失败记 0，与应用和 Service 无关
- 错因：把 up 当成被监控方上报的指标，方向弄反了；up 恰恰是「对方不回话时我自己记一笔」的机制
- 复测：未

## 02 Prometheus/prometheusrule：监控对象整体消失时告警失效，absent 兜底

- 日期：2026-07-19
- 我的答案：知道 Down 不能触发，但说不出机制，也不知道兜底方案（两问均答「不知道」）
- 正确要点：ServiceMonitor 删除 → target 消失 → up 序列经 staleness 蒸发 → `up == 0` 无输入、永不为真。「序列不存在」和「序列值为 0」是两种故障面，`==0` 只覆盖后者；前者用 `absent(up{job="go-api"})` 兜底
- 错因：实验里讲过「up==0 vs 序列蒸发」两种覆盖面，但未内化成可复述的机制和对策
- 复测：未

## 03 Prometheus/prometheusrule：release label 是规则加载的暗号；缺 for 的后果

- 日期：2026-07-19
- 我的答案：只找出 severity/summary 换位（3 处错找到 1 处），且说不出后果
- 正确要点：① metadata.labels 缺 `release: prometheus` → ruleSelector 匹配不到，规则整个不被加载，/alerts 页面上根本看不到（最隐蔽）；② 缺 `for` → 表达式一次为真立即 firing，瞬时毛刺就打扰人；③ severity 必须在 labels（参与告警身份和 Alertmanager 路由），summary 在 annotations（仅展示渲染）
- 错因：release label 在实验中专门讲过（加载链路排查的关键），改错场景下未能迁移；对 for 的默认值行为（无 for = 立即 firing）没有概念
- 复测：未

## 04 Prometheus/prometheusrule：空结果不参与任何比较

- 日期：2026-07-19
- 我的答案：a 答对「空」，但说不出为什么不会误触发（答「不知道」）
- 正确要点：空向量不参与 `>` / `<` 任何比较运算——告警表达式无输出 = 条件不成立，所以基线安静是「无数据」而非「0 < 10」。推论：想告「指标消失/太少」，`< N` 无效，要用 `absent()` 或 `or vector(0)` 补零
- 错因：「空 vs 0」已第四次出现——知道结果是空，但没建立「空不进入比较运算」这条规则，停在现象层没到机制层
- 复测：未
