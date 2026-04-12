import { useQuery } from "@tanstack/react-query";
import { api } from "../api/client";
import { isActiveBatchStatus } from "../lib/ecosystem";

function withDefaultPagination(data, fallbackPage = 1, fallbackPageSize = 20) {
  if (!data) {
    return {
      items: [],
      pagination: {
        page: fallbackPage,
        page_size: fallbackPageSize,
        total_items: 0,
        total_pages: 0,
      },
    };
  }

  if (Array.isArray(data)) {
    return {
      items: data,
      pagination: {
        page: fallbackPage,
        page_size: data.length || fallbackPageSize,
        total_items: data.length,
        total_pages: 1,
      },
    };
  }

  return {
    ...data,
    items: data.items || [],
    pagination: {
      page: data.pagination?.page ?? fallbackPage,
      page_size: data.pagination?.page_size ?? data.pagination?.pageSize ?? fallbackPageSize,
      total_items: data.pagination?.total_items ?? data.pagination?.totalItems ?? data.items?.length ?? 0,
      total_pages: data.pagination?.total_pages ?? data.pagination?.totalPages ?? 1,
    },
  };
}

export function useEcosystemBatchesQuery(filters) {
  return useQuery({
    queryKey: ["ecosystem-batches", filters],
    queryFn: async () => withDefaultPagination(await api.getEcosystemBatches(filters), filters.page, filters.pageSize),
    refetchInterval: (query) =>
      query.state.data?.items?.some((item) => isActiveBatchStatus(item.status)) ? 5000 : false,
    refetchIntervalInBackground: true,
  });
}

export function useEcosystemBatchQuery(batchId) {
  return useQuery({
    queryKey: ["ecosystem-batch", batchId],
    queryFn: () => api.getEcosystemBatch(batchId),
    enabled: Boolean(batchId),
    staleTime: 0,
  });
}

export function useEcosystemBatchSummaryQuery(batchId, shouldPoll) {
  return useQuery({
    queryKey: ["ecosystem-batch-summary", batchId],
    queryFn: () => api.getEcosystemBatchSummary(batchId),
    enabled: Boolean(batchId),
    staleTime: 0,
    refetchInterval: shouldPoll ? 4000 : false,
    refetchIntervalInBackground: true,
  });
}

export function useEcosystemBatchReposQuery(batchId, filters, shouldPoll) {
  return useQuery({
    queryKey: ["ecosystem-batch-repos", batchId, filters],
    queryFn: async () =>
      withDefaultPagination(await api.getEcosystemBatchRepos(batchId, filters), filters.page, filters.pageSize),
    enabled: Boolean(batchId),
    staleTime: 0,
    placeholderData: (previous) => previous,
    refetchInterval: shouldPoll ? 4000 : false,
    refetchIntervalInBackground: true,
  });
}
