import { useNavigate } from "react-router-dom";
import { getBatchId, getBatchLanguage, getBatchProgress } from "../../lib/ecosystem";
import { formatDate } from "../../lib/utils";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";
import ProgressBar from "./ProgressBar";
import StatusBadge from "./StatusBadge";

export default function BatchTable({ items = [] }) {
  const navigate = useNavigate();

  return (
    <div className="panel overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow className="border-white/10 hover:bg-transparent">
            <TableHead>Batch</TableHead>
            <TableHead>Language</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Progress</TableHead>
            <TableHead>Created</TableHead>
            <TableHead>Completed</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {items.map((item) => {
            const progress = getBatchProgress(item);
            return (
              <TableRow
                key={getBatchId(item)}
                className="cursor-pointer border-white/10 hover:bg-white/[0.04]"
                onClick={() => navigate(`/ecosystem/batches/${getBatchId(item)}`)}
              >
                <TableCell>
                  <div>
                    <p className="font-medium text-white">{getBatchId(item)}</p>
                    <p className="mt-1 text-xs text-slate-500">Per-batch aggregation scope</p>
                  </div>
                </TableCell>
                <TableCell className="text-slate-200">{getBatchLanguage(item)}</TableCell>
                <TableCell>
                  <StatusBadge status={item.status} />
                </TableCell>
                <TableCell className="min-w-[220px]">
                  <ProgressBar
                    value={progress.percent}
                    label={
                      <>
                        <span>{progress.completed} scanned</span>
                        <span>{progress.total} total</span>
                      </>
                    }
                    tone={item.status === "completed" ? "success" : item.status === "running" ? "default" : "warning"}
                  />
                </TableCell>
                <TableCell className="text-slate-300">{formatDate(item.created_at)}</TableCell>
                <TableCell className="text-slate-300">{formatDate(item.completed_at)}</TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
