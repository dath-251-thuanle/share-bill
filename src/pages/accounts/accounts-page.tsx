import { useState, type FormEvent } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  User, Mail, Shield, Bell, LogOut, CreditCard,
  ChevronRight, Smartphone, Globe, Moon, Edit3, Camera
} from "lucide-react";

import { useAuth } from "../../features/auth/model/use-auth";
import { userApi } from "../../entities/user/api";
import { eventApi } from "../../entities/event/api";
import { participantApi } from "../../entities/participant/api";
import { Modal } from "../../shared/ui/modal";
import { Input } from "../../shared/ui/input";
import { Button } from "../../shared/ui/button";
import { useToast } from "../../shared/ui/toast";
import { normalizeError } from "../../shared/lib/errors";
import { DEFAULT_AVATAR_URL } from "../../shared/lib/default-avatar";
import {
  BankInfoForm,
  type BankInfoPayload,
} from "../../features/bank-update/ui/bank-info-form";

// --- Components nhỏ cho gọn code ---

// 1. Mục cài đặt đơn lẻ
const SettingItem = ({ icon: Icon, label, value, danger, toggle, onClick }: any) => (
  <div
    className={`flex items-center justify-between p-4 rounded-2xl transition-all cursor-pointer group ${danger ? 'hover:bg-red-50' : 'hover:bg-gray-50'}`}
    onClick={onClick}
  >
    <div className="flex items-center gap-4">
      <div className={`w-10 h-10 rounded-xl flex items-center justify-center transition-colors ${danger ? 'bg-red-100 text-red-500 group-hover:bg-red-200' : 'bg-gray-100 text-gray-600 group-hover:bg-white group-hover:shadow-sm'}`}>
        <Icon size={20} />
      </div>
      <span className={`font-semibold ${danger ? 'text-red-500' : 'text-gray-700'}`}>{label}</span>
    </div>
    
    <div className="flex items-center gap-3">
       {value && <span className="text-sm text-gray-400 font-medium">{value}</span>}
       {toggle ? (
         <div className="w-11 h-6 bg-purple-600 rounded-full relative cursor-pointer">
            <div className="absolute right-1 top-1 w-4 h-4 bg-white rounded-full shadow-sm"></div>
         </div>
       ) : (
         <ChevronRight size={18} className="text-gray-300 group-hover:text-gray-500" />
       )}
    </div>
  </div>
);

// 2. Thẻ Credit Card giả lập (Visual)
const CreditCardVisual = () => (
  <div className="bg-gradient-to-br from-[#2a2a2a] to-black rounded-[24px] p-6 text-white relative overflow-hidden shadow-xl group hover:scale-[1.02] transition-transform duration-300">
     {/* Decoration circles */}
     <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full blur-2xl -mr-10 -mt-10"></div>
     <div className="absolute bottom-0 left-0 w-24 h-24 bg-purple-500/20 rounded-full blur-xl -ml-5 -mb-5"></div>
     
     <div className="relative z-10 flex flex-col justify-between h-40">
        <div className="flex justify-between items-start">
           <div className="text-xs font-bold tracking-widest text-white/60">Debit Card</div>
           <div className="font-bold italic text-lg">VISA</div>
        </div>
        
        <div>
           <div className="text-2xl font-mono tracking-widest mb-1">•••• •••• •••• 4242</div>
           <div className="flex justify-between items-end">
              <div>
                 <div className="text-[10px] text-white/50 uppercase">Card Holder</div>
                 <div className="font-bold text-sm">JIMMY NGUYEN</div>
              </div>
              <div className="text-right">
                 <div className="text-[10px] text-white/50 uppercase">Expires</div>
                 <div className="font-bold text-sm">12/28</div>
              </div>
           </div>
        </div>
     </div>
  </div>
);

