import { useEffect, useMemo, useState } from "react";
import { Plus, Calendar, ArrowRight, TrendingUp, TrendingDown } from "lucide-react";
import { useQuery, useQueries } from "@tanstack/react-query";
import { useNavigate, useSearchParams } from "react-router-dom";

import { useAuth } from "../../features/auth/model/use-auth";
import { useEventStore } from "../../stores/use-event-store";
import { eventApi } from "../../entities/event/api";
import { settlementApi } from "../../entities/settlement/api";
import { participantApi } from "../../entities/participant/api";
import { Modal } from "../../shared/ui/modal";
import { Input } from "../../shared/ui/input";
import { Button } from "../../shared/ui/button";
import { useToast } from "../../shared/ui/toast";
import { normalizeError } from "../../shared/lib/errors";


export default function EventsListPage() {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const user = useAuth((s) => s.user);
  const { selectedEventId: activeEventId, setSelectedEventId: setActiveEventId } =
    useEventStore();
  const toast = useToast();

  // chọn event bằng id (không truyền object hard-code nữa)
  const createOpen = searchParams.get("create") === "1";
  const [createName, setCreateName] = useState("");
  const [createCurrency, setCreateCurrency] = useState("VND");
  const [createDescription, setCreateDescription] = useState("");
  const [createError, setCreateError] = useState<string | null>(null);
  const [createLoading, setCreateLoading] = useState(false);

  const openCreateModal = () => {
    const next = new URLSearchParams(searchParams);
    next.set("create", "1");
    setSearchParams(next, { replace: true });
  };

  const closeCreateModal = () => {
    setCreateError(null);
    const next = new URLSearchParams(searchParams);
    if (next.has("create")) {
      next.delete("create");
      setSearchParams(next, { replace: true });
    }
  };

  // fetch events thật từ BE (đã auto attach token nhờ interceptor)
  const {
    data: events = [],
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery({
    queryKey: ["events"],
    queryFn: eventApi.list,
  });

  useEffect(() => {
    if (!activeEventId && events.length) {
      setActiveEventId(String((events as any[])[0].id));
    }
  }, [activeEventId, events, setActiveEventId]);

  const eventList = events as any[];

  const balanceQueries = useQueries({
    queries: eventList.map((event) => ({
      queryKey: ["event-balance", event.id, user?.id],
      queryFn: async () => {
        const [summary, participantsData] = await Promise.all([
          settlementApi.summary(String(event.id)),
          participantApi.list(String(event.id)),
        ]);
        const participants = Array.isArray(participantsData)
          ? participantsData
          : (participantsData as any)?.participants ?? [];
        const myParticipant = participants.find((p: any) => p.userId === user?.id);
        const balance =
          (summary?.participants ?? []).find((p: any) => p.id === myParticipant?.id)
            ?.balance ?? 0;
        return {
          eventId: String(event.id),
          balance,
          currency: summary?.event?.currency ?? event.currency ?? "VND",
        };
      },
      enabled: !!user?.id,
    })),
  });

  const { youOwe, youAreOwed, currencyLabel, totalsReady } = useMemo(() => {
    let owe = 0;
    let owed = 0;
    const currencies = new Set<string>();
    let hasData = eventList.length === 0;

    balanceQueries.forEach((q) => {
      const balance = q.data?.balance;
      const currency = q.data?.currency;
      if (typeof balance === "number") {
        hasData = true;
        if (balance < 0) owe += Math.abs(balance);
        if (balance > 0) owed += balance;
      }
      if (currency) currencies.add(currency);
    });

    const currencyLabel =
      currencies.size === 1
        ? Array.from(currencies)[0]
        : currencies.size === 0
        ? "VND"
        : "Multi";

    return { youOwe: owe, youAreOwed: owed, currencyLabel, totalsReady: hasData };
  }, [balanceQueries, eventList.length]);

  const fmtMoney = (v: number) =>
    new Intl.NumberFormat("en-US", { maximumFractionDigits: 0 }).format(v);

  async function handleCreateEvent(e: React.FormEvent) {
    e.preventDefault();
    if (!createName.trim()) {
      setCreateError("Event name is required.");
      return;
    }
    setCreateError(null);
    setCreateLoading(true);
    try {
      const created = await eventApi.create({
        name: createName.trim(),
        currency: createCurrency,
        description: createDescription.trim() || undefined,
      });
      toast.push("Event created.");
      closeCreateModal();
      setCreateName("");
      setCreateDescription("");
      await refetch();
      const createdEvent = (created as any)?.event ?? created;
      if (createdEvent?.id) {
        const id = String(createdEvent.id);
        setActiveEventId(id);
        navigate("/app/activity");
      }
    } catch (err) {
      setCreateError(normalizeError(err));
    } finally {
      setCreateLoading(false);
    }
  }

  // nếu đang xem detail -> render detail page (detail sẽ tự fetch theo id)
  return (
    <div className="space-y-8 animate-enter">
      {/* Header & Stats */}
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight">
            Welcome back, {user?.name ?? "buddy"}! 👋
          </h1>
          <p className="text-gray-500 mt-2">
            Here is what's happening with your expenses.
          </p>
        </div>

        <button
          className="bg-gray-900 text-white px-5 py-3 rounded-2xl font-semibold shadow-lg shadow-gray-200 hover:scale-105 transition-transform flex items-center gap-2"
          onClick={openCreateModal}
        >
          <Plus size={20} /> Create Event
        </button>
      </div>

      {/* Stats Cards (tạm giữ hard-code, sau này nối summary endpoint) */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
        <div className="bg-[#FAEBE6] p-6 rounded-[32px] relative overflow-hidden group">
          <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:scale-110 transition-transform">
            <TrendingDown size={80} />
          </div>
          <p className="text-gray-600 font-medium text-sm">You owe (all events)</p>
          <h3 className="text-3xl font-bold text-gray-900 mt-2">
            - {totalsReady ? fmtMoney(youOwe) : "--"}{" "}
            <span className="text-sm font-normal text-gray-500">{currencyLabel}</span>
          </h3>
        </div>

        <div className="bg-[#E6F6FA] p-6 rounded-[32px] relative overflow-hidden group">
          <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:scale-110 transition-transform">
            <TrendingUp size={80} />
          </div>
          <p className="text-gray-600 font-medium text-sm">You are owed (all events)</p>
          <h3 className="text-3xl font-bold text-emerald-600 mt-2">
            + {totalsReady ? fmtMoney(youAreOwed) : "--"}{" "}
            <span className="text-sm font-normal text-gray-500">{currencyLabel}</span>
          </h3>
        </div>
      </div>

      {/* Events Grid */}
      <div>
        <h2 className="text-xl font-bold text-gray-800 mb-6">Your Events</h2>

        {isLoading && <div className="text-gray-500">Loading events...</div>}

        {isError && (
          <div className="rounded-xl border border-rose-200 bg-rose-50 p-3 text-rose-700">
            Failed to load events: {(error as any)?.message ?? "Unknown error"}{" "}
            <button className="underline" onClick={() => refetch()}>
              Retry
            </button>
          </div>
        )}

        {!isLoading && !isError && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {events.map((event: any, i: number) => {
              // Backend thường có id + name + createdAt/created_at
              const eventName = event?.name ?? event?.title ?? "Untitled";
              const rawDate = event?.createdAt ?? event?.created_at ?? event?.date;
              const dateLabel = rawDate
                ? new Date(rawDate).toLocaleDateString("vi-VN")
                : "-";

              return (
                <div
                  key={event.id}
                  onClick={() => {
                    const id = String(event.id);
                    setActiveEventId(id);
                    navigate("/app/activity");
                  }}
                  className="bg-white p-6 rounded-[32px] border border-gray-100 hover:shadow-xl hover:-translate-y-1 transition-all cursor-pointer group animate-enter"
                  // NOTE: tailwind class delay-${i*100} không build được; dùng inline delay cho chắc
                  style={{ animationDelay: `${i * 100}ms` }}
                >
                  <div className="flex justify-between items-start mb-6">
                    <div className="w-14 h-14 rounded-2xl flex items-center justify-center text-2xl bg-gray-100 text-gray-700">
                      🎯
                    </div>
                    <span className="px-3 py-1 rounded-full bg-gray-50 text-xs font-bold uppercase tracking-wider text-gray-500 border border-gray-100">
                      EVENT
                    </span>
                  </div>

                  <h3 className="text-xl font-bold text-gray-900 group-hover:text-purple-600 transition-colors">
                    {eventName}
                  </h3>

                  <div className="flex items-center gap-2 text-gray-400 text-sm mt-2 mb-6">
                    <Calendar size={14} />
                    <span>{dateLabel}</span>
                  </div>

                  <div className="flex items-center justify-between pt-4 border-t border-gray-50">
                    <div className="flex -space-x-2">
                      {[1, 2, 3].map((u) => (
                        <div
                          key={u}
                          className="w-8 h-8 rounded-full border-2 border-white bg-gray-200"
                        />
                      ))}
                    </div>
                    <div className="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-gray-400 group-hover:bg-purple-100 group-hover:text-purple-600 transition-colors">
                      <ArrowRight size={16} />
                    </div>
                  </div>
                </div>
              );
            })}

            {/* Add New Placeholder */}
            <div
              className="border-2 border-dashed border-gray-200 rounded-[32px] flex flex-col items-center justify-center min-h-[200px] text-gray-400 hover:border-purple-300 hover:bg-purple-50 hover:text-purple-600 transition-all cursor-pointer"
              onClick={openCreateModal}
            >
              <Plus size={32} />
              <span className="font-semibold mt-2">New Event</span>
            </div>
          </div>
        )}
      </div>

      <Modal
        open={createOpen}
        onClose={closeCreateModal}
        title="Create event"
      >
        <form onSubmit={handleCreateEvent} className="space-y-4">
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">Event name</label>
            <Input
              placeholder="Da Lat trip"
              value={createName}
              onChange={(e) => setCreateName(e.target.value)}
              required
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">Currency</label>
            <select
              className="h-11 w-full rounded-2xl bg-gray-100 px-4 text-sm text-gray-800 outline-none border border-transparent focus:border-purple-400 focus:bg-white"
              value={createCurrency}
              onChange={(e) => setCreateCurrency(e.target.value)}
            >
              <option value="VND">VND</option>
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
            </select>
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">Description</label>
            <textarea
              className="w-full min-h-[96px] rounded-2xl bg-gray-100 px-4 py-3 text-sm text-gray-800 outline-none border border-transparent focus:border-purple-400 focus:bg-white"
              placeholder="Optional details..."
              value={createDescription}
              onChange={(e) => setCreateDescription(e.target.value)}
            />
          </div>

          {createError && (
            <div className="text-sm text-rose-600">{createError}</div>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="ghost"
              className="text-gray-600"
              onClick={closeCreateModal}
              disabled={createLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createLoading}>
              {createLoading ? "Creating..." : "Create"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
