export type EvmBlock = {
  number: bigint;
  hash: string;
  timestamp: number; // epoch seconds
  transactions?: EvmTransaction[];
};

export type EvmLog = {
  blockTime?: number; // epoch seconds
  blockNumber: bigint;
  blockHash: string;
  address: string;
  removed: boolean;
  data: string;
  transactionHash: string;
  transactionIndex: string;
  topics: string[];
  logIndex: number;
  chainId: number;
  chain: string;
};

export type EvmTransaction = {
  blockHash: string;
  blockNumber: bigint;
  chainId: number;
  from: string;
  gas: string;
  gasPrice: string;
  hash: string;
  input: string;
  maxFeePerGas: string;
  maxPriorityFeePerGas: string;
  nonce: string;
  r: string;
  s: string;
  status?: string;
  to: string;
  transactionIndex: string;
  type: string;
  v: string;
  value: string;
  timestamp: number;
  environment: string;
  chain: string;
  logs: EvmTransactionLog[];
};

export type EvmTransactionLog = { address: string; topics: string[]; data: string };

export type EvmTag = "finalized" | "latest" | "safe";

export type EvmLogFilter = {
  fromBlock: bigint | EvmTag;
  toBlock: bigint | EvmTag;
  addresses: string[];
  topics: string[];
};

export type ReceiptTransaction = {
  status: string;
  transactionHash: string;
  logs: EvmTransactionLog[];
};