// --- MAIN PAGE ---
export default function AccountsPage() {
  const user = useAuth((s) => s.user);
  const setUser = useAuth((s) => s.setUser);
  const queryClient = useQueryClient();
  const toast = useToast();
  const [editOpen, setEditOpen] = useState(false);
  const [bankOpen, setBankOpen] = useState(false);
  const [passwordOpen, setPasswordOpen] = useState(false);
  const [fullName, setFullName] = useState(user?.name ?? "");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [passwordLoading, setPasswordLoading] = useState(false);
  const [currentPassword, setCurrentPassword] = useState("");
  const [nextPassword, setNextPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const displayName = user?.name ?? "Account";
  const displayEmail = user?.email ?? "Unknown";
  const bankInfo = user?.bankInfo ?? {
    bankName: user?.bankName ?? "",
    accountNumber: user?.accountNumber ?? "",
    accountName: user?.accountName ?? "",
  };
  const hasBankInfo =
    !!bankInfo?.bankName && !!bankInfo?.accountNumber && !!bankInfo?.accountName;

  async function syncParticipantNames(userId: string, name: string) {
    if (!userId) return;
    let hadError = false;

    try {
      const events = await eventApi.list();

      await Promise.all(
        (events as any[]).map(async (event) => {
          try {
            const participantsData = await participantApi.list(String(event.id));
            const participants = Array.isArray(participantsData)
              ? participantsData
              : (participantsData as any)?.participants ?? [];
            const match = participants.find(
              (p: any) => String(p.userId) === String(userId)
            );
            if (match && match.name !== name) {
              await participantApi.update(String(match.id), { name });
            }
          } catch {
            hadError = true;
          }
        })
      );

      await queryClient.invalidateQueries({ queryKey: ["participants"] });
      await queryClient.invalidateQueries({ queryKey: ["transactions"] });
    } catch {
      hadError = true;
    }

    if (hadError) {
      toast.push("Name updated, but some events did not sync.");
    }
  }

  async function syncParticipantBankInfo(
    userId: string,
    bankPayload: BankInfoPayload
  ) {
    if (!userId) return;
    let hadError = false;

    try {
      const events = await eventApi.list();

      await Promise.all(
        (events as any[]).map(async (event) => {
          try {
            const participantsData = await participantApi.list(String(event.id));
            const participants = Array.isArray(participantsData)
              ? participantsData
              : (participantsData as any)?.participants ?? [];
            const match = participants.find(
              (p: any) => String(p.userId) === String(userId)
            );
            if (match) {
              await participantApi.update(String(match.id), {
                name: match.name ?? user?.name ?? "Member",
                bankInfo: {
                  bankName: bankPayload.bankName,
                  accountNumber: bankPayload.accountNumber,
                  accountName: bankPayload.accountName,
                },
              });
            }
          } catch {
            hadError = true;
          }
        })
      );

      await queryClient.invalidateQueries({ queryKey: ["participants"] });
    } catch {
      hadError = true;
    }

    if (hadError) {
      toast.push("Bank info updated, but some events did not sync.");
    }
  }

  async function handleSaveName(e: FormEvent) {
    e.preventDefault();
    const nextName = fullName.trim();
    if (!nextName) {
      setError("Full name is required.");
      return;
    }

    setSaving(true);
    setError(null);
    try {
      const updated = await userApi.updateProfile({ name: nextName });
      setUser(updated);
      await syncParticipantNames(String(updated?.id ?? ""), nextName);
      toast.push("Name updated.");
      setEditOpen(false);
    } catch (err) {
      setError(normalizeError(err));
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="animate-enter space-y-8 pb-20">
      
      {/* Header */}
      <div>
        <h1 className="text-3xl font-extrabold text-gray-900 font-display">Account Settings</h1>
        <p className="text-gray-500 mt-1">Manage your profile, payment methods and preferences.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        
        {/* LEFT COLUMN: Profile & Wallet (Chiếm 4 phần) */}
        <div className="lg:col-span-4 space-y-6">
           
           {/* 1. Profile Card */}
           <div className="bg-white p-6 rounded-[32px] border border-gray-100 shadow-sm text-center relative overflow-hidden">
              <div className="absolute top-0 left-0 w-full h-24 bg-gradient-to-r from-teal-400 to-purple-400 opacity-20"></div>
              
              <div className="relative z-10 mt-8 mb-4">
                 <div className="w-28 h-28 mx-auto rounded-full p-1 bg-white shadow-md relative group cursor-pointer">
                    <img 
                        src={user?.avatar ?? user?.avatarUrl ?? DEFAULT_AVATAR_URL}
                        alt="Profile" 
                        className="w-full h-full rounded-full object-cover border-4 border-white"
                    />
                    <div className="absolute bottom-0 right-1 bg-gray-900 text-white p-2 rounded-full hover:bg-purple-600 transition-colors">
                        <Camera size={14} />
                    </div>
                 </div>
              </div>
              
              <h2 className="text-xl font-bold text-gray-900">{displayName}</h2>
              <p className="text-gray-500 text-sm mb-6">{displayEmail}</p>
              
              <div className="flex justify-center gap-4 mb-6">
                 <div className="text-center px-4 py-2 bg-gray-50 rounded-2xl">
                    <div className="font-bold text-gray-900">14</div>
                    <div className="text-xs text-gray-400 font-bold uppercase">Friends</div>
                 </div>
                 <div className="text-center px-4 py-2 bg-gray-50 rounded-2xl">
                    <div className="font-bold text-purple-600">8</div>
                    <div className="text-xs text-gray-400 font-bold uppercase">Groups</div>
                 </div>
              </div>

              <button
                className="w-full py-3 rounded-xl border border-gray-200 font-bold text-gray-700 hover:bg-gray-50 flex items-center justify-center gap-2 transition-colors"
                onClick={() => {
                  setFullName(user?.name ?? "");
                  setError(null);
                  setEditOpen(true);
                }}
              >
                 <Edit3 size={16} /> Edit name
              </button>
           </div>

           {/* 2. My Wallet */}
           <div className="bg-white p-6 rounded-[32px] border border-gray-100 shadow-sm">
              <div className="flex justify-between items-center mb-4">
                 <h3 className="font-bold text-gray-900">My Wallet</h3>
                 <button className="text-purple-600 text-sm font-bold hover:underline">Manage</button>
              </div>
              <CreditCardVisual />
              <div className="mt-4 flex gap-2">
                 <div className="flex-1 bg-gray-50 p-3 rounded-xl flex items-center gap-2 cursor-pointer hover:bg-gray-100">
                    <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold text-xs">P</div>
                    <span className="text-sm font-semibold text-gray-700">PayPal</span>
                 </div>
                 <div className="flex-1 bg-gray-50 p-3 rounded-xl flex items-center gap-2 cursor-pointer hover:bg-gray-100 border border-dashed border-gray-300">
                    <div className="w-8 h-8 rounded-full bg-gray-200 text-gray-400 flex items-center justify-center">+</div>
                    <span className="text-sm font-semibold text-gray-500">Add</span>
                 </div>
              </div>
           </div>

           <div className="bg-white p-6 rounded-[32px] border border-gray-100 shadow-sm">
              <div className="flex items-center justify-between">
                 <h3 className="font-bold text-gray-900">Bank transfer</h3>
                 <button
                   className="text-purple-600 text-sm font-bold hover:underline"
                   onClick={() => setBankOpen(true)}
                 >
                   {hasBankInfo ? "Edit" : "Add"}
                 </button>
              </div>
              {hasBankInfo ? (
                <div className="mt-4 space-y-1 text-sm">
                  <div className="font-semibold text-gray-900">
                    {bankInfo.bankName}
                  </div>
                  <div className="text-gray-700">{bankInfo.accountNumber}</div>
                  <div className="text-gray-500">{bankInfo.accountName}</div>
                </div>
              ) : (
                <div className="mt-4 text-sm text-gray-500">
                  Add bank info to receive VietQR payments.
                </div>
              )}
              <div className="mt-3 text-xs text-gray-400">
                Use VietQR bank code (VCB, ACB, TCB).
              </div>
           </div>

        </div>

        {/* RIGHT COLUMN: Settings List (Chiếm 8 phần) */}
        <div className="lg:col-span-8 space-y-6">
           
           {/* Group 1: General */}
           <div className="bg-white p-2 rounded-[32px] border border-gray-100 shadow-sm">
              <div className="px-6 py-4 border-b border-gray-50">
                 <h3 className="font-bold text-gray-900">General</h3>
              </div>
              <div className="p-2 space-y-1">
                 <SettingItem icon={User} label="Personal Information" value="Updated" />
                 <SettingItem icon={Smartphone} label="My QR Code" value="Ready" />
                 <SettingItem icon={Globe} label="Language" value="English (US)" />
                 <SettingItem icon={Moon} label="Dark Mode" toggle={true} />
              </div>
           </div>

           {/* Group 2: Notifications & Security */}
           <div className="bg-white p-2 rounded-[32px] border border-gray-100 shadow-sm">
              <div className="px-6 py-4 border-b border-gray-50">
                 <h3 className="font-bold text-gray-900">Security & Notifications</h3>
              </div>
              <div className="p-2 space-y-1">
                 <SettingItem icon={Bell} label="Push Notifications" toggle={true} />
                 <SettingItem icon={Mail} label="Email Newsletter" toggle={false} />
                 <SettingItem
                   icon={Shield}
                   label="Change Password"
                   onClick={() => {
                     setPasswordError(null);
                     setCurrentPassword("");
                     setNextPassword("");
                     setConfirmPassword("");
                     setPasswordOpen(true);
                   }}
                 />
                 <SettingItem icon={CreditCard} label="Payment Security" value="FaceID" />
              </div>
           </div>

           {/* Group 3: Danger Zone */}
           <div className="bg-white p-2 rounded-[32px] border border-gray-100 shadow-sm">
               <SettingItem icon={LogOut} label="Log out" danger />
           </div>

           <div className="text-center text-sm text-gray-400 pt-4">
              Sharever App v1.0.2 • Built with ❤️ by Jimmy
           </div>

        </div>

      </div>

      <Modal
        open={editOpen}
        onClose={() => {
          setEditOpen(false);
          setError(null);
        }}
        title="Update full name"
      >
        <form className="space-y-4" onSubmit={handleSaveName}>
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">Full name</label>
            <Input
              placeholder="Your full name"
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              required
            />
          </div>

          {error && <div className="text-sm text-rose-600">{error}</div>}

          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="ghost"
              className="text-gray-600"
              onClick={() => setEditOpen(false)}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={saving}>
              {saving ? "Saving..." : "Save"}
            </Button>
          </div>
        </form>
      </Modal>

      <BankInfoForm
        open={bankOpen}
        onClose={() => setBankOpen(false)}
        initial={{
          bankName: bankInfo?.bankName ?? "",
          accountNumber: bankInfo?.accountNumber ?? "",
          accountName: bankInfo?.accountName ?? "",
        }}
        onSubmit={async (payload) => {
          const bankPayload = {
            bankName: payload.bankName.trim(),
            accountNumber: payload.accountNumber.trim(),
            accountName: payload.accountName.trim(),
          };
          const updated = await userApi.updateProfile(bankPayload);
          setUser(updated);
          await syncParticipantBankInfo(String(updated?.id ?? ""), bankPayload);
        }}
      />

      <Modal
        open={passwordOpen}
        onClose={() => {
          setPasswordOpen(false);
          setPasswordError(null);
        }}
        title="Change password"
      >
        <form
          className="space-y-4"
          onSubmit={async (e) => {
            e.preventDefault();
            if (!currentPassword || !nextPassword || !confirmPassword) {
              setPasswordError("Please fill in all fields.");
              return;
            }
            if (nextPassword.length < 6) {
              setPasswordError("New password must be at least 6 characters.");
              return;
            }
            if (nextPassword !== confirmPassword) {
              setPasswordError("Passwords do not match.");
              return;
            }

            setPasswordLoading(true);
            setPasswordError(null);
            try {
              await userApi.changePassword({
                currentPassword,
                newPassword: nextPassword,
              });
              toast.push("Password updated.");
              setPasswordOpen(false);
            } catch (err) {
              setPasswordError(normalizeError(err));
            } finally {
              setPasswordLoading(false);
            }
          }}
        >
          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">
              Current password
            </label>
            <Input
              type="password"
              placeholder="••••••••"
              value={currentPassword}
              onChange={(e) => setCurrentPassword(e.target.value)}
              required
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">
              New password
            </label>
            <Input
              type="password"
              placeholder="At least 6 characters"
              value={nextPassword}
              onChange={(e) => setNextPassword(e.target.value)}
              required
            />
          </div>

          <div className="space-y-1">
            <label className="text-sm font-medium text-gray-700">
              Confirm new password
            </label>
            <Input
              type="password"
              placeholder="Repeat new password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />
          </div>

          {passwordError && (
            <div className="text-sm text-rose-600">{passwordError}</div>
          )}

          <div className="flex justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="ghost"
              className="text-gray-600"
              onClick={() => setPasswordOpen(false)}
              disabled={passwordLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={passwordLoading}>
              {passwordLoading ? "Saving..." : "Update password"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
