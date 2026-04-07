import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { api } from "../api/client";
import { Button } from "../components/ui/button";
import { Input } from "../components/ui/input";

export default function NewScanPage() {
  const navigate = useNavigate();
  const [owner, setOwner] = useState("");
  const [repo, setRepo] = useState("");

  const mutation = useMutation({
    mutationFn: api.startScan,
    onSuccess: (job) => {
      toast.success("Scan queued");
      navigate(`/scans/${job.job_id}`);
    },
  });

  return (
    <div className="mx-auto max-w-2xl space-y-6 pb-8">
      <div>
        <p className="text-sm uppercase tracking-[0.28em] text-cyan-300">New Scan</p>
        <h1 className="mt-3 text-3xl font-semibold text-white">Launch a repository scan</h1>
      </div>

      <div className="panel p-6">
        <div className="grid gap-4">
          <Input placeholder="GitHub owner" value={owner} onChange={(event) => setOwner(event.target.value)} />
          <Input placeholder="Repository name" value={repo} onChange={(event) => setRepo(event.target.value)} />
          <Button onClick={() => mutation.mutate({ owner, repo })} disabled={mutation.isPending || !owner || !repo}>
            {mutation.isPending ? "Starting scan..." : "Start Scan"}
          </Button>
          {mutation.isError ? <p className="text-sm text-red-300">{mutation.error.message}</p> : null}
        </div>
      </div>
    </div>
  );
}
