import { afterEach, describe, it, expect, jest } from "@jest/globals";
import { PollEvmLogsMetadata, PollEvm, PollEvmLogsConfig } from "../../../../src/domain/actions";
import {
  EvmBlockRepository,
  MetadataRepository,
  StatRepository,
} from "../../../../src/domain/repositories";
import { EvmBlock, EvmLog, ReceiptTransaction } from "../../../../src/domain/entities";
import { thenWaitForAssertion } from "../../../wait-assertion";

let cfg = PollEvmLogsConfig.fromBlock("acala", 0n);

let getBlocksSpy: jest.SpiedFunction<EvmBlockRepository["getBlocks"]>;
let getLogsSpy: jest.SpiedFunction<EvmBlockRepository["getFilteredLogs"]>;
let handlerSpy: jest.SpiedFunction<(logs: EvmLog[]) => Promise<void>>;
let metadataSaveSpy: jest.SpiedFunction<MetadataRepository<PollEvmLogsMetadata>["save"]>;

let metadataRepo: MetadataRepository<PollEvmLogsMetadata>;
let evmBlockRepo: EvmBlockRepository;
let statsRepo: StatRepository;

let handlers = {
  working: (logs: EvmLog[]) => Promise.resolve(),
  failing: (logs: EvmLog[]) => Promise.reject(),
};
let pollEvm: PollEvm;

describe("PollEvm", () => {
  afterEach(async () => {
    await pollEvm.stop();
  });

  it("should be able to read logs from latest block when no fromBlock is configured", async () => {
    const currentHeight = 10n;
    const blocksAhead = 1n;
    givenEvmBlockRepository(currentHeight, blocksAhead);
    givenMetadataRepository();
    givenStatsRepository();
    givenPollEvmLogs();

    await whenPollEvmLogsStarts();

    await thenWaitForAssertion(
      () => expect(getBlocksSpy).toHaveReturnedTimes(1),
      () =>
        expect(getBlocksSpy).toHaveBeenCalledWith(
          "acala",
          new Set([currentHeight, currentHeight + 1n]),
          false
        ),
      () =>
        expect(getLogsSpy).toBeCalledWith("acala", {
          addresses: cfg.filters[0].addresses,
          topics: cfg.filters[0].topics,
          fromBlock: currentHeight + blocksAhead,
          toBlock: currentHeight + blocksAhead,
        })
    );
  });

  it("should be able to read logs from last known block when configured from is before", async () => {
    const lastExtractedBlock = 10n;
    const blocksAhead = 10n;
    givenEvmBlockRepository(lastExtractedBlock, blocksAhead);
    givenMetadataRepository({ lastBlock: lastExtractedBlock });
    givenStatsRepository();
    givenPollEvmLogs(lastExtractedBlock - 10n);

    await whenPollEvmLogsStarts();

    await thenWaitForAssertion(
      () => () =>
        expect(getBlocksSpy).toHaveBeenCalledWith(
          new Set([lastExtractedBlock, lastExtractedBlock + 1n])
        ),
      () =>
        expect(getLogsSpy).toBeCalledWith("acala", {
          addresses: cfg.filters[0].addresses,
          topics: cfg.filters[0].topics,
          fromBlock: lastExtractedBlock + 1n,
          toBlock: lastExtractedBlock + blocksAhead,
        })
    );
  });

  it("should pass logs to handlers and persist metadata", async () => {
    const currentHeight = 10n;
    const blocksAhead = 1n;
    givenEvmBlockRepository(currentHeight, blocksAhead);
    givenMetadataRepository();
    givenStatsRepository();
    givenPollEvmLogs(currentHeight);

    await whenPollEvmLogsStarts();

    await thenWaitForAssertion(
      () => expect(handlerSpy).toHaveBeenCalledWith(expect.any(Array)),
      () =>
        expect(metadataSaveSpy).toBeCalledWith("watch-evm-logs", {
          lastBlock: currentHeight + blocksAhead,
        })
    );
  });
});

const givenEvmBlockRepository = (height?: bigint, blocksAhead?: bigint) => {
  const logsResponse: EvmLog[] = [];
  const blocksResponse: Record<string, EvmBlock> = {};
  const receiptResponse: Record<string, ReceiptTransaction> = {};
  if (height) {
    for (let index = 0n; index <= (blocksAhead ?? 1n); index++) {
      logsResponse.push({
        blockNumber: height + index,
        blockHash: `0x0${index}`,
        blockTime: 0,
        address: "",
        removed: false,
        data: "",
        transactionHash: "",
        transactionIndex: "",
        topics: [],
        logIndex: 0,
        chainId: 2,
        chain: "ethereum",
      });
      blocksResponse[`0x0${index}`] = {
        timestamp: 0,
        hash: `0x0${index}`,
        number: height + index,
      };
      receiptResponse[`0x0${index}`] = {
        status: "0x1",
        transactionHash: `0x0${index}`,
        logs: [
          {
            address: "0xf890982f9310df57d00f659cf4fd87e65aded8d7",
            topics: ["0xbccc00b713f54173962e7de6098f643d8ebf53d488d71f4b2a5171496d038f9e"],
            data: "0x",
          },
        ],
      };
    }
  }

  evmBlockRepo = {
    getBlocks: () => Promise.resolve(blocksResponse),
    getBlockHeight: () => Promise.resolve(height ? height + (blocksAhead ?? 10n) : 10n),
    getFilteredLogs: () => Promise.resolve(logsResponse),
    getTransactionReceipt: () => Promise.resolve(receiptResponse),
    getBlock: () => Promise.resolve(blocksResponse[0]),
  };

  getBlocksSpy = jest.spyOn(evmBlockRepo, "getBlocks");
  getLogsSpy = jest.spyOn(evmBlockRepo, "getFilteredLogs");
  handlerSpy = jest.spyOn(handlers, "working");
};

const givenMetadataRepository = (data?: PollEvmLogsMetadata) => {
  metadataRepo = {
    get: () => Promise.resolve(data),
    save: () => Promise.resolve(),
  };
  metadataSaveSpy = jest.spyOn(metadataRepo, "save");
};

const givenStatsRepository = () => {
  statsRepo = {
    count: () => {},
    measure: () => {},
    report: () => Promise.resolve(""),
  };
};

const givenPollEvmLogs = (from?: bigint) => {
  cfg.setFromBlock(from);
  pollEvm = new PollEvm(evmBlockRepo, metadataRepo, statsRepo, cfg, "GetEvmLogs");
};

const whenPollEvmLogsStarts = async () => {
  pollEvm.run([handlers.working]);
};
