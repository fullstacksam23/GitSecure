import { useEffect, useMemo, useState } from "react";
import { CheckCircle2, Clock3, LoaderCircle, ShieldX } from "lucide-react";
import ScanForm from "../components/ScanForm";
import PageHeader from "../components/PageHeader";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "../components/ui/dialog";
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card";
import { Button } from "../components/ui/button";
import { api } from "../services/api";

const statusConfig = {
  queued: {
    icon: Clock3,
    title: "Queued for scanning",
    description: "The job has been accepted and is waiting for the worker to pick it up.",
  },
  running: {
    icon: LoaderCircle,
    title: "Scan in progress",
    description: "SBOM extraction and vulnerability analysis are currently running.",
  },
  completed: {
    icon: CheckCircle2,
    title: "Scan complete",
    description: "Results are ready to review in the vulnerability explorer.",
  },
  failed: {
    icon: ShieldX,
    title: "Scan failed",
    description: "The backend reported a failure and the job needs attention.",
  },
};

export default function ScanPage() {
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [job, setJob] = useState(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  const status = job?.status || "queued";
  const details = statusConfig[status] || statusConfig.queued;
  const StatusIcon = details.icon;

  async function handleScan(parsed) {
    setSubmitting(true);
    setError("");

    try {
      const createdJob = await api.startScan(parsed);
      setJob(createdJob);
      setDialogOpen(true);
    } catch (err) {
      setError(err.message || "Unable to start the scan.");
    } finally {
      setSubmitting(false);
    }
  }

  useEffect(() => {
    if (!job?.job_id || ["completed", "failed"].includes(job.status)) return undefined;

    const interval = window.setInterval(async () => {
      try {
        const nextJob = await api.getScan(job.job_id);
        setJob(nextJob);
      } catch {
        window.clearInterval(interval);
      }
    }, 4000);

    return () => window.clearInterval(interval);
  }, [job]);

  const statusSteps = useMemo(
    () => [
      { label: "Queued", active: ["queued", "running", "completed"].includes(status) },
      { label: "Scanning", active: ["running", "completed"].includes(status) },
      { label: "Ready", active: status === "completed" },
    ],
    [status]
  );

  return (
    <div className="space-y-6">
      <PageHeader
        eyebrow="Start Scan"
        title="Launch a new vulnerability scan"
        description="Trigger the backend pipeline for a public GitHub repository and track the job from queue to completion."
      />

      <div className="grid gap-6 xl:grid-cols-[1.3fr_0.9fr]">
        <ScanForm onSubmit={handleScan} loading={submitting} />

        <Card>
          <CardHeader>
            <CardTitle>Scan status</CardTitle>
          </CardHeader>
          <CardContent>
            {job ? (
              <div className="space-y-5">
                <div className="flex items-center gap-3 rounded-2xl bg-muted/70 p-4">
                  <div className="rounded-2xl bg-white p-3 shadow-sm">
                    <StatusIcon className={`h-5 w-5 ${status === "running" ? "animate-spin" : ""}`} />
                  </div>
                  <div>
                    <p className="font-semibold text-slate-950">{details.title}</p>
                    <p className="text-sm text-muted-foreground">{details.description}</p>
                  </div>
                </div>

                <div className="grid gap-3">
                  {statusSteps.map((step) => (
                    <div
                      key={step.label}
                      className={`rounded-2xl px-4 py-3 text-sm font-medium ${
                        step.active ? "bg-slate-950 text-white" : "bg-muted text-muted-foreground"
                      }`}
                    >
                      {step.label}
                    </div>
                  ))}
                </div>

                <div className="space-y-3 rounded-2xl border border-border/70 bg-white/70 p-4 text-sm">
                  <div className="flex items-center justify-between gap-4">
                    <span className="text-muted-foreground">Repository</span>
                    <span className="font-semibold text-slate-950">{job.repo}</span>
                  </div>
                  <div className="flex items-center justify-between gap-4">
                    <span className="text-muted-foreground">Job ID</span>
                    <span className="font-mono text-xs text-slate-950">{job.job_id}</span>
                  </div>
                  <div className="flex items-center justify-between gap-4">
                    <span className="text-muted-foreground">Status</span>
                    <span className="font-semibold uppercase tracking-[0.18em] text-slate-950">
                      {job.status}
                    </span>
                  </div>
                  {job.commit_hash ? (
                    <div className="flex items-center justify-between gap-4">
                      <span className="text-muted-foreground">Commit</span>
                      <span className="font-mono text-xs text-slate-950">{job.commit_hash}</span>
                    </div>
                  ) : null}
                </div>
              </div>
            ) : (
              <div className="rounded-3xl bg-gradient-to-br from-sky-50 to-amber-50 p-6 text-sm leading-6 text-slate-600">
                A started scan returns a `202 Accepted` response from `POST /scan`. The page then polls
                the job endpoint described in the frontend API contract so users can track progress.
              </div>
            )}

            {error ? (
              <div className="mt-4 rounded-2xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
                {error}
              </div>
            ) : null}
          </CardContent>
        </Card>
      </div>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Scan accepted</DialogTitle>
            <DialogDescription>
              The backend has queued the job. This dialog becomes more useful once `GET /scans/:jobId`
              is implemented and starts returning live status updates.
            </DialogDescription>
          </DialogHeader>

          {job ? (
            <div className="mt-6 space-y-4">
              <div className="rounded-2xl bg-slate-950 p-4 text-white">
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-sky-200">Job ID</p>
                <p className="mt-2 font-mono text-sm">{job.job_id}</p>
              </div>

              <Button className="w-full" onClick={() => setDialogOpen(false)}>
                Continue monitoring
              </Button>
            </div>
          ) : null}
        </DialogContent>
      </Dialog>
    </div>
  );
}
