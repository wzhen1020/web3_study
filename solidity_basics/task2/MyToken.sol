// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is Ownable {
    // 代币元数据
    string private _name;
    string private _symbol;
    uint8 private _decimals;

    // 代币总供应量
    uint256 private _totalSupply;

    // 余额映射
    mapping(address => uint256) private _balances;

    // 授权映射
    mapping(address => mapping(address => uint256)) private _allowances;

    // 事件定义
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);


    // 构造函数
    constructor(uint256 initialSupply) Ownable(msg.sender) {
        _name = "MyToken";
        _symbol = "MTK";
        _decimals = 18;
        
        _mint(msg.sender, initialSupply * 10 ** _decimals);
    }

    // 获取代币名称
    function name() public view returns (string memory) {
        return _name;
    }
    
    // 获取代币符号
    function symbol() public view returns (string memory) {
        return _symbol;
    }
    
    // 获取小数位数
    function decimals() public view returns (uint8) {
        return _decimals;
    }
    
    // 获取总供应量
    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }
    
    // 获取账户余额
    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }
    
    // 转账功能
    function transfer(address to, uint256 value) public returns (bool) {
        address owner = msg.sender;
        _transfer(owner, to, value);
        return true;
    }
    
    // 授权功能
    function approve(address spender, uint256 value) public returns (bool) {
        address owner = msg.sender;
        _approve(owner, spender, value);
        return true;
    }
    
    // 查询授权额度
    function allowance(address owner, address spender) public view returns (uint256) {
        return _allowances[owner][spender];
    }
    
    // 代扣转账功能
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        address spender = msg.sender;
        // 转出地址 余额验证
        _spendAllowance(spender,from, value);
        _transfer(from, to, value);
        return true;
    }
    
    // 增发代币（仅所有者）
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }
    
    // 内部转账函数
    function _transfer(address from, address to, uint256 value) internal {
        require(from != address(0), "ERC20: transfer from the zero address");
        require(to != address(0), "ERC20: transfer to the zero address");
        
        uint256 fromBalance = _balances[from];
        require(fromBalance >= value, "ERC20: transfer amount exceeds balance");
        
        unchecked {
            _balances[from] = fromBalance - value;
        }
        _balances[to] += value;
        
        emit Transfer(from, to, value);
    }

    // 内部铸造函数
    function _mint(address account, uint256 value) internal {
        require(account != address(0), "ERC20: mint to the zero address");
        
        _totalSupply += value;
        _balances[account] += value;
        
        emit Transfer(address(0), account, value);
    }

    // 内部授权函数
    function _approve(address owner, address spender, uint256 value) internal {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");
        // 授权
        _allowances[owner][spender] = value;
        emit Approval(owner, spender, value);
    }

    // 内部消耗授权额度函数
    function _spendAllowance(address owner, address spender, uint256 value) internal {
        // 查询转出方授权额度
        uint256 currentAllowance = allowance(owner, spender);
        if (currentAllowance != type(uint256).max) {
            // 剩余额度是否大于转帐额度
            require(currentAllowance >= value, "ERC20: insufficient allowance");
            unchecked {
                _approve(owner, spender, currentAllowance - value);
            }
        }
    }
}




