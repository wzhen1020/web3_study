const { ethers, upgrades } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();

  console.log("Deploying with account:", deployer.address);

  const ERC20 = await ethers.getContractFactory("MyToken");
  const taskToken = await ERC20.deploy();
  await taskToken.mint(deployer.address, ethers.utils.parseEther("1000000"));

  const TaskStake = await ethers.getContractFactory("TaskStake");
  const taskStake = await upgrades.deployProxy(
    TaskStake,
    [taskToken.address, 1, 1000000, ethers.utils.parseEther("1")],
    { kind: "uups" }
  );

  console.log("TaskStake deployed to:", taskStake.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
