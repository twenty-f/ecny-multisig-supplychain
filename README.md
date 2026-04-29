# 🏭 e-CNY Supply Chain Multi-Sig | 供应链金融阈值多签智能合约

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Threshold_Sig-ff69b4)
![Security](https://img.shields.io/badge/Security-RBAC-critical)

## 📖 项目背景
在大宗商品采购（如千万级化肥供应链）中，传统的银团信用证手续繁杂且成本高昂。本项目利用数字人民币底层可编程性，实现了一个**基于多方共识的去中心化资金拨付系统**。仅当物流、质检等关键节点完成数字签名并达到预设阈值时，千万级资金才会自动瞬间划拨，彻底解决供应链“三角债”与信任危机。

## 🏗️ 核心架构与亮点
- **基于角色的访问控制 (RBAC) 与白名单验证**：采用二维嵌套 Map `map[RoleType]map[string]bool` 物理隔离不同角色（物流/质检）的权限边界，精准防御越权签名。
- **M-of-N 阈值触发机制 (Threshold Trigger)**：支持高度可定制的多签策略（如“物流 >=1 且 质检 >=1 且 总数 >=3”）。摒弃手动拉取打款（Pull Payment），实现状态满足瞬间的自动事件驱动放款。
- **严密的全量状态机防御**：构建 `PendingFund` (待汇款锁仓) -> `Signing` (授权收集中) -> `Paid` (已拨付) 的完整生命周期。
- **工程化依赖注入 (Dependency Injection)**：重构 `InitContract`，强制通过参数传入阈值与白名单，完美规避 Go 语言特有的结构体“零值陷阱 (Zero Value)”，提升合约工厂的通用性。
- **防重放攻击 (Replay Attack Prevention)**：在签名校验层实现 O(1) 级别的幂等性检查，阻断单节点恶意重复签名刷单。