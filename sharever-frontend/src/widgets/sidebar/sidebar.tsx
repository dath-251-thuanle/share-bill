import { useEffect } from "react";
import { NavLink, Link, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Home, Activity, CreditCard, Users, Plus } from "lucide-react";

import { eventApi } from "../../entities/event/api";
import { useEventStore } from "../../stores/use-event-store";

// 1) Menu với Icon Lucide
const navItems = [
  { to: "/app", end: true, label: "Home", icon: Home },
  { to: "/app/activity", label: "Activity", icon: Activity },
  { to: "/app/expenses", label: "Expenses", icon: CreditCard },
  { to: "/app/accounts", label: "Accounts", icon: Users },
];

export function Sidebar() {
  const navigate = useNavigate();
  const { selectedEventId, setSelectedEventId } = useEventStore();

  const { data: events = [], isLoading } = useQuery({
    queryKey: ["events"],
    queryFn: eventApi.list,
  });

  // auto-select event đầu tiên nếu chưa có
  useEffect(() => {
    if (!events.length) {
      if (selectedEventId) setSelectedEventId(null);
      return;
    }

    const match = events.find(
      (event: any) => String(event.id) === String(selectedEventId)
    );
    if (!match) {
      setSelectedEventId(String(events[0].id));
    }
  }, [events, selectedEventId, setSelectedEventId]);

  function handleSelect(eventId: string) {
    setSelectedEventId(String(eventId));
    // chọn event xong -> nhảy sang expenses (giống Codex)
    navigate("/app/activity");
  }

  return (
    <aside className="w-[280px] h-full bg-white border-r border-gray-100 flex flex-col z-20">
      <div className="flex flex-col h-full p-6">
        {/* LOGO */}
        <Link
          to="/app"
          className="flex items-center gap-3 mb-10 cursor-pointer hover:opacity-80 transition-opacity"
        >
          <div className="h-10 w-10 rounded-full bg-gradient-to-tr from-teal-400 to-blue-500 flex items-center justify-center text-white font-bold shadow-lg shadow-blue-200">
            S
          </div>
          <div className="text-2xl font-bold italic tracking-tight font-display text-gray-900">
            Sharever
          </div>
        </Link>

        {/* NAV */}
        <nav className="space-y-2 flex-1">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                `relative flex items-center gap-4 px-4 py-3 rounded-xl transition-all duration-200 group ${
                  isActive
                    ? "bg-gray-100 text-gray-900 font-bold"
                    : "text-gray-500 hover:text-gray-900 hover:bg-gray-50"
                }`
              }
            >
              {({ isActive }) => (
                <>
                  <item.icon
                    size={24}
                    strokeWidth={isActive ? 2.5 : 2}
                    className="transition-all"
                  />
                  <span className="text-sm font-medium">{item.label}</span>

                  {isActive && (
                    <div className="absolute left-0 top-1/2 -translate-y-1/2 h-3/5 w-1 bg-gray-900 rounded-r-full" />
                  )}
                </>
              )}
            </NavLink>
          ))}
        </nav>

        {/* EVENTS */}
        <div className="mt-8">
          <div className="flex items-center justify-between mb-4">
            <div className="text-lg font-bold text-gray-900">Your events</div>
            <button
              className="bg-purple-600 text-white px-3 py-1.5 rounded-lg text-xs font-bold flex items-center gap-1 shadow-lg shadow-purple-200 hover:scale-105 transition-transform"
              onClick={() => navigate("/app?create=1")}
            >
              <Plus size={14} /> New
            </button>
          </div>

          {isLoading && (
            <div className="text-xs text-gray-400">Loading events...</div>
          )}

          {!isLoading && events.length === 0 && (
            <div className="text-xs text-gray-400">No events yet.</div>
          )}

          <div className="space-y-2">
            {events.map((event: any) => {
              const isActiveEvent = String(event.id) === String(selectedEventId);
              const iconChar = (event.name || "?").slice(0, 1).toUpperCase();

              return (
                <button
                  key={event.id}
                  onClick={() => handleSelect(event.id)}
                  className={`w-full flex items-center gap-3 px-3 py-3 rounded-xl text-left transition-all ${
                    isActiveEvent
                      ? "bg-purple-50 text-purple-700 font-semibold border border-purple-100"
                      : "bg-transparent text-gray-600 border border-transparent hover:bg-gray-50"
                  }`}
                >
                  <div
                    className={`h-8 w-8 rounded-lg flex items-center justify-center text-sm ${
                      isActiveEvent ? "bg-purple-200" : "bg-gray-200"
                    }`}
                  >
                    {iconChar}
                  </div>
                  <span className="text-sm truncate">{event.name}</span>
                </button>
              );
            })}
          </div>
        </div>
      </div>
    </aside>
  );
}
