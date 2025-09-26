## 优化策略

### 优化策略 1：减少 require 检查（如果安全可控）

	> src/MathOptimized1.sol

```sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract MathOptimized1 {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        return a + b;
    }

    // 使用 unchecked 避免 Solidity 自动溢出检查
    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked {
            return a - b;
        }
    }
}

```

`unchecked` 可以省掉 Solidity 的安全溢出检查，减少 Gas。

适用于你确信 `b <= a` 的情况。

### 优化策略 2：内联函数或直接返回表达式

> `src/Math.sol`：

```sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract MathOptimized2 {
    function add(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked { return a + b; }
    }

    function sub(uint256 a, uint256 b) external pure returns (uint256) {
        unchecked { return a - b; }
    }
}

```

减少局部变量，直接返回表达式。

同时使用 `unchecked` 进一步优化。

## 测试

更新测试合约：

> `test/Math.t.sol`：

```sol
 math = new MathOptimized1();
// 或 math = new MathOptimized2();

```

执行测试命令

```lua
forge test --match-path test/Math.t.sol
```



## 优化效果分析



1. **优化前**

   ```lua
   [PASS] testAdd() (gas: 8994)
   [PASS] testSub() (gas: 9041)
   ```

   

2. **优化1（unchecked）**

   - 去掉溢出检查，`sub` Gas 降幅明显。
   - `add` Gas 减少不明显，因为加法原本没有 require。

   ```lua
   [PASS] testAdd() (gas: 8994)
   [PASS] testSub() (gas: 8836)
   ```

   

3. **优化2（inline + unchecked）**

   - 局部变量减少 + unchecked，`add` 和 `sub` 都节省了 Gas。
   - 对于大批量计算或者循环调用，可显著降低交易成本。

   ```lua
   [PASS] testAdd() (gas: 8815)
   [PASS] testSub() (gas: 8836)
   ```

   

| 版本  | Add Gas | Sub Gas |
| ----- | ------- | ------- |
| 原始  | 8994    | 9041    |
| 优化1 | 8994    | 8836    |
| 优化2 | 8815    | 8836    |

