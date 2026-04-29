package main

import (
	"errors"
	"fmt"
)

// RoleType 表示签名角色类型。
type RoleType int

const (
	RoleLogistics RoleType = iota // 0: 物流
	RoleQuality                   // 1: 质检
)

func (r RoleType) String() string {
	switch r {
	case RoleLogistics:
		return "Logistics"
	case RoleQuality:
		return "Quality"
	default:
		return "Unknown"
	}
}

// ContractStatus 表示合约当前状态。
type ContractStatus string

const (
	StatusPending   ContractStatus = "Pending"   //等待付款
	StatusSigning   ContractStatus = "Signing"   //等待签名
	StatusCompleted ContractStatus = "Completed" //合约已完成
	StatusRefunded  ContractStatus = "Refunded"  //合约未完成已回款
)

// MultiSigContract 供应链金融多签合约骨架。
type MultiSigContract struct {
	ContractID string
	Amount     int

	// 白名单：按角色维护可签名主体
	Whitelist map[RoleType]map[string]bool
	// 签名记录：按角色维护已签名主体
	Signatures map[RoleType]map[string]bool

	// 当前计数
	LogisticsCount int
	QualityCount   int
	TotalCount     int

	// 阈值配置
	LogisticsMin int
	QualityMin   int
	TotalMin     int

	Status ContractStatus
}

// InitContract 初始化合约与白名单，设置 1+1+3 阈值。
func InitContract(contractID string, amount int, logisticsmin int, qualitymin int, totalmin int) *MultiSigContract {
	m := &MultiSigContract{
		ContractID: contractID,
		Amount:     amount,
		Whitelist: map[RoleType]map[string]bool{
			RoleLogistics: {
				"A": true,
				"B": true,
				"C": true,
			},
			RoleQuality: {
				"a": true,
				"b": true,
				"c": true,
			},
		},
		Signatures: map[RoleType]map[string]bool{
			RoleLogistics: {},
			RoleQuality:   {},
		},
		LogisticsCount: 0,
		QualityCount:   0,
		TotalCount:     0,
		LogisticsMin:   logisticsmin,
		QualityMin:     qualitymin,
		TotalMin:       totalmin,
		Status:         StatusPending,
	}

	return m
}

func (m *MultiSigContract) Deposit(amount int) error {
	// 拦截 1：状态不对
	if m.Status != StatusPending {
		return errors.New("❌ 汇款失败：合约当前不处于待汇款状态")
	}

	// 拦截 2：金额不对（防止少打钱或者多打钱扯皮）
	if amount != m.Amount {
		return fmt.Errorf("❌ 汇款失败：金额不匹配。需锁仓 %d，实际支付 %d", m.Amount, amount)
	}

	// 资金到位，状态流转
	m.Status = StatusSigning
	fmt.Printf("💰 【资金锁仓成功】采购方已将 %d 元数字人民币注入合约 eCNY-MSIG-SC-20260429-0001！\n", amount)
	fmt.Println("⏳ 合约状态变更为：待签名 (StatusSigning)，请各节点开始授权...")
	return nil
}

func (m *MultiSigContract) Sign(name string, role RoleType) error {
	// 合约处于待签名
	if m.Status != StatusSigning {
		return errors.New("❌ 验证失败：合约不在待签名状态")
	}
	// 签名者是白名单
	if !m.Whitelist[role][name] {
		return errors.New("❌ 验证失败：你不是该角色的合法节点")
	}
	// 签名者不是已签名节点
	if m.Signatures[role][name] {
		return errors.New("❌ 验证失败：你在该合约中已经签过名了")
	}
	// 确认无误，添加签名记录
	m.Signatures[role][name] = true
	m.TotalCount++
	if role == RoleLogistics {
		m.LogisticsCount++
	} else {
		m.QualityCount++

	}
	if m.LogisticsCount >= m.LogisticsMin && m.QualityCount >= m.QualityMin && m.TotalCount >= m.TotalMin {
		m.Status = "StatusPaid" // 状态瞬间变为已打款
		fmt.Println("🎉 触发智能合约：多签阈值已达标，1000万数字人民币自动划拨至供应商账户！")
	}
	return nil
}

// 手动验证是否可以打款
func (m *MultiSigContract) CheakAndPay() error {
	// 合约处于待签名
	if m.Status != StatusSigning {
		return errors.New("❌ 验证失败：合约不在待签名状态")
	}
	if m.LogisticsCount < m.LogisticsMin {
		return errors.New("❌ 验证失败：物流节点签名量不足")
	}
	if m.QualityCount < m.QualityMin {
		return errors.New("❌ 验证失败：质检节点签名量不足")
	}
	if m.TotalCount < m.TotalMin {
		return errors.New("❌ 验证失败：总节点签名量不足")
	}
	m.Status = StatusCompleted
	// 执行打款操作
	fmt.Println("🎉 触发智能合约：多签阈值已达标，1000万数字人民币自动划拨至供应商账户！")
	return nil
}

func main() {
	// 1. 初始化你的合约 (这里假设你已经写好了 Init 函数，或者直接手动构建一个)
	// 记得把阈值设为：物流至少 1，质检至少 1，总数至少 3
	contract := InitContract("eCNY-MSIG-SC-20260429-0001", 10000000, 1, 1, 3)

	fmt.Println("==== 浙江农业集团：千万级化肥采购案多签启动 ====")

	// 2. 模拟真实业务流转
	err := contract.Deposit(10000000)
	if err != nil {
		fmt.Println(err)
		return // 钱没到位，直接终止
	}

	err1 := contract.Sign("A", RoleLogistics) // 顺丰快递小哥 A 确认发货
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println("✅ 物流节点 A 签名成功！")
	}

	err11 := contract.Sign("A", RoleLogistics) // 顺丰快递小哥 A 重复确认发货
	if err11 != nil {
		fmt.Println(err11)
	} else {
		fmt.Println("✅ 物流节点 A 签名成功！")
	}

	err2 := contract.Sign("a", RoleQuality) // 质检员 a 确认合格
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("✅ 质检节点 a 签名成功！")
	}

	err3 := contract.Sign("B", RoleLogistics) // 菜鸟驿站站长 B 再次确认（满足总数3的最后一块拼图）
	if err3 != nil {
		fmt.Println(err3)
	} else {
		fmt.Println("✅ 物流节点 B 签名成功！")
	}

}
