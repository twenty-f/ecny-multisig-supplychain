package main

import (
	"testing"
)

// 测试函数必须以大写的 Test 开头，并传入 *testing.T 指针
func TestDeposit(t *testing.T) {
	// 1. 靶机准备：初始化一个需要 1000 万的合约
	contract := InitContract("TEST-MSIG-001", 10000000, 1, 1, 3)

	// 2. 第一轮攻击：故意只打 500 万，测试系统是否拦截
	err := contract.Deposit(5000000)
	if err == nil {
		// t.Errorf 是测试专属的报错，一旦执行，这个测试就被判定为 FAIL
		t.Errorf("🚨 致命漏洞：汇款金额不足，系统居然没有拦截！")
	} else {
		t.Logf("✅ 防御成功，成功拦截非法汇款：%v", err)
	}

	// 3. 第二轮正常操作：打入足额 1000 万
	err = contract.Deposit(10000000)
	if err != nil {
		t.Errorf("🚨 致命漏洞：足额汇款被意外拦截，错误信息：%v", err)
	}

	// 4. 终极断言：检查状态机的齿轮是否正确转动
	if contract.Status != StatusSigning {
		t.Errorf("🚨 状态机异常：足额汇款后，状态应为 %s，但当前为 %s", StatusSigning, contract.Status)
	} else {
		t.Log("✅ 状态机流转正常，已进入待签名状态！")
	}
}

func TestSignFlow(t *testing.T) {
	// 1. 初始化一個 1+1+3 的合約並完成匯款
	contract := InitContract("TEST-FLOW-001", 1000, 1, 1, 3)
	contract.Deposit(1000)

	// 2. 攻擊測試：重複簽名攔截
	contract.Sign("A", RoleLogistics)
	err := contract.Sign("A", RoleLogistics) // A 又簽了一次
	if err == nil {
		t.Errorf("🚨 漏洞：重複簽名未被攔截！")
	}

	// 3. 攻擊測試：非法角色攔截
	err = contract.Sign("Stranger", RoleLogistics)
	if err == nil {
		t.Errorf("🚨 漏洞：非白名單人員居然簽名成功了！")
	}

	// 4. 正常路徑：補齊剩餘簽名
	contract.Sign("a", RoleQuality) // 質檢 a 簽名 (此時 1+1，總數 2)

	if contract.Status == StatusCompleted {
		t.Errorf("🚨 邏輯錯誤：總數才 2 個，怎麼就打款了？")
	}

	contract.Sign("B", RoleLogistics) // 物流 B 簽名 (此時 1+1，總數 3，觸發！)

	// 5. 終極斷言：檢查是否自動打款
	if contract.Status != StatusCompleted {
		t.Errorf("🚨 邏輯錯誤：多簽已達標，但合約狀態仍為 %s", contract.Status)
	} else {
		t.Log("🎉 完美！自動化測試證實：1+1+3 多簽流程完美閉環，資金已安全撥付。")
	}
}
