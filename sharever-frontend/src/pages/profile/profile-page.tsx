import { User, Mail, Shield, Bell, LogOut, ChevronRight, CreditCard } from "lucide-react";
import { DEFAULT_AVATAR_URL } from "../../shared/lib/default-avatar";

const ProfileItem = ({ icon: Icon, label, value, danger }: any) => (
  <div className={`flex items-center justify-between p-4 rounded-2xl transition-colors cursor-pointer ${danger ? 'hover:bg-red-50' : 'hover:bg-gray-50'}`}>
    <div className="flex items-center gap-4">
      <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${danger ? 'bg-red-100 text-red-500' : 'bg-gray-100 text-gray-600'}`}>
        <Icon size={20} />
      </div>
      <span className={`font-semibold ${danger ? 'text-red-500' : 'text-gray-700'}`}>{label}</span>
    </div>
    <div className="flex items-center gap-3">
       {value && <span className="text-sm text-gray-400">{value}</span>}
       <ChevronRight size={16} className="text-gray-300" />
    </div>
  </div>
);

export default function ProfilePage() {
  return (
    <div className="max-w-2xl mx-auto animate-enter py-6">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Account</h1>

      {/* Avatar Section */}
      <div className="bg-white p-6 rounded-[32px] border border-gray-100 shadow-sm mb-6 flex items-center gap-6">
        <div className="w-24 h-24 rounded-full p-1 border-2 border-dashed border-purple-200">
            <img 
                src={DEFAULT_AVATAR_URL}
                alt="Profile" 
                className="w-full h-full rounded-full object-cover"
            />
        </div>
        <div>
            <h2 className="text-2xl font-bold text-gray-900">Jimmy Nguyen</h2>
            <p className="text-gray-500">jimmy@thuanle.me</p>
            <button className="text-sm text-purple-600 font-bold mt-2 hover:underline">Change Picture</button>
        </div>
      </div>

      {/* Settings List */}
      <div className="bg-white p-2 rounded-[32px] border border-gray-100 shadow-sm space-y-1">
         <ProfileItem icon={User} label="Personal Information" />
         <ProfileItem icon={CreditCard} label="Payment Methods" value="Visa **4242" />
         <ProfileItem icon={Bell} label="Notifications" value="On" />
         <ProfileItem icon={Shield} label="Security & Privacy" />
      </div>

      <div className="mt-6 bg-white p-2 rounded-[32px] border border-gray-100 shadow-sm">
         <ProfileItem icon={LogOut} label="Log out" danger />
      </div>
    </div>
  );
}
