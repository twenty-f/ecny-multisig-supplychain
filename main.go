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
func InitContract(contractID string, amount int) *MultiSigContract {
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
		LogisticsMin:   1,
		QualityMin:     1,
		TotalMin:       3,
		Status:         StatusPending,
	}

	return m
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
		return nil
	}
	m.QualityCount++
	return nil
}

func main() {
	contract := InitContract("SCF-MULTISIG-0001", 100000)

	fmt.Printf("初始化完成: ContractID=%s, Amount=%d, Status=%s\n", contract.ContractID, contract.Amount, contract.Status)
	fmt.Printf("阈值配置: LogisticsMin=%d, QualityMin=%d, TotalMin=%d\n", contract.LogisticsMin, contract.QualityMin, contract.TotalMin)

	// 模拟几次签名调用（仅骨架演示，不包含核心逻辑）
	_ = contract.Sign("A", RoleLogistics)
	_ = contract.Sign("a", RoleQuality)
	_ = contract.Sign("B", RoleLogistics)

	fmt.Println("---- 合约状态快照 ----")
	fmt.Printf("Status=%s\n", contract.Status)
	fmt.Printf("Count: logistics=%d, quality=%d, total=%d\n", contract.LogisticsCount, contract.QualityCount, contract.TotalCount)
	fmt.Printf("Signatures(logistics)=%v\n", contract.Signatures[RoleLogistics])
	fmt.Printf("Signatures(quality)=%v\n", contract.Signatures[RoleQuality])
}
