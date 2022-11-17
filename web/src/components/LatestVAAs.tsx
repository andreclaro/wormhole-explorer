import { ChainId, tryHexToNativeString } from "@certusone/wormhole-sdk";
import { parseVaa } from "@certusone/wormhole-sdk/lib/cjs/vaa";
import { parseTransferPayload } from "@certusone/wormhole-sdk/lib/cjs/utils";
import { ChevronRight } from "@mui/icons-material";
import { Card, IconButton, Typography } from "@mui/material";
import { Box } from "@mui/system";
import {
  createColumnHelper,
  getCoreRowModel,
  getExpandedRowModel,
  Row,
  useReactTable,
} from "@tanstack/react-table";
import { BigNumber } from "ethers";
import { ReactElement } from "react";
import useLatestNonPythNetVAAs from "../hooks/useLatestNonPythNetVAAs";
import { VAAsResponse } from "../hooks/useLatestVAAs";
import Table from "./Table";

const columnHelper = createColumnHelper<VAAsResponse>();

const columns = [
  columnHelper.display({
    id: "_expand",
    cell: ({ row }) =>
      row.getCanExpand() ? (
        <IconButton
          size="small"
          {...{
            onClick: row.getToggleExpandedHandler(),
            style: { cursor: "pointer" },
          }}
        >
          <ChevronRight
            sx={{
              transition: ".2s",
              transform: row.getIsExpanded() ? "rotate(90deg)" : undefined,
            }}
          />
        </IconButton>
      ) : null,
  }),
  columnHelper.accessor("_id", {
    id: "chain",
    header: () => "Chain",
    cell: (info) => info.getValue().split("/")[0],
  }),
  columnHelper.accessor("_id", {
    id: "emitter",
    header: () => "Emitter",
    cell: (info) => info.getValue().split("/")[1],
  }),
  columnHelper.accessor("_id", {
    id: "sequence",
    header: () => "Sequence",
    cell: (info) => info.getValue().split("/")[2],
  }),
  columnHelper.accessor("updatedAt", {
    header: () => "Observed At",
    cell: (info) => new Date(info.getValue()).toLocaleString(),
  }),
];

function VAADetails({ row }: { row: Row<VAAsResponse> }): ReactElement {
  const parsedVaa = parseVaa(
    new Uint8Array(Buffer.from(row.original.vaas, "base64"))
  );
  const payload = parsedVaa.payload;
  const parsedPayload = parseTransferPayload(payload);
  let token = parsedPayload.originAddress;
  // FromChain is a misnomer - actually OriginChain
  if (parsedPayload.originAddress && parsedPayload.originChain)
    try {
      token = tryHexToNativeString(
        parsedPayload.originAddress,
        parsedPayload.originChain as ChainId
      );
    } catch (e) {}
  return (
    <>
      Version: {parsedVaa.version}
      <br />
      Timestamp: {new Date(parsedVaa.timestamp * 1000).toLocaleString()}
      <br />
      Consistency: {parsedVaa.consistencyLevel}
      <br />
      Nonce: {parsedVaa.nonce}
      <br />
      Origin: {parsedPayload.originChain}
      <br />
      Token: {token}
      <br />
      Amount: {BigNumber.from(parsedPayload.amount).toString()}
      <br />
    </>
  );
}

function LatestVAAs() {
  const vaas = useLatestNonPythNetVAAs();
  const table = useReactTable({
    columns,
    data: vaas,
    getRowId: (vaa) => vaa._id,
    getRowCanExpand: () => true,
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    enableSorting: false,
  });
  return (
    <Box m={2}>
      <Card>
        <Box m={2}>
          <Typography variant="h5">Latest Messages</Typography>
        </Box>
        <Table<VAAsResponse> table={table} renderSubComponent={VAADetails} />
      </Card>
    </Box>
  );
}
export default LatestVAAs;
