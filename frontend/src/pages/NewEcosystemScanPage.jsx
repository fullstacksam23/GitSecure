import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Layers3, Radar } from "lucide-react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { api } from "../api/client";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";

const languageOptions = [
  { value: "go", label: "Go" },
  { value: "javascript", label: "JavaScript" },
  { value: "typescript", label: "TypeScript" },
  { value: "python", label: "Python" },
  { value: "rust", label: "Rust" },
  { value: "java", label: "Java" },
];

export default function NewEcosystemScanPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [language, setLanguage] = useState("go");
  const [repoCount, setRepoCount] = useState("20");

  const mutation = useMutation({
    mutationFn: api.startBatchScan,
    onSuccess: (batch) => {
      queryClient.invalidateQueries({ queryKey: ["ecosystem-batches"] });
      queryClient.invalidateQueries({ queryKey: ["ecosystem-batch", batch.batch_id] });
      queryClient.invalidateQueries({ queryKey: ["ecosystem-batch-summary", batch.batch_id] });
      queryClient.invalidateQueries({ queryKey: ["ecosystem-batch-repos", batch.batch_id] });
      toast.success("Ecosystem batch queued");
      navigate(`/ecosystem/batches/${batch.batch_id}`);
    },
  });

  const parsedRepoCount = Number(repoCount);
  const isValid = language && Number.isFinite(parsedRepoCount) && parsedRepoCount > 0;

  return (
    <div className="mx-auto max-w-3xl space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">New Ecosystem Scan</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Launch a batch-scoped ecosystem scan</h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-400">
          Pick a language and how many high-signal repositories to scan. The backend will create one batch and aggregate progress and metrics only within that batch.
        </p>
      </div>

      <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
        <div className="panel p-6">
          <div className="mb-6 flex items-center gap-3">
            <div className="rounded-2xl border border-cyan-400/20 bg-cyan-400/10 p-3 text-cyan-200">
              <Radar className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-white">Batch configuration</h2>
              <p className="text-sm text-slate-400">This request starts `GET /batch/scan` with query parameters.</p>
            </div>
          </div>

          <div className="grid gap-4">
            <label className="grid gap-2 text-sm text-slate-300">
              Language
              <select
                className="h-11 rounded-2xl border border-white/10 bg-white/[0.04] px-4 text-sm text-slate-200"
                value={language}
                onChange={(event) => setLanguage(event.target.value)}
              >
                {languageOptions.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </label>

            <label className="grid gap-2 text-sm text-slate-300">
              Number of repositories
              <Input
                type="number"
                min="1"
                max="100"
                step="1"
                value={repoCount}
                onChange={(event) => setRepoCount(event.target.value)}
                placeholder="20"
              />
            </label>

            <Button
              onClick={() => mutation.mutate({ language, repoCount: parsedRepoCount })}
              disabled={!isValid || mutation.isPending}
            >
              {mutation.isPending ? "Starting ecosystem scan..." : "Start ecosystem scan"}
            </Button>

            {mutation.isError ? <p className="text-sm text-red-300">{mutation.error.message}</p> : null}
          </div>
        </div>

        <div className="panel p-6">
          <div className="mb-5 flex items-center gap-3">
            <div className="rounded-2xl border border-white/10 bg-white/[0.05] p-3 text-slate-200">
              <Layers3 className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-white">What happens next</h2>
              <p className="text-sm text-slate-400">The frontend takes you directly into the batch detail page after queueing.</p>
            </div>
          </div>

          <div className="space-y-3 text-sm leading-6 text-slate-400">
            <div className="surface p-4">A new batch id is created for the selected language.</div>
            <div className="surface p-4">Top repositories are fetched and queued as batch scan jobs.</div>
            <div className="surface p-4">Progress and aggregate metrics update on the batch detail page.</div>
          </div>
        </div>
      </div>
    </div>
  );
}
