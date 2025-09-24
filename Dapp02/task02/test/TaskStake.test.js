const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("TaskStake", function () {
  let TaskStake, TaskStakeV2, taskStake, taskToken;
  let owner, alice, bob, upgrader, pauser;

  beforeEach(async function () {
    [owner, alice, bob, upgrader, pauser] = await ethers.getSigners();

    // 部署一个奖励 token (ERC20)
    const ERC20Mock = await ethers.getContractFactory("MyToken");
    taskToken = await ERC20Mock.deploy();
    await taskToken.mint(owner.address, ethers.utils.parseEther("1000000"));

    // 部署 TaskStake (UUPS 代理)
    TaskStake = await ethers.getContractFactory("TaskStake");
    taskStake = await upgrades.deployProxy(
      TaskStake,
      [taskToken.address, 1, 1000000, ethers.utils.parseEther("1")],
      { kind: "uups" }
    );

    // 授权角色
    const UPGRADER_ROLE = await taskStake.UPGRADER_ROLE();
    const PAUSER_ROLE = await taskStake.PAUSER_ROLE();

    await taskStake.grantRole(UPGRADER_ROLE, upgrader.address);
    await taskStake.grantRole(PAUSER_ROLE, pauser.address);

    // 添加一个 ERC20 池
    await taskStake.addPool(taskToken.address, 100, ethers.utils.parseEther("1"), 10);

    // 给 Alice 一些 Token
    await taskToken.mint(alice.address, ethers.utils.parseEther("1000"));
    await taskToken.connect(alice).approve(taskStake.address, ethers.utils.parseEther("1000"));
  });

  it("部署后检查初始参数", async function () {
    expect(await taskStake.taskToken()).to.equal(taskToken.address);
    expect(await taskStake.hasRole(await taskStake.DEFAULT_ADMIN_ROLE(), owner.address)).to.be.true;
  });

  it("Alice 可以质押、领奖、解除质押", async function () {
    // Alice deposit
    await expect(taskStake.connect(alice).deposit(1, ethers.utils.parseEther("10")))
      .to.emit(taskStake, "Deposit");

    // 快进几个区块
    for (let i = 0; i < 20; i++) {
      await ethers.provider.send("evm_mine");
    }

    // 领取奖励
    await expect(taskStake.connect(alice).claim(1)).to.emit(taskStake, "Claim");

    // 请求解除质押
    await expect(taskStake.connect(alice).unstake(1, ethers.utils.parseEther("5")))
      .to.emit(taskStake, "UnstakeRequested");
  });

  it("PAUSER_ROLE 可以暂停/恢复操作", async function () {
    await taskStake.connect(pauser).pauseStaking();

    await expect(
      taskStake.connect(alice).deposit(1, ethers.utils.parseEther("1"))
    ).to.be.revertedWith("TaskStake: staking is paused");

    await taskStake.connect(pauser).unpauseStaking();
    await expect(taskStake.connect(alice).deposit(1, ethers.utils.parseEther("1")))
      .to.emit(taskStake, "Deposit");
  });

  it("UPGRADER_ROLE 可以升级合约", async function () {
    TaskStakeV2 = await ethers.getContractFactory("TaskStakeV2");
    await expect(
      upgrades.upgradeProxy(taskStake.address, TaskStakeV2.connect(upgrader))
    ).to.not.be.reverted;

    const upgraded = await ethers.getContractAt("TaskStakeV2", taskStake.address);
    expect(await upgraded.version()).to.equal("V2");
  });

  it("非授权账户不能升级或暂停", async function () {
    const UPGRADER_ROLE = await taskStake.UPGRADER_ROLE();
    const PAUSER_ROLE = await taskStake.PAUSER_ROLE();

    expect(await taskStake.hasRole(UPGRADER_ROLE, alice.address)).to.be.false;
    expect(await taskStake.hasRole(PAUSER_ROLE, bob.address)).to.be.false;

    await expect(
      upgrades.upgradeProxy(taskStake.address, TaskStakeV2.connect(alice))
    ).to.be.revertedWithCustomError;

    await expect(taskStake.connect(bob).pauseStaking()).to.be.reverted;
  });
});
