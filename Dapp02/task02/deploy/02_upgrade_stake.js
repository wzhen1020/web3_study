const { ethers, upgrades } = require("hardhat");

async function main() {
  const proxyAddress = "YOUR_DEPLOYED_PROXY_ADDRESS"; // 填写实际地址

  const TaskStakeV2 = await ethers.getContractFactory("TaskStakeV2");
  const upgraded = await upgrades.upgradeProxy(proxyAddress, TaskStakeV2);

  console.log("TaskStake upgraded to V2 at:", upgraded.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
