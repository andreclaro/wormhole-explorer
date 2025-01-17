import { SnsConfig } from "../../src/infrastructure/repositories";
import { Config, ChainRPCConfig } from "../../src/infrastructure/config";

export const configMock = (): Config => {
  const chainsRecord: Record<string, ChainRPCConfig> = {
    solana: {
      name: "solana",
      network: "devnet",
      chainId: 1,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    bsc: {
      name: "bsc",
      network: "BNB Smart Chain testnet",
      chainId: 4,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    polygon: {
      name: "polygon",
      network: "mumbai",
      chainId: 5,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    avalanche: {
      name: "avalanche",
      network: "testnet",
      chainId: 6,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    oasis: {
      name: "oasis",
      network: "emerald",
      chainId: 7,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    fantom: {
      name: "fantom",
      network: "testnet",
      chainId: 10,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    karura: {
      name: "karura",
      network: "testnet",
      chainId: 11,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    acala: {
      name: "acala",
      network: "testnet",
      chainId: 12,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    klaytn: {
      name: "klaytn",
      network: "baobab",
      chainId: 13,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    celo: {
      name: "celo",
      network: "alfajores",
      chainId: 14,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    moonbeam: {
      name: "moonbeam",
      network: "testnet",
      chainId: 16,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    injective: {
      name: "injective",
      network: "testnet",
      chainId: 19,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    osmosis: {
      name: "osmosis",
      network: "testnet",
      chainId: 20,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    aptos: {
      name: "aptos",
      network: "testnet",
      chainId: 22,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    arbitrum: {
      name: "arbitrum",
      network: "goerli",
      chainId: 23,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    optimism: {
      name: "optimism",
      network: "goerli",
      chainId: 24,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    base: {
      name: "base",
      network: "goerli",
      chainId: 30,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    sei: {
      name: "sei",
      network: "mainnet",
      chainId: 32,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    scroll: {
      name: "scroll",
      network: "testnet",
      chainId: 34,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    mantle: {
      name: "mantle",
      network: "testnet",
      chainId: 35,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    blast: {
      name: "blast",
      network: "testnet",
      chainId: 36,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    evmos: {
      name: "evmos",
      network: "testnet",
      chainId: 4001,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    kujira: {
      name: "kujira",
      network: "testnet",
      chainId: 4002,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    xlayer: {
      name: "xlayer",
      network: "testnet",
      chainId: 37,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    "ethereum-sepolia": {
      name: "ethereum-sepolia",
      network: "sepolia",
      chainId: 10002,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    "arbitrum-sepolia": {
      name: "arbitrum-sepolia",
      network: "sepolia",
      chainId: 10003,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    "base-sepolia": {
      name: "base-sepolia",
      network: "sepolia",
      chainId: 10004,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    "optimism-sepolia": {
      name: "optimism-sepolia",
      network: "sepolia",
      chainId: 10005,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    sui: {
      name: "sui",
      network: "testnet",
      chainId: 21,
      rpcs: ["https://fullnode.testnet.sui.io:443"],
      timeout: 10000,
    },
    "ethereum-holesky": {
      name: "ethereum-holesky",
      network: "holesky",
      chainId: 10006,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    "polygon-sepolia": {
      name: "polygon-sepolia",
      network: "sepolia",
      chainId: 10007,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
    wormchain: {
      name: "wormchain",
      network: "testnet",
      chainId: 3104,
      rpcs: ["http://localhost"],
      timeout: 10000,
    },
  };

  const snsConfig: SnsConfig = {
    region: "us-east",
    topicArn: "123333223232s",
    subject: "",
    groupId: "1",
    credentials: {
      accessKeyId: "212312312323",
      secretAccessKey: "244122wdsd",
      url: "",
    },
  };

  const cfg: Config = {
    environment: "testnet",
    port: 999,
    logLevel: "info",
    dryRun: false,
    sns: snsConfig,
    metadata: {
      dir: "./metadata-repo/jobs",
    },
    jobs: {
      dir: "./metadata-repo/jobs",
    },
    chains: chainsRecord,
    enabledPlatforms: ["solana", "evm", "sui", "aptos", "wormchain", "sei"],
  };

  return cfg;
};
