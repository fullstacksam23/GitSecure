import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useMemo, useState } from "react";
import { api } from "../api/client";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";

function parseGithubRepo(value) {
  const trimmed = String(value || "").trim();
  if (!trimmed) return { owner: "", repo: "", valid: false };

  const normalized = trimmed
    .replace(/^https?:\/\/(www\.)?github\.com\//i, "")
    .replace(/^github\.com\//i, "")
    .replace(/\/+$/, "");

  const [owner = "", repo = ""] = normalized.split("/");
  const cleanRepo = repo.replace(/\.git$/i, "");

  return {
    owner,
    repo: cleanRepo,
    valid: Boolean(owner && cleanRepo),
  };
}

export default function NewScanPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [githubUrl, setGithubUrl] = useState("");
  const [owner, setOwner] = useState("");
  const [repo, setRepo] = useState("");
  const parsedUrlRepo = useMemo(() => parseGithubRepo(githubUrl), [githubUrl]);
  const manualOwner = owner.trim();
  const manualRepo = repo.trim();
  const hasManualRepo = Boolean(manualOwner && manualRepo);
  const hasUrlRepo = parsedUrlRepo.valid;
  const canSubmit = hasUrlRepo || hasManualRepo;

  const mutation = useMutation({
    mutationFn: api.startScan,
    onSuccess: (job) => {
      queryClient.invalidateQueries({ queryKey: ["scans"] });
      queryClient.invalidateQueries({ queryKey: ["sidebar-scans"] });
      queryClient.invalidateQueries({ queryKey: ["history"] });
      queryClient.invalidateQueries({ queryKey: ["compare-scans-list"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard-summary"] });
      toast.success("Scan queued");
      navigate(`/scans/${job.job_id}`);
    },
  });

  function handleStartScan() {
    if (mutation.isPending || !canSubmit) return;

    const payload = hasManualRepo
      ? { owner: manualOwner, repo: manualRepo }
      : { owner: parsedUrlRepo.owner, repo: parsedUrlRepo.repo };

    mutation.mutate(payload);
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">New Scan</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Launch a repository scan</h1>
        <p className="mt-2 max-w-2xl text-sm leading-6 text-slate-400">
          Start a scan with either a full GitHub repository URL or the owner and repository name.
        </p>
      </div>

      <div className="panel p-6">
        <div className="grid gap-4">
          <div className="surface p-4">
            <p className="text-sm font-medium text-white">GitHub repository URL</p>
            <p className="mt-1 text-sm text-slate-400">Paste a full link like `https://github.com/owner/repo`.</p>
            <Input
              className="mt-3"
              placeholder="https://github.com/vercel/next.js"
              value={githubUrl}
              onChange={(event) => setGithubUrl(event.target.value)}
            />
          </div>

          <div className="flex items-center gap-3 py-1">
            <div className="h-px flex-1 bg-white/10" />
            <span className="text-xs uppercase tracking-[0.22em] text-slate-500">Or enter it manually</span>
            <div className="h-px flex-1 bg-white/10" />
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <div>
              <p className="mb-2 text-sm font-medium text-white">GitHub owner</p>
              <Input placeholder="vercel" value={owner} onChange={(event) => setOwner(event.target.value)} />
            </div>
            <div>
              <p className="mb-2 text-sm font-medium text-white">Repository name</p>
              <Input placeholder="next.js" value={repo} onChange={(event) => setRepo(event.target.value)} />
            </div>
          </div>

          <div className="rounded-2xl border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-slate-400">
            {hasManualRepo
              ? `Scanning ${manualOwner}/${manualRepo} using the owner and repository fields.`
              : hasUrlRepo
                ? `Scanning ${parsedUrlRepo.owner}/${parsedUrlRepo.repo} from the GitHub URL.`
                : "Enter either a GitHub URL or both the owner and repository name to enable scanning."}
          </div>

          <Button onClick={handleStartScan} disabled={mutation.isPending || !canSubmit}>
            {mutation.isPending ? "Starting scan..." : "Start Scan"}
          </Button>
          {mutation.isError ? <p className="text-sm text-red-300">{mutation.error.message}</p> : null}
        </div>
      </div>
    </div>
  );
}
